package app

import (
	"booky-backend/internal/cart"
	"booky-backend/internal/inventory"
	"booky-backend/internal/product"
	"booky-backend/internal/shared"

	// "booky-backend/internal/checkout"
	"booky-backend/internal/config"
	"booky-backend/internal/db"
	"booky-backend/internal/order"
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
	db *db.DB
}

func (app *App) initHandlers(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	txRunner := db.NewTxRunner(app.db)

	// inventory
	inventoryRepo := inventory.NewPostgresRepository()

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

	shared.Log(shared.DEBUG, "Graceful Shutdown")
}

func (app *App) Run() error {
	cfg := config.Load()

	var err error
	app.db, err = db.ConnectDB(context.Background(), cfg)
	if err != nil {
		return err
	}

	router := gin.Default()
	app.initHandlers(router)

	app.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.SvPort),
		Handler: router,
	}

	shared.Log(shared.DEBUG, fmt.Sprintf("Server started on port %s", cfg.SvPort))
	return app.server.ListenAndServe()
}
