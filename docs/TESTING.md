# API Testing Guide

## Overview
This document provides complete testing instructions for the Video Processing API.

## Prerequisites
1. Server running on `http://localhost:8080`
2. PostgreSQL database configured
3. FFmpeg installed
4. Sample video file for testing

## Quick Start
```bash
# Start the server
./bin/video-processor

# Or in background
./bin/video-processor > /tmp/server.log 2>&1 &
```

## API Endpoints

### 1. Health Check
```bash
curl http://localhost:8080/api/v1/health | jq .
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2025-11-18T09:44:33+05:30",
    "version": "1.0.0"
  }
}
```

### 2. Upload Video
```bash
curl -X POST http://localhost:8080/api/v1/videos/upload \
  -F "video=@/path/to/your/video.mp4"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Video uploaded successfully",
  "data": {
    "video_id": "92503788-0298-4e79-9212-898ee9e66cc6",
    "message": "Video uploaded and metadata extracted"
  }
}
```

**Note:** Save the `video_id` for subsequent operations.

### 3. List All Videos
```bash
curl http://localhost:8080/api/v1/videos | jq .
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "videos": [
      {
        "id": "92503788-0298-4e79-9212-898ee9e66cc6",
        "original_filename": "sample.mp4",
        "file_path": "storage/uploads/92503788-0298-4e79-9212-898ee9e66cc6.mp4",
        "file_size": 5819906,
        "mime_type": "application/octet-stream",
        "duration": 13.858866,
        "resolution": "718x1280",
        "codec": "h264",
        "created_at": "2025-11-18T09:38:01.984684Z",
        "updated_at": "2025-11-18T09:38:01.984684Z"
      }
    ],
    "total": 1
  }
}
```

### 4. Get Video Details
```bash
VIDEO_ID="92503788-0298-4e79-9212-898ee9e66cc6"  # Replace with your video ID
curl http://localhost:8080/api/v1/videos/$VIDEO_ID | jq .
```

### 5. Generate Thumbnails
```bash
VIDEO_ID="92503788-0298-4e79-9212-898ee9e66cc6"  # Replace with your video ID
curl -X POST http://localhost:8080/api/v1/videos/$VIDEO_ID/thumbnail \
  -H "Content-Type: application/json" \
  -d '{"count": 3}'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Thumbnail job created",
  "data": {
    "id": "95a546de-8de6-4e0f-87aa-7acef5ea56c7",
    "video_id": "92503788-0298-4e79-9212-898ee9e66cc6",
    "operation": "thumbnail",
    "status": "PENDING",
    "progress": 0,
    "parameters": "{\"count\":3}"
  }
}
```

### 6. Check Job Status
```bash
JOB_ID="95a546de-8de6-4e0f-87aa-7acef5ea56c7"  # Replace with your job ID
curl http://localhost:8080/api/v1/jobs/$JOB_ID | jq .
```

### 7. Transcode Video
```bash
VIDEO_ID="92503788-0298-4e79-9212-898ee9e66cc6"  # Replace with your video ID
curl -X POST http://localhost:8080/api/v1/videos/$VIDEO_ID/transcode \
  -H "Content-Type: application/json" \
  -d '{
    "format": "mp4",
    "codec": "h264",
    "resolution": "1280x720",
    "bitrate": "2000k"
  }'
```

### 8. Compress Video
```bash
VIDEO_ID="92503788-0298-4e79-9212-898ee9e66cc6"  # Replace with your video ID
curl -X POST http://localhost:8080/api/v1/videos/$VIDEO_ID/compress \
  -H "Content-Type: application/json" \
  -d '{
    "quality": "medium",
    "target_size": 5242880
  }'
```

**Quality Options:** `low`, `medium`, `high`

### 9. Stream Video
```bash
VIDEO_ID="92503788-0298-4e79-9212-898ee9e66cc6"  # Replace with your video ID
curl http://localhost:8080/api/v1/videos/$VIDEO_ID/stream -o output.mp4
```

### 10. Delete Video
```bash
VIDEO_ID="92503788-0298-4e79-9212-898ee9e66cc6"  # Replace with your video ID
curl -X DELETE http://localhost:8080/api/v1/videos/$VIDEO_ID | jq .
```

## Verify Generated Files

### Check Uploaded Videos
```bash
ls -lh storage/uploads/
```

### Check Generated Thumbnails
```bash
VIDEO_ID="92503788-0298-4e79-9212-898ee9e66cc6"  # Replace with your video ID
ls -lh storage/thumbnails/$VIDEO_ID/
```

### Check Processed Videos
```bash
ls -lh storage/processed/
```

## Common Issues & Solutions

### Issue 1: "No video file provided"
**Solution:** Make sure you're NOT including `-H "Content-Type: multipart/form-data"` in your upload request. Use `-F` flag only:
```bash
# CORRECT
curl -X POST http://localhost:8080/api/v1/videos/upload -F "video=@/path/to/video.mp4"

# WRONG
curl -X POST http://localhost:8080/api/v1/videos/upload \
  -H "Content-Type: multipart/form-data" \
  -F "video=@/path/to/video.mp4"
```

### Issue 2: "invalid MIME type for video"
**Solution:** This has been fixed. The API now accepts `application/octet-stream` MIME type for video files with valid video extensions (.mp4, .avi, .mov, etc.).

### Issue 3: "Invalid parameters: EOF"
**Solution:** Make sure to include the required JSON parameters with proper Content-Type header:
```bash
curl -X POST http://localhost:8080/api/v1/videos/$VIDEO_ID/thumbnail \
  -H "Content-Type: application/json" \
  -d '{"count": 3}'
```

### Issue 4: Server Logs
View server logs:
```bash
# If running in background
tail -f /tmp/server.log

# Or check recent logs
tail -30 /tmp/server.log
```

## Complete Test Workflow

```bash
#!/bin/bash

# 1. Upload a video
echo "Uploading video..."
VIDEO_ID=$(curl -s -X POST http://localhost:8080/api/v1/videos/upload \
  -F "video=@sample.mp4" | jq -r '.data.video_id')
echo "Video ID: $VIDEO_ID"

# 2. Verify upload
echo "\nVideo details:"
curl -s http://localhost:8080/api/v1/videos/$VIDEO_ID | jq .

# 3. Generate thumbnails
echo "\nGenerating thumbnails..."
THUMB_JOB=$(curl -s -X POST http://localhost:8080/api/v1/videos/$VIDEO_ID/thumbnail \
  -H "Content-Type: application/json" \
  -d '{"count": 3}' | jq -r '.data.id')
echo "Thumbnail Job: $THUMB_JOB"

# 4. Wait and check thumbnails
sleep 3
echo "\nThumbnail files:"
ls -lh storage/thumbnails/$VIDEO_ID/

# 5. Transcode video
echo "\nTranscoding video..."
TRANS_JOB=$(curl -s -X POST http://localhost:8080/api/v1/videos/$VIDEO_ID/transcode \
  -H "Content-Type: application/json" \
  -d '{"format":"mp4","codec":"h264","resolution":"1280x720","bitrate":"2000k"}' | jq -r '.data.id')
echo "Transcode Job: $TRANS_JOB"

# 6. Check metrics
echo "\nSystem metrics:"
curl -s http://localhost:8080/api/v1/metrics | jq .

echo "\nTest complete!"
```

## Monitoring

### Check System Metrics
```bash
curl http://localhost:8080/api/v1/metrics | jq .
```

### Worker Pool Status
Check server logs for worker activity:
```bash
grep "Worker" /tmp/server.log | tail -20
```

### Job Processing
Check job-related logs:
```bash
grep "job" /tmp/server.log | tail -20
```

## Performance Testing

### Concurrent Uploads
```bash
for i in {1..5}; do
  curl -X POST http://localhost:8080/api/v1/videos/upload \
    -F "video=@sample.mp4" &
done
wait
```

### Load Testing
```bash
# Install Apache Bench if needed
# apt-get install apache2-utils

ab -n 100 -c 10 http://localhost:8080/api/v1/health
```

## Troubleshooting

### Check if server is running
```bash
ps aux | grep video-processor
netstat -tuln | grep 8080
```

### Restart server
```bash
pkill -f video-processor
./bin/video-processor > /tmp/server.log 2>&1 &
```

### Clean storage
```bash
# Be careful - this deletes all uploaded/processed files
rm -rf storage/uploads/*
rm -rf storage/processed/*
rm -rf storage/thumbnails/*
```

### Rebuild project
```bash
make build
```

## API Error Codes

| Status Code | Meaning |
|-------------|---------|
| 200 | Success |
| 201 | Resource created |
| 400 | Bad request (invalid parameters) |
| 404 | Resource not found |
| 500 | Internal server error |

## Notes

1. **File Uploads:** Always use absolute paths or paths relative to current directory
2. **Job Processing:** Jobs are processed asynchronously by a worker pool (4 workers by default)
3. **MIME Types:** The API accepts both proper video MIME types and `application/octet-stream`
4. **File Extensions:** Supported formats: mp4, avi, mov, mkv, webm, flv, wmv, m4v, 3gp
5. **Max File Size:** Default is 500MB (configurable in .env)
