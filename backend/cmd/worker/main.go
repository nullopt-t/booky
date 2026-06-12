package main

import (
	"booky-backend/internal/notifier"
	"booky-backend/pkg/config"
	"booky-backend/pkg/log"
	"booky-backend/pkg/mail"
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisCfg.Host, cfg.RedisCfg.Port),
	})

	fmt.Println("redis client created")

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		fmt.Println("redis ping failed:", err)
		return
	}
	fmt.Println("redis ping successful")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handle Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		cancel()
	}()

	jobQueue := notifier.NewRedisJobQueue(redisClient)
	logger := log.NewConsoleLogger()

	mailer := mail.NewMailer(&mail.Config{
		Host:     cfg.SMTPCfg.Host,
		Port:     cfg.SMTPCfg.Port,
		Username: cfg.SMTPCfg.Username,
		Password: cfg.SMTPCfg.Password,
	})

	renderer, err := notifier.NewRenderer()
	if err != nil {
		fmt.Println("renderer creation failed:", err)
		return
	}

	worker := notifier.NewNotifierWorker(
		jobQueue,
		logger,
		mailer,
		renderer,
	)

	worker.Start(ctx)
}
