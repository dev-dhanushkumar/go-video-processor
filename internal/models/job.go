package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "PENDING"
	JobStatusProcessing JobStatus = "PROCESSING"
	JobStatusCompleted  JobStatus = "COMPLETED"
	JobStatusFailed     JobStatus = "FAILED"
)

type ProcessingJob struct {
	ID           string     `gorm:"type:varchar(36);primary_key" json:"id"`
	VideoID      string     `gorm:"type:varchar(36);not null;index" json:"video_id"`
	Operation    string     `gorm:"type:varchar(50);not null" json:"operation"`
	Status       JobStatus  `gorm:"type:varchar(20);default:'PENDING'" json:"status"`
	Progress     int        `gorm:"default:0" json:"progress"`
	OutputPath   string     `gorm:"type:varchar(512)" json:"output_path,omitempty"`
	ErrorMessage string     `gorm:"type:text" json:"error_message,omitempty"`
	Parameters   string     `gorm:"type:text" json:"parameters,omitempty"` // JSON string for operation params
	StartedAt    *time.Time `json:"started_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`

	// Relationships
	Video Video `gorm:"foreignKey:VideoID;constraint:OnDelete:CASCADE" json:"video,omitempty"`
}

// BeforeCreate hook to generate UUID
func (j *ProcessingJob) BeforeCreate(tx *gorm.DB) error {
	if j.ID == "" {
		j.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (ProcessingJob) TableName() string {
	return "processing_jobs"
}
