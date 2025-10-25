package usecase

import (
	"fmt"
	"go-jobqueue/internal/domain"
)

type Executor interface{ Execute(job *domain.Job) error }

type JobDispatcher struct {
	executors map[string]Executor
}

// FIX: Change NewJobDispatcher to accept the map of executors.
func NewJobDispatcher(executors map[string]Executor) *JobDispatcher {
	return &JobDispatcher{
		// Use the map passed via dependency injection
		executors: executors,
	}
}

func (d *JobDispatcher) Dispatch(job *domain.Job) error {
	if executor, ok := d.executors[job.Type]; ok {
		return executor.Execute(job)
	}
	return fmt.Errorf("unknown job type: %s", job.Type)
}
