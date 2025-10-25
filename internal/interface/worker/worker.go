package worker

import (
	"context"
	"go-jobqueue/internal/domain"
	"log"
	"time" // ⬅️ Must import time for sleeping
)

type JobProcessor interface {
	Execute(job *domain.Job) error
}

type Worker struct {
	repo      domain.JobRepository
	processor JobProcessor
}

func NewWorker(repo domain.JobRepository, processor JobProcessor) *Worker {
	return &Worker{repo: repo, processor: processor}
}

// Start now runs continuously, respects the context for graceful shutdown,
// and correctly handles an empty queue.
func (w *Worker) Start(ctx context.Context) {
	log.Println("👷 Worker started...")

	// 1. Worker must run in a continuous loop
	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Worker stopped gracefully")
			return
		default:
			// 2. Pull job from Redis queue
			job, err := w.repo.Dequeue()

			if err != nil {
				// Log and pause on infrastructure error (e.g., Redis down)
				log.Printf("⚠️ Worker Dequeue Error: %v\n", err)
				time.Sleep(time.Second * 5) // Longer sleep for serious errors
				continue
			}

			// 3. CRITICAL: Check for empty queue (Dequeue returns nil, nil)
			if job == nil {
				time.Sleep(time.Second) // Pause briefly to avoid a tight loop
				continue
			}

			// --- Job Found and Ready to Process ---
			log.Printf("📥 Dequeued job %s (%s)\n", job.ID, job.Type)

			// 4. Delegate to the Use Case Processor
			if err := w.processor.Execute(job); err != nil {
				// TODO: The Use Case should ideally handle the status update to JobFailed
				// If the processor failed, log it and continue to the next job.
				log.Printf("❌ Job %s FAILED: %v\n", job.ID, err)
			} else {
				log.Printf("✅ Job %s completed\n", job.ID)
			}
		}
	}
}
