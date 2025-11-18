-- Drop trigger
DROP TRIGGER IF EXISTS update_videos_updated_at ON videos;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop videos table
DROP TABLE IF EXISTS videos;
