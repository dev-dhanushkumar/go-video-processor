package models

// Standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Video upload response
type UploadResponse struct {
	VideoID string `json:"video_id"`
	JobID   string `json:"job_id,omitempty"`
	Message string `json:"message"`
}

// Job status response
type JobStatusResponse struct {
	JobID        string    `json:"job_id"`
	Status       JobStatus `json:"status"`
	Progress     int       `json:"progress"`
	ErrorMessage string    `json:"error_message,omitempty"`
	OutputPath   string    `json:"output_path,omitempty"`
}

// Video list response
type VideoListResponse struct {
	Videos []Video `json:"videos"`
	Total  int64   `json:"total"`
}

// Health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// Metrics response
type MetricsResponse struct {
	TotalVideos   int64 `json:"total_videos"`
	TotalJobs     int64 `json:"total_jobs"`
	ActiveJobs    int64 `json:"active_jobs"`
	CompletedJobs int64 `json:"completed_jobs"`
	FailedJobs    int64 `json:"failed_jobs"`
	StorageUsed   int64 `json:"storage_used_bytes"`
	ActiveWorkers int   `json:"active_workers"`
	QueueDepth    int   `json:"queue_depth"`
}

// Error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
