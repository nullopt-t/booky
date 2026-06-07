package app

import (
	// "booky-backend/internal/cart"
	"booky-backend/internal/http/swagger"
	"booky-backend/internal/middleware"

	// "booky-backend/internal/inventory"
	// "booky-backend/internal/product"
	"booky-backend/internal/user"

	// "booky-backend/internal/checkout"
	// "booky-backend/internal/order"
	"booky-backend/pkg/config"
	"booky-backend/pkg/database"
	"booky-backend/pkg/log"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	// database
	db *database.DB
}

func (app *App) setupRoutes(config *config.Config, router *gin.Engine) {
	// setup middlewares
	router.Use(middleware.ErrorHandler(app.logger))

	v1 := router.Group("/api/v1")
	swagger.SetUpDocs(v1)

	txRunner := database.NewTxRunner(app.db)

	// user
	userRepo := user.NewPostgresRepository()
	userService := user.NewService(txRunner, userRepo, app.logger)
	userHandler := user.NewHandler(userService, config)
	userRouter := user.NewRouter(userHandler, config)
	userRouter.MapRoutes(v1)

	// inventory
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
	cfg := config.Load()

	app.logger = log.NewConsoleLogger()

	var err error
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
