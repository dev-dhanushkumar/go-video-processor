package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dev-dhanushkumar/go-video-processor/internal/api"
	"github.com/dev-dhanushkumar/go-video-processor/internal/api/handlers"
	"github.com/dev-dhanushkumar/go-video-processor/internal/processor"
	"github.com/dev-dhanushkumar/go-video-processor/internal/queue"
	"github.com/dev-dhanushkumar/go-video-processor/internal/repository"
	"github.com/dev-dhanushkumar/go-video-processor/internal/services"
	"github.com/dev-dhanushkumar/go-video-processor/internal/utils"
	"github.com/dev-dhanushkumar/go-video-processor/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize logger
	logger := utils.NewLogger(cfg.Logging.Level, cfg.Logging.Format)
	logger.Info("Starting Video Processing System...")

	// Ensure storage directories exist
	if err := utils.EnsureDirectory(cfg.Storage.UploadDir); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to create upload directory: %v", err))
	}
	if err := utils.EnsureDirectory(cfg.Storage.ProcessedDir); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to create processed directory: %v", err))
	}
	if err := utils.EnsureDirectory(cfg.Storage.ThumbnailDir); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to create thumbnail directory: %v", err))
	}
	if err := utils.EnsureDirectory(cfg.Processing.TempDir); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to create temp directory: %v", err))
	}

	// Initialize database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Database,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Auto migrate database schema
	logger.Info("Running database migrations...")
	// Note: Migrations should be run using SQL migration files
	// For now, we'll skip auto-migration and rely on manual migrations
	logger.Info("Database migrations skipped - use migrations/*.sql files")

	// Initialize FFmpeg processor
	ffmpegProc := processor.NewFFmpegProcessor(cfg.FFmpeg.Path, cfg.FFmpeg.FFprobePath)

	// Validate FFmpeg installation
	if err := ffmpegProc.ValidateFFmpeg(); err != nil {
		logger.Fatal(fmt.Sprintf("FFmpeg validation failed: %v", err))
	}
	logger.Info("FFmpeg validated successfully")

	// Initialize repositories
	videoRepo := repository.NewVideoRepository(db)
	jobRepo := repository.NewJobRepository(db)

	// Initialize worker pool
	workerPool := queue.NewWorkerPool(
		cfg.Processing.WorkerCount,
		cfg.Processing.MaxConcurrentJobs,
		logger,
	)
	workerPool.Start()
	logger.Info(fmt.Sprintf("Worker pool started with %d workers", cfg.Processing.WorkerCount))

	// Initialize services
	videoService := services.NewVideoService(videoRepo, ffmpegProc, cfg, logger)
	processingService := services.NewProcessingService(jobRepo, videoRepo, ffmpegProc, workerPool, cfg, logger)
	storageService := services.NewStorageService(videoRepo, cfg, logger)

	// Initialize handlers
	videoHandler := handlers.NewVideoHandler(videoService, logger)
	jobHandler := handlers.NewJobHandler(processingService, logger)
	healthHandler := handlers.NewHealthHandler(storageService, jobRepo, workerPool)

	// Setup router
	router := api.NewRouter(videoHandler, jobHandler, healthHandler, logger)
	engine := router.SetupRoutes()

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info(fmt.Sprintf("Server listening on %s", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(fmt.Sprintf("Server failed to start: %v", err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop worker pool
	workerPool.Stop()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	logger.Info("Server stopped")
}
