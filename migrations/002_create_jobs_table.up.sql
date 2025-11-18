-- Create processing_jobs table
CREATE TABLE IF NOT EXISTS processing_jobs (
    id VARCHAR(36) PRIMARY KEY,
    video_id VARCHAR(36) NOT NULL,
    operation VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING',
    progress INT DEFAULT 0,
    output_path VARCHAR(512),
    error_message TEXT,
    parameters TEXT,
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_jobs_video_id ON processing_jobs(video_id);
CREATE INDEX IF NOT EXISTS idx_jobs_status ON processing_jobs(status);
CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON processing_jobs(created_at);
