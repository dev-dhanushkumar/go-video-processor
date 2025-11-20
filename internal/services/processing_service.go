package services

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/dev-dhanushkumar/go-video-processor/internal/models"
	"github.com/dev-dhanushkumar/go-video-processor/internal/processor"
	"github.com/dev-dhanushkumar/go-video-processor/internal/queue"
	"github.com/dev-dhanushkumar/go-video-processor/internal/repository"
	"github.com/dev-dhanushkumar/go-video-processor/internal/utils"
	"github.com/dev-dhanushkumar/go-video-processor/pkg/config"
	"github.com/google/uuid"
)

type ProcessingService struct {
	jobRepo    repository.JobRepository
	videoRepo  repository.VideoRepository
	processor  *processor.FFmpegProcessor
	workerPool *queue.WorkerPool
	config     *config.Config
	logger     *utils.Logger
}

func NewProcessingService(
	jobRepo repository.JobRepository,
	videoRepo repository.VideoRepository,
	ffmpegProc *processor.FFmpegProcessor,
	workerPool *queue.WorkerPool,
	cfg *config.Config,
	logger *utils.Logger,
) *ProcessingService {
	return &ProcessingService{
		jobRepo:    jobRepo,
		videoRepo:  videoRepo,
		processor:  ffmpegProc,
		workerPool: workerPool,
		config:     cfg,
		logger:     logger,
	}
}

// TranscodeParams represents parameters for transcoding
type TranscodeParams struct {
	Resolution string `json:"resolution"`
	Format     string `json:"format"`
	Quality    int    `json:"quality"`
}

// CompressParams represents parameters for compression
type CompressParams struct {
	Quality int `json:"quality"`
}

// ThumbnailParams represents parameters for thumbnail generation
type ThumbnailParams struct {
	Count int `json:"count"`
}

// TranscodeVideo creates a transcoding job
func (s *ProcessingService) TranscodeVideo(ctx context.Context, videoID string, params TranscodeParams) (*models.ProcessingJob, error) {
	// Verify video exists
	video, err := s.videoRepo.GetByID(ctx, videoID)
	if err != nil {
		return nil, fmt.Errorf("video not found: %w", err)
	}

	// Serialize parameters
	paramsJSON, _ := json.Marshal(params)

	// Create job record
	job := &models.ProcessingJob{
		ID:         uuid.New().String(),
		VideoID:    videoID,
		Operation:  "transcode",
		Status:     models.JobStatusPending,
		Parameters: string(paramsJSON),
		CreatedAt:  time.Now(),
	}

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	// Submit job to worker pool
	queueJob := &queue.Job{
		ID:        job.ID,
		VideoID:   videoID,
		Operation: "transcode",
		Handler: func(jobCtx context.Context) error {
			return s.executeTranscode(jobCtx, job.ID, video, params)
		},
	}

	if err := s.workerPool.Submit(queueJob); err != nil {
		s.jobRepo.UpdateStatus(ctx, job.ID, models.JobStatusFailed)
		return nil, fmt.Errorf("failed to submit job: %w", err)
	}

	s.logger.InfoWithFields("Transcode job created", map[string]interface{}{
		"job_id":   job.ID,
		"video_id": videoID,
		"params":   params,
	})

	return job, nil
}

// CompressVideo creates a compression job
func (s *ProcessingService) CompressVideo(ctx context.Context, videoID string, params CompressParams) (*models.ProcessingJob, error) {
	video, err := s.videoRepo.GetByID(ctx, videoID)
	if err != nil {
		return nil, fmt.Errorf("video not found: %w", err)
	}

	paramsJSON, _ := json.Marshal(params)

	job := &models.ProcessingJob{
		ID:         uuid.New().String(),
		VideoID:    videoID,
		Operation:  "compress",
		Status:     models.JobStatusPending,
		Parameters: string(paramsJSON),
		CreatedAt:  time.Now(),
	}

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	queueJob := &queue.Job{
		ID:        job.ID,
		VideoID:   videoID,
		Operation: "compress",
		Handler: func(jobCtx context.Context) error {
			return s.executeCompress(jobCtx, job.ID, video, params)
		},
	}

	if err := s.workerPool.Submit(queueJob); err != nil {
		s.jobRepo.UpdateStatus(ctx, job.ID, models.JobStatusFailed)
		return nil, fmt.Errorf("failed to submit job: %w", err)
	}

	s.logger.InfoWithFields("Compress job created", map[string]interface{}{
		"job_id":   job.ID,
		"video_id": videoID,
	})

	return job, nil
}

// GenerateThumbnail creates a thumbnail generation job
func (s *ProcessingService) GenerateThumbnail(ctx context.Context, videoID string, params ThumbnailParams) (*models.ProcessingJob, error) {
	video, err := s.videoRepo.GetByID(ctx, videoID)
	if err != nil {
		return nil, fmt.Errorf("video not found: %w", err)
	}

	paramsJSON, _ := json.Marshal(params)

	job := &models.ProcessingJob{
		ID:         uuid.New().String(),
		VideoID:    videoID,
		Operation:  "thumbnail",
		Status:     models.JobStatusPending,
		Parameters: string(paramsJSON),
		CreatedAt:  time.Now(),
	}

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	queueJob := &queue.Job{
		ID:        job.ID,
		VideoID:   videoID,
		Operation: "thumbnail",
		Handler: func(jobCtx context.Context) error {
			return s.executeThumbnail(jobCtx, job.ID, video, params)
		},
	}

	if err := s.workerPool.Submit(queueJob); err != nil {
		s.jobRepo.UpdateStatus(ctx, job.ID, models.JobStatusFailed)
		return nil, fmt.Errorf("failed to submit job: %w", err)
	}

	s.logger.InfoWithFields("Thumbnail job created", map[string]interface{}{
		"job_id":   job.ID,
		"video_id": videoID,
	})

	return job, nil
}

// GetJobStatus retrieves job status
func (s *ProcessingService) GetJobStatus(ctx context.Context, jobID string) (*models.ProcessingJob, error) {
	return s.jobRepo.GetByID(ctx, jobID)
}

// executeTranscode performs the actual transcoding
func (s *ProcessingService) executeTranscode(ctx context.Context, jobID string, video *models.Video, params TranscodeParams) error {
	now := time.Now()
	s.jobRepo.Update(ctx, &models.ProcessingJob{
		ID:        jobID,
		Status:    models.JobStatusProcessing,
		StartedAt: &now,
	})

	// Generate output filename
	ext := params.Format
	if ext == "" {
		ext = "mp4"
	}
	outputFilename := fmt.Sprintf("%s_transcoded_%s.%s", filepath.Base(video.FilePath[:len(video.FilePath)-len(filepath.Ext(video.FilePath))]), params.Resolution, ext)
	outputPath := filepath.Join(s.config.Storage.ProcessedDir, outputFilename)

	// Determine codecs based on output format
	videoCodec := s.config.FFmpeg.DefaultCodec
	audioCodec := s.config.FFmpeg.DefaultAudioCodec

	// WebM format requires VP8/VP9 for video and Vorbis/Opus for audio
	if ext == "webm" {
		videoCodec = "libvpx-vp9" // VP9 codec for WebM
		audioCodec = "libopus"    // Opus codec for WebM
	} else if ext == "mp4" {
		videoCodec = "libx264" // H.264 for MP4
		audioCodec = "aac"     // AAC for MP4
	} else if ext == "mkv" {
		videoCodec = "libx265" // H.265 for MKV
		audioCodec = "aac"
	}

	// Prepare transcode options
	opts := processor.TranscodeOptions{
		OutputPath: outputPath,
		VideoCodec: videoCodec,
		AudioCodec: audioCodec,
		Resolution: processor.GetResolutionString(params.Resolution),
		Preset:     "medium",
	}

	if params.Quality > 0 {
		opts.CRF = params.Quality
	} else {
		// Set default CRF based on format
		if ext == "webm" {
			opts.CRF = 31 // VP9 default quality (0-63 range)
		} else {
			opts.CRF = 23 // H.264/H.265 default quality (0-51 range)
		}
	}

	// Execute transcoding
	if err := s.processor.Transcode(ctx, video.FilePath, opts); err != nil {
		completed := time.Now()
		s.jobRepo.Update(ctx, &models.ProcessingJob{
			ID:           jobID,
			Status:       models.JobStatusFailed,
			ErrorMessage: err.Error(),
			CompletedAt:  &completed,
		})
		return err
	}

	// Update job as completed
	completed := time.Now()
	s.jobRepo.Update(ctx, &models.ProcessingJob{
		ID:          jobID,
		Status:      models.JobStatusCompleted,
		OutputPath:  outputPath,
		Progress:    100,
		CompletedAt: &completed,
	})

	return nil
}

// executeCompress performs the actual compression
func (s *ProcessingService) executeCompress(ctx context.Context, jobID string, video *models.Video, params CompressParams) error {
	now := time.Now()
	s.jobRepo.Update(ctx, &models.ProcessingJob{
		ID:        jobID,
		Status:    models.JobStatusProcessing,
		StartedAt: &now,
	})

	outputFilename := fmt.Sprintf("%s_compressed.mp4", filepath.Base(video.FilePath[:len(video.FilePath)-len(filepath.Ext(video.FilePath))]))
	outputPath := filepath.Join(s.config.Storage.ProcessedDir, outputFilename)

	quality := params.Quality
	if quality == 0 {
		quality = 23
	}

	opts := processor.CompressOptions{
		OutputPath: outputPath,
		Quality:    quality,
		Preset:     "medium",
	}

	if err := s.processor.Compress(ctx, video.FilePath, opts); err != nil {
		completed := time.Now()
		s.jobRepo.Update(ctx, &models.ProcessingJob{
			ID:           jobID,
			Status:       models.JobStatusFailed,
			ErrorMessage: err.Error(),
			CompletedAt:  &completed,
		})
		return err
	}

	completed := time.Now()
	s.jobRepo.Update(ctx, &models.ProcessingJob{
		ID:          jobID,
		Status:      models.JobStatusCompleted,
		OutputPath:  outputPath,
		Progress:    100,
		CompletedAt: &completed,
	})

	return nil
}

// executeThumbnail performs thumbnail generation
func (s *ProcessingService) executeThumbnail(ctx context.Context, jobID string, video *models.Video, params ThumbnailParams) error {
	now := time.Now()
	s.jobRepo.Update(ctx, &models.ProcessingJob{
		ID:        jobID,
		Status:    models.JobStatusProcessing,
		StartedAt: &now,
	})

	// Create thumbnail directory for this video
	thumbDir := filepath.Join(s.config.Storage.ThumbnailDir, video.ID)
	if err := utils.EnsureDirectory(thumbDir); err != nil {
		completed := time.Now()
		s.jobRepo.Update(ctx, &models.ProcessingJob{
			ID:           jobID,
			Status:       models.JobStatusFailed,
			ErrorMessage: err.Error(),
			CompletedAt:  &completed,
		})
		return err
	}

	count := params.Count
	if count == 0 {
		count = 4
	}

	thumbnails, err := s.processor.GenerateAutoThumbnails(ctx, video.FilePath, thumbDir, count)
	if err != nil {
		completed := time.Now()
		s.jobRepo.Update(ctx, &models.ProcessingJob{
			ID:           jobID,
			Status:       models.JobStatusFailed,
			ErrorMessage: err.Error(),
			CompletedAt:  &completed,
		})
		return err
	}

	completed := time.Now()
	s.jobRepo.Update(ctx, &models.ProcessingJob{
		ID:          jobID,
		Status:      models.JobStatusCompleted,
		OutputPath:  thumbDir,
		Progress:    100,
		CompletedAt: &completed,
	})

	s.logger.InfoWithFields("Thumbnails generated", map[string]interface{}{
		"job_id": jobID,
		"count":  len(thumbnails),
	})

	return nil
}
