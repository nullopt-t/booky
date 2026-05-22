package app

import (
	"booky-backend/internal/cart"
	"booky-backend/internal/http/swagger"
	"booky-backend/internal/inventory"
	"booky-backend/internal/product"
	"booky-backend/pkg/logger"

	// "booky-backend/internal/checkout"
	"booky-backend/internal/order"
	"booky-backend/pkg/config"
	"booky-backend/pkg/database"
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

	// database
	db *database.DB
}

func (app *App) initHandlers(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	swagger.SetUpDocs(v1)

	txRunner := database.NewTxRunner(app.db)

	// inventory
	inventoryRepo := inventory.NewPostgresRepository()
	inventoryService := inventory.NewService(txRunner, inventoryRepo)
	inventoryHandler := inventory.NewInventoryHandler(inventoryService)
	inventory.MapRoutes(v1, inventoryHandler)

	// product
	productRepo := product.NewPostgresRepository()
	productService := product.NewService(txRunner, productRepo, inventoryRepo)
	productHandler := product.NewHandler(productService)
	product.MapRoutes(v1, productHandler)

	// cart
	cartRepo := cart.NewPostgresRepository()
	cartService := cart.NewService(txRunner, cartRepo, productRepo)
	cartHandler := cart.NewHandler(cartService)
	cart.MapRoutes(v1, cartHandler)

	// order
	orderRepo := order.NewPostgresRepository()
	orderService := order.NewService(txRunner, orderRepo)
	orderHandler := order.NewHandler(orderService)
	order.MapRoutes(v1, orderHandler)

	// // checkout
	// checkoutService := checkout.NewService(app.db.GetPool(), orderRepo, cartRepo)
	// checkoutHandler := checkout.NewHandler(checkoutService)
	// checkout.RegisterRoutes(checkoutHandler, apiV1.Group("/checkout"), app.db.GetPool())
}

func (app *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if app.server != nil {
		app.server.Shutdown(ctx)
	}

	if app.db != nil {
		app.db.Close()
	}

	logger.Log(logger.DEBUG, "Graceful Shutdown")
}

func (app *App) Run() error {
	cfg := config.Load()

	var err error
	app.db, err = database.ConnectDB(context.Background(), cfg)
	if err != nil {
		return err
	}

	router := gin.Default()
	app.initHandlers(router)

	app.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.SvPort),
		Handler: router,
	}

	logger.Log(logger.DEBUG, fmt.Sprintf("Server started on port %s", cfg.SvPort))
	return app.server.ListenAndServe()
}
