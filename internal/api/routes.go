package api

import (
	"github.com/dev-dhanushkumar/go-video-processor/internal/api/handlers"
	"github.com/dev-dhanushkumar/go-video-processor/internal/api/middleware"
	"github.com/dev-dhanushkumar/go-video-processor/internal/utils"
	"github.com/gin-gonic/gin"
)

type Router struct {
	videoHandler  *handlers.VideoHandler
	jobHandler    *handlers.JobHandler
	healthHandler *handlers.HealthHandler
	logger        *utils.Logger
}

func NewRouter(
	videoHandler *handlers.VideoHandler,
	jobHandler *handlers.JobHandler,
	healthHandler *handlers.HealthHandler,
	logger *utils.Logger,
) *Router {
	return &Router{
		videoHandler:  videoHandler,
		jobHandler:    jobHandler,
		healthHandler: healthHandler,
		logger:        logger,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware(r.logger))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", r.healthHandler.HealthCheck)
		v1.GET("/metrics", r.healthHandler.GetMetrics)

		// Video management
		videos := v1.Group("/videos")
		{
			videos.POST("/upload", r.videoHandler.UploadVideo)
			videos.GET("", r.videoHandler.ListVideos)
			videos.GET("/:id", r.videoHandler.GetVideo)
			videos.GET("/:id/stream", r.videoHandler.StreamVideo)
			videos.DELETE("/:id", r.videoHandler.DeleteVideo)

			// Processing operations
			videos.POST("/:id/transcode", r.jobHandler.TranscodeVideo)
			videos.POST("/:id/compress", r.jobHandler.CompressVideo)
			videos.POST("/:id/thumbnail", r.jobHandler.GenerateThumbnail)
		}

		// Job management
		jobs := v1.Group("/jobs")
		{
			jobs.GET("/:id", r.jobHandler.GetJobStatus)
		}
	}

	return router
}
