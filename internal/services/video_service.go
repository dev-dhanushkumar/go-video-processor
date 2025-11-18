package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/dev-dhanushkumar/go-video-processor/internal/models"
	"github.com/dev-dhanushkumar/go-video-processor/internal/processor"
	"github.com/dev-dhanushkumar/go-video-processor/internal/repository"
	"github.com/dev-dhanushkumar/go-video-processor/internal/utils"
	"github.com/dev-dhanushkumar/go-video-processor/pkg/config"
	"github.com/google/uuid"
)

type VideoService struct {
	videoRepo repository.VideoRepository
	processor *processor.FFmpegProcessor
	config    *config.Config
	logger    *utils.Logger
}

func NewVideoService(
	videoRepo repository.VideoRepository,
	ffmpegProc *processor.FFmpegProcessor,
	cfg *config.Config,
	logger *utils.Logger,
) *VideoService {
	return &VideoService{
		videoRepo: videoRepo,
		processor: ffmpegProc,
		config:    cfg,
		logger:    logger,
	}
}

// UploadVideo handles video upload and metadata extraction
func (s *VideoService) UploadVideo(ctx context.Context, file *multipart.FileHeader) (*models.Video, error) {
	// Validate video file
	if err := utils.ValidateVideoFile(file, s.config.Storage.MaxFileSize); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate unique ID for video
	videoID := uuid.New().String()

	// Sanitize filename
	sanitizedFilename := utils.SanitizeFilename(file.Filename)
	ext := utils.GetFileExtension(sanitizedFilename)

	// Create file path
	filename := fmt.Sprintf("%s.%s", videoID, ext)
	filePath := filepath.Join(s.config.Storage.UploadDir, filename)

	// Save uploaded file
	if err := utils.SaveUploadedFile(file, filePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Extract metadata using FFmpeg
	metadata, err := s.processor.GetMetadata(ctx, filePath)
	if err != nil {
		// Clean up uploaded file on error
		utils.DeleteFile(filePath)
		return nil, fmt.Errorf("failed to extract metadata: %w", err)
	}

	// Create video record
	video := &models.Video{
		ID:               videoID,
		OriginalFilename: file.Filename,
		FilePath:         filePath,
		FileSize:         file.Size,
		MimeType:         file.Header.Get("Content-Type"),
		Duration:         metadata.Duration,
		Resolution:       metadata.Resolution,
		Codec:            metadata.Codec,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Save to database
	if err := s.videoRepo.Create(ctx, video); err != nil {
		// Clean up uploaded file on error
		utils.DeleteFile(filePath)
		return nil, fmt.Errorf("failed to save video record: %w", err)
	}

	s.logger.InfoWithFields("Video uploaded successfully", map[string]interface{}{
		"video_id":   videoID,
		"filename":   file.Filename,
		"size":       utils.FormatFileSize(file.Size),
		"duration":   metadata.Duration,
		"resolution": metadata.Resolution,
	})

	return video, nil
}

// GetVideo retrieves a video by ID
func (s *VideoService) GetVideo(ctx context.Context, id string) (*models.Video, error) {
	return s.videoRepo.GetByID(ctx, id)
}

// ListVideos retrieves a paginated list of videos
func (s *VideoService) ListVideos(ctx context.Context, limit, offset int) ([]models.Video, int64, error) {
	return s.videoRepo.List(ctx, limit, offset)
}

// DeleteVideo deletes a video and its associated file
func (s *VideoService) DeleteVideo(ctx context.Context, id string) error {
	// Get video record
	video, err := s.videoRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("video not found: %w", err)
	}

	// Delete file from storage
	if err := utils.DeleteFile(video.FilePath); err != nil {
		s.logger.WarnWithFields("Failed to delete video file", map[string]interface{}{
			"video_id": id,
			"path":     video.FilePath,
			"error":    err.Error(),
		})
	}

	// Delete from database
	if err := s.videoRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete video record: %w", err)
	}

	s.logger.InfoWithFields("Video deleted successfully", map[string]interface{}{
		"video_id": id,
	})

	return nil
}

// GetVideoPath returns the file path of a video
func (s *VideoService) GetVideoPath(ctx context.Context, id string) (string, error) {
	video, err := s.videoRepo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	return video.FilePath, nil
}
