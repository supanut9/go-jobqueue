package worker

import (
	"context"
	"log"
	"time"

	"go-jobqueue/internal/domain"
)

type Worker struct {
	repo domain.JobRepository
}

func NewWorker(repo domain.JobRepository) *Worker {
	return &Worker{repo: repo}
}

func (w *Worker) Start(ctx context.Context) {
	log.Println("👷 Worker started...")

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Worker stopped")
			return
		default:
			// 1. Pull job from Redis queue
			job, err := w.repo.Dequeue()
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			if job == nil {
				time.Sleep(time.Second)
				continue
			}

			// 2. Mark as processing
			job.Status = domain.JobProcessing
			_ = w.repo.UpdateStatus(job.ID, job.Status)

			// 3. Simulate job work
			log.Printf("⚙️ Processing job %s (%s)\n", job.ID, job.Type)
			time.Sleep(3 * time.Second) // simulate work

			// 4. Mark as success
			job.Status = domain.JobSuccess
			_ = w.repo.UpdateStatus(job.ID, job.Status)

			log.Printf("✅ Job %s completed\n", job.ID)
		}
	}
}
