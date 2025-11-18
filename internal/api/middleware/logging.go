package middleware

import (
	"time"

	"github.com/dev-dhanushkumar/go-video-processor/internal/utils"
	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(logger *utils.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		logger.InfoWithFields("Request processed", map[string]interface{}{
			"method":   method,
			"path":     path,
			"status":   statusCode,
			"duration": duration.String(),
			"ip":       c.ClientIP(),
		})
	}
}
