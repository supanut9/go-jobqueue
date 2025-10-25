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
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on OS environment variables.")
	}

	cfg := infrastructure.LoadConfig()

	// Setup Redis connection
	redisClient := infrastructure.NewRedisClient(cfg.RedisAddr)
	emailGateway := infrastructure.NewSMTPGateway(cfg.SMTP)

	// Create repository and use case
	redisJobRepository := infrastructure.NewRedisJobRepository(redisClient)
	emailExecutor := usecase.NewEmailJobExecutor(emailGateway)
	enqueueJobUsecase := usecase.NewEnqueueJobUsecase(redisJobRepository)
	dispatcher := usecase.NewJobDispatcher(map[string]usecase.Executor{
		"send-email": emailExecutor,
	})
	processJobUsecase := usecase.NewProcessJobUsecase(redisJobRepository, dispatcher)

	// Start worker in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker.NewWorker(redisJobRepository, processJobUsecase).Start(ctx)

	app := fiber.New()
	handler := http.NewJobHandler(enqueueJobUsecase)

	// Register routes
	app.Post("/jobs", handler.Enqueue)

	// Start server
	log.Println("🚀 Server running on :8080")
	app.Listen(":8080")

	// Give worker time to finish gracefully
	time.Sleep(time.Second)
}
