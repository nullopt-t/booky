package app

import (
	// "booky-backend/internal/cart"
	"booky-backend/internal/http/swagger"
	"booky-backend/internal/middleware"
	"booky-backend/internal/notifier"
	"booky-backend/internal/shared/jwt"
	"booky-backend/internal/user/otp"

	// "booky-backend/internal/inventory"
	// "booky-backend/internal/product"
	"booky-backend/internal/user"

	// "booky-backend/internal/checkout"
	// "booky-backend/internal/order"

	"booky-backend/pkg/api/security"
	"booky-backend/pkg/config"
	"booky-backend/pkg/database"
	"booky-backend/pkg/log"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type application interface {
	Run() error
	Shutdown()
}

type App struct {
	// http server
	server *http.Server

	// logger
	logger *log.ConsoleLogger

	// cache
	redisClient *redis.Client

	// database
	db *database.DB
}

func (app *App) setupRoutes(config *config.Config, router *gin.Engine) {
	// setup middlewares
	router.Use(middleware.ErrorHandler(app.logger))

	v1 := router.Group("/api/v1")
	swagger.SetUpDocs(v1)

	txRunner := database.NewTxRunner(app.db)

	jobQueue := notifier.NewRedisJobQueue(
		app.redisClient,
	)

	notifier := notifier.NewNotifier(
		jobQueue,
		app.logger,
	)
	// otp
	otpRepo := otp.NewOTPStore(
		app.redisClient,
		app.logger,
	)

	rateLimiter := security.NewRateLimiter(
		app.redisClient,
	)
	otpGen := otp.NewOTPGenerator()

	otpService := otp.NewService(
		otpRepo,
		otpGen,
		app.logger,
		notifier,
	)
	// user
	userRepo := user.NewPostgresRepository()
	userService := user.NewService(
		txRunner,
		otpService,
		userRepo,
		app.logger,
	)

	jwtService := jwt.NewJWTManager(
		config.KeysCfg,
	)

	authService := user.NewAuthService(
		app.logger,
		userService,
		jwtService,
		otpService,
		app.redisClient,
		notifier,
	)

	userHandler := user.NewHandler(
		userService,
		authService,
		rateLimiter,
		config.KeysCfg,
		app.logger,
	)

	userRouter := user.NewRouter(
		userHandler,
		config,
	)

	userRouter.MapRoutes(v1)

	// inventoryRepo := inventory.NewPostgresRepository()
	// inventoryService := inventory.NewService(txRunner, inventoryRepo)
	// inventoryHandler := inventory.NewHandler(inventoryService)
	// inventoryRouter := inventory.NewRouter(inventoryHandler)
	// inventoryRouter.MapRoutes(v1.Group("/inventories"))

	// // product
	// productRepo := product.NewPostgresRepository()
	// productService := product.NewService(txRunner, productRepo, inventoryRepo)
	// productHandler := product.NewHandler(productService)
	// productRouter := product.NewRouter(productHandler)
	// productRouter.MapRoutes(v1.Group("/products"))

	// // cart
	// cartRepo := cart.NewPostgresRepository()
	// cartService := cart.NewService(txRunner, cartRepo, productRepo)
	// cartHandler := cart.NewHandler(cartService)
	// cartRouter := cart.NewRouter(cartHandler)
	// cartRouter.MapRoutes(v1.Group("/carts"))

	// // order
	// orderRepo := order.NewPostgresRepository()
	// orderService := order.NewService(txRunner, orderRepo)
	// orderHandler := order.NewHandler(orderService)
	// orderRouter := order.NewRouter(orderHandler)
	// orderRouter.MapRoutes(v1.Group("/orders"))

	// // checkout
	// checkoutService := checkout.NewService(app.db.GetPool(), orderRepo, cartRepo)
	// checkoutHandler := checkout.NewHandler(checkoutService)
	// checkout.RegisterRoutes(checkoutHandler, apiV1.Group("/checkout"), app.db.GetPool())
}

func (app *App) Shutdown() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if app.server != nil {
		app.server.Shutdown(ctx)
	}

	if app.db != nil {
		app.db.Close()
	}

	app.logger.Debug("Graceful Shutdown")
}

func (app *App) Run() error {
	app.logger = log.NewConsoleLogger()
	err := godotenv.Load()
	if err != nil {
		app.logger.Warn(
			"failed to load .env file",
			log.Meta{
				"Error": err.Error(),
			},
		)
	}

	cfg := config.Load()

	app.db, err = database.ConnectDB(context.Background(), cfg)
	if err != nil {
		app.logger.Error(
			"database connection issue",
			log.Meta{
				"Error": err.Error(),
			},
		)
		return err
	}

	if err := app.db.Ping(context.Background()); err != nil {
		app.logger.Warn(
			"database is not live",
			log.Meta{
				"Error": err.Error(),
			},
		)
	}

	app.redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.RedisCfg.Host,
			cfg.RedisCfg.Port,
		),
	})

	if err := app.redisClient.Ping(context.Background()); err != nil {
		app.logger.Warn(
			"redis connection issue",
			log.Meta{
				"Error": err,
			},
		)
	}

	router := gin.Default()
	app.setupRoutes(cfg, router)

	app.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.SvPort),
		Handler: router,
	}

	app.logger.Info(
		"Server started",
		log.Meta{
			"URL":  fmt.Sprintf("http://localhost:%s", cfg.SvPort),
			"Port": cfg.SvPort,
		},
	)
	return app.server.ListenAndServe()
}
