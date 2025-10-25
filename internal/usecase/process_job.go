package usecase

import (
	"go-jobqueue/internal/domain"
	"log"
	// Assuming Executor and JobDispatcher types are defined here or imported
)

type ProcessJobUsecase struct {
	repo       domain.JobRepository
	dispatcher *JobDispatcher // <-- NEW FIELD
}

// FIX: NewProcessJobUsecase now accepts the repository AND the dispatcher.
func NewProcessJobUsecase(repo domain.JobRepository, dispatcher *JobDispatcher) *ProcessJobUsecase {
	return &ProcessJobUsecase{repo: repo, dispatcher: dispatcher}
}

func (uc *ProcessJobUsecase) Execute(job *domain.Job) error {
	// 1. Mark as processing (using the injected repo)
	job.Status = domain.JobProcessing
	if err := uc.repo.UpdateStatus(job.ID, job.Status); err != nil {
		return err
	}

	// 2. Delegate the actual work to the dispatcher
	if err := uc.dispatcher.Dispatch(job); err != nil {
		// Log the dispatch error and mark job as failed
		log.Printf("Job %s dispatch failed: %v", job.ID, err)
		job.Status = domain.JobFailed
		_ = uc.repo.UpdateStatus(job.ID, job.Status) // Attempt to update status despite original error
		return err
	}

	// 3. Mark as success (if dispatch succeeded)
	log.Printf("⚙️ Job %s work completed.\n", job.ID)
	job.Status = domain.JobSuccess
	return uc.repo.UpdateStatus(job.ID, job.Status)
}
