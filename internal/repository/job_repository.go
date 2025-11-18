package repository

import (
	"context"

	"github.com/dev-dhanushkumar/go-video-processor/internal/models"
	"gorm.io/gorm"
)

type JobRepository interface {
	Create(ctx context.Context, job *models.ProcessingJob) error
	GetByID(ctx context.Context, id string) (*models.ProcessingJob, error)
	GetByVideoID(ctx context.Context, videoID string) ([]models.ProcessingJob, error)
	List(ctx context.Context, limit, offset int) ([]models.ProcessingJob, int64, error)
	ListByStatus(ctx context.Context, status models.JobStatus) ([]models.ProcessingJob, error)
	Update(ctx context.Context, job *models.ProcessingJob) error
	UpdateStatus(ctx context.Context, id string, status models.JobStatus) error
	UpdateProgress(ctx context.Context, id string, progress int) error
	Delete(ctx context.Context, id string) error
	CountByStatus(ctx context.Context, status models.JobStatus) (int64, error)
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(ctx context.Context, job *models.ProcessingJob) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *jobRepository) GetByID(ctx context.Context, id string) (*models.ProcessingJob, error) {
	var job models.ProcessingJob
	err := r.db.WithContext(ctx).
		Preload("Video").
		Where("id = ?", id).
		First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) GetByVideoID(ctx context.Context, videoID string) ([]models.ProcessingJob, error) {
	var jobs []models.ProcessingJob
	err := r.db.WithContext(ctx).
		Where("video_id = ?", videoID).
		Order("created_at DESC").
		Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) List(ctx context.Context, limit, offset int) ([]models.ProcessingJob, int64, error) {
	var jobs []models.ProcessingJob
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&models.ProcessingJob{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.WithContext(ctx).
		Preload("Video").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error

	return jobs, total, err
}

func (r *jobRepository) ListByStatus(ctx context.Context, status models.JobStatus) ([]models.ProcessingJob, error) {
	var jobs []models.ProcessingJob
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at ASC").
		Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) Update(ctx context.Context, job *models.ProcessingJob) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *jobRepository) UpdateStatus(ctx context.Context, id string, status models.JobStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.ProcessingJob{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *jobRepository) UpdateProgress(ctx context.Context, id string, progress int) error {
	return r.db.WithContext(ctx).
		Model(&models.ProcessingJob{}).
		Where("id = ?", id).
		Update("progress", progress).Error
}

func (r *jobRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.ProcessingJob{}, "id = ?", id).Error
}

func (r *jobRepository) CountByStatus(ctx context.Context, status models.JobStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.ProcessingJob{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}
