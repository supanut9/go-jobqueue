package domain

import "time"

type JobStatus string

const (
	JobPending    JobStatus = "pending"
	JobProcessing JobStatus = "processing"
	JobSuccess    JobStatus = "success"
	JobFailed     JobStatus = "failed"
)

type Job struct {
	ID        string
	Type      string
	Payload   map[string]interface{}
	Status    JobStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
