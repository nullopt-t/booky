package main

import (
	"booky-backend/internal/config"
	"booky-backend/internal/db"
	"booky-backend/internal/product"
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := db.NewDatabase(cfg.DBCfg)
	if err := db.Connect(context.Background()); err != nil {
		log.Fatalf("connecting database is failed : %W", err)
	}

	router := gin.Default()

	repo := product.NewPostgresRepo(db)
	service := product.NewService(repo)
	handler := product.NewHandler(service)
	handler.RegisterRoutes(router)

	router.Run(":8080")
}
