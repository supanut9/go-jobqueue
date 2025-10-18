package domain

type JobRepository interface {
	Enqueue(job *Job) error
	Dequeue() (*Job, error)
	UpdateStatus(id string, status JobStatus) error
}
