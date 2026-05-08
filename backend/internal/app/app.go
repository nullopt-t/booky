package app

import (
	"booky-backend/internal/config"
	"booky-backend/internal/db"
	"booky-backend/internal/order"
	"booky-backend/internal/payment"
	"booky-backend/internal/product"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type application interface {
	Run() error
	Close()
}

type App struct {
	// http server
	server *http.Server

	// database
	db *db.DB
}

func (app *App) initHandlers(router *gin.Engine) *gin.Engine {
	repo := product.NewPostgresRepo(app.db)
	service := product.NewService(repo)
	handler := product.NewHandler(service)
	handler.RegisterRoutes(router)

	orepo := order.NewPostgresRepo(app.db)
	oservice := order.NewService(orepo)
	ohandler := order.NewHandler(oservice)
	ohandler.RegisterRoutes(router)

	prepo := payment.NewPostgresRepo(app.db)
	pservice := payment.NewService(prepo, orepo)
	phandler := payment.NewHandler(pservice)
	phandler.RegisterRoutes(router)
	return router
}

func (app *App) Run() error {
	cfg := config.Load()
	db := db.NewDatabase(cfg.DBCfg)
	err := db.Connect(context.Background())
	if err != nil {
		return err
	}
	app.db = db

	router := app.initHandlers(gin.Default())

	app.server = &http.Server{
		Addr:    cfg.SvPort,
		Handler: router,
	}

	return app.server.ListenAndServe()
}

func (app *App) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if app.server != nil {
		app.server.Shutdown(ctx)
	}

	if app.db != nil {
		app.db.Close()
	}
}
