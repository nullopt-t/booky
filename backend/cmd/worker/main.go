package main

import (
	"booky-backend/internal/notifier"
	"booky-backend/pkg/config"
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisCfg.Addr,
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
	worker := notifier.NewNotifierWorker(jobQueue)

	worker.Start(ctx)
}
