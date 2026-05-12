package app

import (
	"booky-backend/internal/cart"
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

func (app *App) initHandlers(router *gin.Engine) *gin.Engine {
	apiV1 := router.Group("/api/v1")

	txRunner := db.NewTxRunner(app.db)

	// cart
	cartRepo := cart.NewPostgresRepository()
	cartService := cart.NewService(cartRepo, txRunner)
	cartHandler := cart.NewHandler(cartService)
	cart.RegisterRoutes(apiV1.Group("/cart"), cartHandler)

	// orders
	orderRepo := order.NewPostgresRepo()
	orderService := order.NewService(txRunner, orderRepo)
	orderHandler := order.NewHandler(orderService)
	order.RegisterRoutes(apiV1.Group("/orders"), orderHandler)

	// // checkout
	// checkoutService := checkout.NewService(app.db.GetPool(), orderRepo, cartRepo)
	// checkoutHandler := checkout.NewHandler(checkoutService)
	// checkout.RegisterRoutes(checkoutHandler, apiV1.Group("/checkout"), app.db.GetPool())
	return router
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

	fmt.Println("Graceful Shutdown")
}

func (app *App) Run() error {
	cfg := config.Load()

	var err error
	app.db, err = db.ConnectDB(context.Background(), cfg)
	if err != nil {
		return err
	}

	router := app.initHandlers(gin.Default())

	app.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.SvPort),
		Handler: router,
	}

	return app.server.ListenAndServe()
}
