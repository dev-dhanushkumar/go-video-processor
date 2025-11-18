package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Video struct {
	ID               string    `gorm:"type:varchar(36);primary_key" json:"id"`
	OriginalFilename string    `gorm:"type:varchar(255);not null" json:"original_filename"`
	FilePath         string    `gorm:"type:varchar(512);not null" json:"file_path"`
	FileSize         int64     `gorm:"not null" json:"file_size"`
	MimeType         string    `gorm:"type:varchar(100)" json:"mime_type"`
	Duration         float64   `json:"duration"`
	Resolution       string    `gorm:"type:varchar(20)" json:"resolution"`
	Codec            string    `gorm:"type:varchar(50)" json:"codec"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// BeforeCreate hook to generate UUID
func (v *Video) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (Video) TableName() string {
	return "videos"
}
