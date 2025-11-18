# Quick Reference Guide

## Common Commands

### Development
```bash
# Run application
make run

# Build application
make build

# Run tests
make test

# Format code
make format

# Install dependencies
make install-deps

# Clean build artifacts
make clean
```

### Database
```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down
```

### Docker
```bash
# Build image
make docker-build

# Run container
make docker-run

# Start with compose
make docker-compose-up

# Stop compose
make docker-compose-down
```

## API Quick Reference

### Upload Video
```bash
curl -X POST http://localhost:8080/api/v1/videos/upload \
  -F "video=@video.mp4"
```

### List Videos
```bash
curl http://localhost:8080/api/v1/videos
```

### Get Video
```bash
curl http://localhost:8080/api/v1/videos/{id}
```

### Stream Video
```bash
curl http://localhost:8080/api/v1/videos/{id}/stream -o output.mp4
```

### Transcode
```bash
curl -X POST http://localhost:8080/api/v1/videos/{id}/transcode \
  -H "Content-Type: application/json" \
  -d '{"resolution":"720p","format":"mp4"}'
```

### Compress
```bash
curl -X POST http://localhost:8080/api/v1/videos/{id}/compress \
  -H "Content-Type: application/json" \
  -d '{"quality":23}'
```

### Thumbnails
```bash
curl -X POST http://localhost:8080/api/v1/videos/{id}/thumbnail \
  -H "Content-Type: application/json" \
  -d '{"count":4}'
```

### Job Status
```bash
curl http://localhost:8080/api/v1/jobs/{job_id}
```

### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

### Metrics
```bash
curl http://localhost:8080/api/v1/metrics
```

## Environment Variables

### Server
- `SERVER_HOST` - Server host (default: localhost)
- `SERVER_PORT` - Server port (default: 8080)
- `SERVER_READ_TIMEOUT` - Read timeout (default: 30s)
- `SERVER_WRITE_TIMEOUT` - Write timeout (default: 30s)

### Database
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `DB_MAX_CONNECTIONS` - Max connections (default: 25)

### Storage
- `STORAGE_UPLOAD_DIR` - Upload directory
- `STORAGE_PROCESSED_DIR` - Processed files directory
- `STORAGE_THUMBNAIL_DIR` - Thumbnails directory
- `STORAGE_MAX_FILE_SIZE` - Max file size in bytes (default: 500MB)

### Processing
- `PROCESSING_WORKER_COUNT` - Number of workers (default: 4)
- `PROCESSING_MAX_CONCURRENT_JOBS` - Max concurrent jobs (default: 10)
- `PROCESSING_TEMP_DIR` - Temp directory
- `PROCESSING_SUPPORTED_FORMATS` - Supported formats (comma-separated)

### FFmpeg
- `FFMPEG_PATH` - FFmpeg binary path
- `FFPROBE_PATH` - FFprobe binary path
- `FFMPEG_DEFAULT_CODEC` - Default video codec (default: libx264)
- `FFMPEG_DEFAULT_AUDIO_CODEC` - Default audio codec (default: aac)

### Logging
- `LOG_LEVEL` - Log level (debug, info, warn, error)
- `LOG_FORMAT` - Log format (json, text)
- `LOG_OUTPUT` - Log output (stdout, file)

## Resolution Presets

- `1080p` → 1920x1080
- `720p` → 1280x720
- `480p` → 854x480
- `360p` → 640x360
- `240p` → 426x240

## Video Formats

Supported formats:
- mp4
- avi
- mov
- mkv
- webm
- flv
- wmv
- m4v
- 3gp

## CRF Quality Values

- 18 - Very high quality (large file)
- 23 - Default (balanced)
- 28 - Lower quality (small file)
- Range: 0-51 (lower is better)

## Job Statuses

- `PENDING` - Waiting in queue
- `PROCESSING` - Currently processing
- `COMPLETED` - Successfully completed
- `FAILED` - Processing failed

## File Locations

```
storage/
├── uploads/           # Original videos
├── processed/         # Processed videos
└── thumbnails/        # Video thumbnails
    └── {video_id}/    # Per-video thumbnails
```

## Database Tables

### videos
- id (UUID)
- original_filename
- file_path
- file_size
- mime_type
- duration
- resolution
- codec
- created_at
- updated_at

### processing_jobs
- id (UUID)
- video_id
- operation
- status
- progress
- output_path
- error_message
- parameters (JSON)
- started_at
- completed_at
- created_at

## Common Workflows

### Basic Upload & Transcode
```bash
# 1. Upload
UPLOAD_RESP=$(curl -s -X POST http://localhost:8080/api/v1/videos/upload \
  -F "video=@sample.mp4")
VIDEO_ID=$(echo $UPLOAD_RESP | jq -r '.data.video_id')

# 2. Transcode
JOB_RESP=$(curl -s -X POST "http://localhost:8080/api/v1/videos/$VIDEO_ID/transcode" \
  -H "Content-Type: application/json" \
  -d '{"resolution":"720p"}')
JOB_ID=$(echo $JOB_RESP | jq -r '.data.id')

# 3. Check status
curl "http://localhost:8080/api/v1/jobs/$JOB_ID"
```

### Multiple Operations
```bash
# Upload once, process multiple ways
VIDEO_ID="your-video-id"

# Transcode to 720p
curl -X POST "http://localhost:8080/api/v1/videos/$VIDEO_ID/transcode" \
  -H "Content-Type: application/json" \
  -d '{"resolution":"720p"}'

# Compress
curl -X POST "http://localhost:8080/api/v1/videos/$VIDEO_ID/compress" \
  -H "Content-Type: application/json" \
  -d '{"quality":23}'

# Generate thumbnails
curl -X POST "http://localhost:8080/api/v1/videos/$VIDEO_ID/thumbnail" \
  -H "Content-Type: application/json" \
  -d '{"count":6}'
```

## Troubleshooting

### Check logs
```bash
# View recent logs
tail -f /var/log/video-processor.log

# Search for errors
grep ERROR /var/log/video-processor.log
```

### Database connection
```bash
psql -h localhost -U postgres -d video_processing
```

### Storage usage
```bash
du -sh storage/*
```

### Active processes
```bash
ps aux | grep video-processor
```

### Port usage
```bash
lsof -i :8080
```

## Performance Tips

1. **Adjust worker count** based on CPU cores
   ```env
   PROCESSING_WORKER_COUNT=6  # For 8-core system
   ```

2. **Increase max file size** if needed
   ```env
   STORAGE_MAX_FILE_SIZE=1048576000  # 1GB
   ```

3. **Use SSD** for storage directories

4. **Monitor queue depth**
   ```bash
   curl http://localhost:8080/api/v1/metrics | jq .data.queue_depth
   ```

## Security Checklist

- [ ] Change default database password
- [ ] Limit file upload size
- [ ] Validate file types
- [ ] Sanitize filenames
- [ ] Use HTTPS in production
- [ ] Add authentication (future)
- [ ] Enable rate limiting (future)
- [ ] Regular security updates

## Useful Scripts

### Delete old processed files
```bash
find storage/processed -type f -mtime +30 -delete
```

### Backup database
```bash
pg_dump video_processing > backup_$(date +%Y%m%d).sql
```

### Monitor active jobs
```bash
watch -n 2 'curl -s http://localhost:8080/api/v1/metrics | jq .data.active_jobs'
```

---

For detailed documentation, see:
- `README.md` - Overview
- `docs/API.md` - Full API documentation
- `docs/ARCHITECTURE.md` - System architecture
- `docs/GETTING_STARTED.md` - Setup guide
