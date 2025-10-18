package usecase

import (
	"go-jobqueue/internal/domain"
	"time"

	"github.com/google/uuid"
)

type EnqueueJobUsecase struct {
	repo domain.JobRepository
}

func NewEnqueueJobUsecase(repo domain.JobRepository) *EnqueueJobUsecase {
	return &EnqueueJobUsecase{repo: repo}
}

func (uc *EnqueueJobUsecase) Execute(job *domain.Job) error {
	job.ID = uuid.New().String()
	job.Status = domain.JobPending
	job.CreatedAt = time.Now()

	return uc.repo.Enqueue(job)
}
