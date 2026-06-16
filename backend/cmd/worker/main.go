package main

import (
	"booky-backend/internal/notifier"
	"booky-backend/internal/shared/html"
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
	logger := log.NewConsoleLogger()

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisCfg.Host, cfg.RedisCfg.Port),
	})

	logger.Info("user client created")

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Error("redis ping failed", log.Meta{
			"Error": err,
		})
		return
	}
	logger.Info("redis ping successful")

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

	mailer := mail.NewMailer(&mail.Config{
		Host:     cfg.SMTPCfg.Host,
		Port:     cfg.SMTPCfg.Port,
		Username: cfg.SMTPCfg.Username,
		Password: cfg.SMTPCfg.Password,
	})

	renderer, err := html.NewRenderer()
	if err != nil {
		logger.Error(
			"renderer creation failed:", log.Meta{
				"Error": err,
			})
		return
	}

	emailMessageHandler := notifier.NewEmailHandler(renderer, mailer)

	dispatcher := notifier.NewMessageDispatcher()
	dispatcher.Register(notifier.MessageTypeEmailOTP, emailMessageHandler.SendEmailOTP)

	worker := notifier.NewNotifierWorker(
		jobQueue,
		dispatcher,
		logger,
		cfg.ClientCfg,
	)

	worker.Start(ctx)
}
