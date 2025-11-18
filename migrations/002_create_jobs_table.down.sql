-- Drop indexes (if they exist separately)
DROP INDEX IF EXISTS idx_jobs_created_at;
DROP INDEX IF EXISTS idx_jobs_status;
DROP INDEX IF EXISTS idx_jobs_video_id;

-- Drop processing_jobs table
DROP TABLE IF EXISTS processing_jobs;
