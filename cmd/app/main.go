package main

import (
	"context"
	"go-jobqueue/internal/infrastructure"
	"go-jobqueue/internal/interface/http"
	"go-jobqueue/internal/interface/worker"
	"go-jobqueue/internal/usecase"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Setup Redis connection
	redisClient := infrastructure.NewRedisClient()

	// Create repository and use case
	repo := infrastructure.NewRedisJobRepository(redisClient)
	enqueueUC := usecase.NewEnqueueJobUsecase(repo)

	// Start worker in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker.NewWorker(repo).Start(ctx)

	app := fiber.New()
	handler := http.NewJobHandler(enqueueUC)

	// Register routes
	app.Post("/jobs", handler.Enqueue)

	// Start server
	log.Println("🚀 Server running on :8080")
	app.Listen(":8080")

	// Give worker time to finish gracefully
	time.Sleep(time.Second)
}
