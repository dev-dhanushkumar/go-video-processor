package handlers

import (
	"net/http"

	"github.com/dev-dhanushkumar/go-video-processor/internal/models"
	"github.com/dev-dhanushkumar/go-video-processor/internal/services"
	"github.com/dev-dhanushkumar/go-video-processor/internal/utils"
	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	processingService *services.ProcessingService
	logger            *utils.Logger
}

func NewJobHandler(processingService *services.ProcessingService, logger *utils.Logger) *JobHandler {
	return &JobHandler{
		processingService: processingService,
		logger:            logger,
	}
}

// TranscodeVideo creates a transcode job
func (h *JobHandler) TranscodeVideo(c *gin.Context) {
	videoID := c.Param("id")

	var params services.TranscodeParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid parameters: " + err.Error(),
		})
		return
	}

	job, err := h.processingService.TranscodeVideo(c.Request.Context(), videoID, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Transcode job created",
		Data:    job,
	})
}

// CompressVideo creates a compress job
func (h *JobHandler) CompressVideo(c *gin.Context) {
	videoID := c.Param("id")

	var params services.CompressParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid parameters: " + err.Error(),
		})
		return
	}

	job, err := h.processingService.CompressVideo(c.Request.Context(), videoID, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Compress job created",
		Data:    job,
	})
}

// GenerateThumbnail creates a thumbnail job
func (h *JobHandler) GenerateThumbnail(c *gin.Context) {
	videoID := c.Param("id")

	var params services.ThumbnailParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid parameters: " + err.Error(),
		})
		return
	}

	job, err := h.processingService.GenerateThumbnail(c.Request.Context(), videoID, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Thumbnail job created",
		Data:    job,
	})
}

// GetJobStatus retrieves job status
func (h *JobHandler) GetJobStatus(c *gin.Context) {
	jobID := c.Param("id")

	job, err := h.processingService.GetJobStatus(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Job not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: models.JobStatusResponse{
			JobID:        job.ID,
			Status:       job.Status,
			Progress:     job.Progress,
			ErrorMessage: job.ErrorMessage,
			OutputPath:   job.OutputPath,
		},
	})
}
