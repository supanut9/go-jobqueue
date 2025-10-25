package usecase

import (
	"fmt"
	"go-jobqueue/internal/domain"
)

// New Dependency: The Use Case layer defines the contract for an external service.
// This ensures the executor is decoupled from the SMTP library.
type EmailService interface {
	Send(to string, subject string, body string) error
}

type EmailJobExecutor struct {
	// The executor depends on the interface, not the concrete implementation.
	sender EmailService
}

// NewEmailJobExecutor is used for dependency injection.
func NewEmailJobExecutor(sender EmailService) *EmailJobExecutor {
	return &EmailJobExecutor{sender: sender}
}

// Execute contains the business logic for email jobs.
func (e *EmailJobExecutor) Execute(job *domain.Job) error {
	// 1. Validate and Parse Payload
	recipient, ok := job.Payload["recipient"].(string)
	if !ok {
		return fmt.Errorf("email job requires a valid 'recipient' in payload")
	}

	subject, ok := job.Payload["subject"].(string)
	if !ok {
		subject = "Job Queue Notification" // Default subject
	}

	// Assume body is required
	body, ok := job.Payload["body"].(string)
	if !ok {
		return fmt.Errorf("email job requires 'body' in payload")
	}

	// 2. Delegate to the injected EmailService (the "real" part)
	if err := e.sender.Send(recipient, subject, body); err != nil {
		return fmt.Errorf("failed to send email to %s: %w", recipient, err)
	}

	return nil
}
