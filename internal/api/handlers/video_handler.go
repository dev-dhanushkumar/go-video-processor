package handlers

import (
	"net/http"
	"strconv"

	"github.com/dev-dhanushkumar/go-video-processor/internal/models"
	"github.com/dev-dhanushkumar/go-video-processor/internal/services"
	"github.com/dev-dhanushkumar/go-video-processor/internal/utils"
	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videoService *services.VideoService
	logger       *utils.Logger
}

func NewVideoHandler(videoService *services.VideoService, logger *utils.Logger) *VideoHandler {
	return &VideoHandler{
		videoService: videoService,
		logger:       logger,
	}
}

// UploadVideo handles video upload
func (h *VideoHandler) UploadVideo(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "No video file provided",
		})
		return
	}

	video, err := h.videoService.UploadVideo(c.Request.Context(), file)
	if err != nil {
		h.logger.Error("Failed to upload video: " + err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Video uploaded successfully",
		Data: models.UploadResponse{
			VideoID: video.ID,
			Message: "Video uploaded and metadata extracted",
		},
	})
}

// GetVideo retrieves video details
func (h *VideoHandler) GetVideo(c *gin.Context) {
	videoID := c.Param("id")

	video, err := h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Video not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    video,
	})
}

// ListVideos retrieves all videos
func (h *VideoHandler) ListVideos(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	videos, total, err := h.videoService.ListVideos(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to retrieve videos",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: models.VideoListResponse{
			Videos: videos,
			Total:  total,
		},
	})
}

// DeleteVideo deletes a video
func (h *VideoHandler) DeleteVideo(c *gin.Context) {
	videoID := c.Param("id")

	if err := h.videoService.DeleteVideo(c.Request.Context(), videoID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Video deleted successfully",
	})
}

// StreamVideo streams video content
func (h *VideoHandler) StreamVideo(c *gin.Context) {
	videoID := c.Param("id")

	videoPath, err := h.videoService.GetVideoPath(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Video not found",
		})
		return
	}

	// Stream the video file with range support
	c.File(videoPath)
}
