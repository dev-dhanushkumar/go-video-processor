package handlers

import (
	"net/http"
	"time"

	"github.com/dev-dhanushkumar/go-video-processor/internal/models"
	"github.com/dev-dhanushkumar/go-video-processor/internal/queue"
	"github.com/dev-dhanushkumar/go-video-processor/internal/repository"
	"github.com/dev-dhanushkumar/go-video-processor/internal/services"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	storageService *services.StorageService
	jobRepo        repository.JobRepository
	workerPool     *queue.WorkerPool
}

func NewHealthHandler(
	storageService *services.StorageService,
	jobRepo repository.JobRepository,
	workerPool *queue.WorkerPool,
) *HealthHandler {
	return &HealthHandler{
		storageService: storageService,
		jobRepo:        jobRepo,
		workerPool:     workerPool,
	}
}

// HealthCheck returns system health status
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: models.HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now().Format(time.RFC3339),
			Version:   "1.0.0",
		},
	})
}

// GetMetrics returns system metrics
func (h *HealthHandler) GetMetrics(c *gin.Context) {
	metrics, err := h.storageService.GetMetrics(c.Request.Context(), h.jobRepo, h.workerPool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to retrieve metrics",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    metrics,
	})
}
