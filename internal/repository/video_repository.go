package repository

import (
	"context"

	"github.com/dev-dhanushkumar/go-video-processor/internal/models"
	"gorm.io/gorm"
)

type VideoRepository interface {
	Create(ctx context.Context, video *models.Video) error
	GetByID(ctx context.Context, id string) (*models.Video, error)
	List(ctx context.Context, limit, offset int) ([]models.Video, int64, error)
	Update(ctx context.Context, video *models.Video) error
	Delete(ctx context.Context, id string) error
}

type videoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) VideoRepository {
	return &videoRepository{db: db}
}

func (r *videoRepository) Create(ctx context.Context, video *models.Video) error {
	return r.db.WithContext(ctx).Create(video).Error
}

func (r *videoRepository) GetByID(ctx context.Context, id string) (*models.Video, error) {
	var video models.Video
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) List(ctx context.Context, limit, offset int) ([]models.Video, int64, error) {
	var videos []models.Video
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&models.Video{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&videos).Error

	return videos, total, err
}

func (r *videoRepository) Update(ctx context.Context, video *models.Video) error {
	return r.db.WithContext(ctx).Save(video).Error
}

func (r *videoRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Video{}, "id = ?", id).Error
}
