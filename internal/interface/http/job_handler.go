package http

import (
	"go-jobqueue/internal/domain"
	"go-jobqueue/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type JobHandler struct {
	usecase *usecase.EnqueueJobUsecase
}

func NewJobHandler(usecase *usecase.EnqueueJobUsecase) *JobHandler {
	return &JobHandler{usecase: usecase}
}

func (h *JobHandler) Enqueue(c *fiber.Ctx) error {
	var job domain.Job
	if err := c.BodyParser(&job); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.usecase.Execute(&job); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "queued", "id": job.ID})
}
