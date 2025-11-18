package services

import (
	"context"
	"fmt"

	"github.com/dev-dhanushkumar/go-video-processor/internal/models"
	"github.com/dev-dhanushkumar/go-video-processor/internal/queue"
	"github.com/dev-dhanushkumar/go-video-processor/internal/repository"
	"github.com/dev-dhanushkumar/go-video-processor/internal/utils"
	"github.com/dev-dhanushkumar/go-video-processor/pkg/config"
)

type StorageService struct {
	videoRepo repository.VideoRepository
	config    *config.Config
	logger    *utils.Logger
}

func NewStorageService(
	videoRepo repository.VideoRepository,
	cfg *config.Config,
	logger *utils.Logger,
) *StorageService {
	return &StorageService{
		videoRepo: videoRepo,
		config:    cfg,
		logger:    logger,
	}
}

// GetStorageUsage calculates total storage used by videos
func (s *StorageService) GetStorageUsage(ctx context.Context) (int64, error) {
	uploadSize, err := utils.GetDirectorySize(s.config.Storage.UploadDir)
	if err != nil {
		return 0, fmt.Errorf("failed to get upload directory size: %w", err)
	}

	processedSize, err := utils.GetDirectorySize(s.config.Storage.ProcessedDir)
	if err != nil {
		return 0, fmt.Errorf("failed to get processed directory size: %w", err)
	}

	thumbnailSize, err := utils.GetDirectorySize(s.config.Storage.ThumbnailDir)
	if err != nil {
		return 0, fmt.Errorf("failed to get thumbnail directory size: %w", err)
	}

	return uploadSize + processedSize + thumbnailSize, nil
}

// GetMetrics returns system metrics
func (s *StorageService) GetMetrics(ctx context.Context, jobRepo repository.JobRepository, workerPool *queue.WorkerPool) (*models.MetricsResponse, error) {
	// Count total videos
	videos, totalVideos, err := s.videoRepo.List(ctx, 1, 0)
	if err != nil {
		return nil, err
	}
	_ = videos

	// Count jobs by status
	_, totalJobs, err := jobRepo.List(ctx, 1, 0)
	if err != nil {
		return nil, err
	}

	completedJobs, _ := jobRepo.CountByStatus(ctx, models.JobStatusCompleted)
	failedJobs, _ := jobRepo.CountByStatus(ctx, models.JobStatusFailed)
	activeJobs, _ := jobRepo.CountByStatus(ctx, models.JobStatusProcessing)

	// Get storage usage
	storageUsed, _ := s.GetStorageUsage(ctx)

	metrics := &models.MetricsResponse{
		TotalVideos:   totalVideos,
		TotalJobs:     totalJobs,
		ActiveJobs:    activeJobs,
		CompletedJobs: completedJobs,
		FailedJobs:    failedJobs,
		StorageUsed:   storageUsed,
		ActiveWorkers: workerPool.ActiveJobs(),
		QueueDepth:    workerPool.QueueDepth(),
	}

	return metrics, nil
}
