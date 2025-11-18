# API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Currently, no authentication is required (will be added in future versions).

---

## Endpoints

### Health & Monitoring

#### 1. Health Check
Check if the service is running.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2024-11-18T10:30:00Z",
    "version": "1.0.0"
  }
}
```

#### 2. System Metrics
Get system metrics and statistics.

**Endpoint:** `GET /metrics`

**Response:**
```json
{
  "success": true,
  "data": {
    "total_videos": 150,
    "total_jobs": 450,
    "active_jobs": 3,
    "completed_jobs": 420,
    "failed_jobs": 27,
    "storage_used_bytes": 5368709120,
    "active_workers": 4,
    "queue_depth": 2
  }
}
```

---

### Video Management

#### 3. Upload Video
Upload a new video file.

**Endpoint:** `POST /videos/upload`

**Content-Type:** `multipart/form-data`

**Parameters:**
- `video` (file, required): Video file to upload

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/videos/upload \
  -F "video=@sample.mp4"
```

**Response:**
```json
{
  "success": true,
  "message": "Video uploaded successfully",
  "data": {
    "video_id": "550e8400-e29b-41d4-a716-446655440000",
    "message": "Video uploaded and metadata extracted"
  }
}
```

#### 4. List Videos
Get a paginated list of all videos.

**Endpoint:** `GET /videos`

**Query Parameters:**
- `limit` (integer, optional): Number of videos per page (default: 10)
- `offset` (integer, optional): Offset for pagination (default: 0)

**Example:**
```bash
curl "http://localhost:8080/api/v1/videos?limit=10&offset=0"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "videos": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "original_filename": "sample.mp4",
        "file_path": "./storage/uploads/550e8400-e29b-41d4-a716-446655440000.mp4",
        "file_size": 10485760,
        "mime_type": "video/mp4",
        "duration": 120.5,
        "resolution": "1920x1080",
        "codec": "h264",
        "created_at": "2024-11-18T10:00:00Z",
        "updated_at": "2024-11-18T10:00:00Z"
      }
    ],
    "total": 1
  }
}
```

#### 5. Get Video Details
Retrieve details of a specific video.

**Endpoint:** `GET /videos/{id}`

**Path Parameters:**
- `id` (string, required): Video ID

**Example:**
```bash
curl "http://localhost:8080/api/v1/videos/550e8400-e29b-41d4-a716-446655440000"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "original_filename": "sample.mp4",
    "file_path": "./storage/uploads/550e8400-e29b-41d4-a716-446655440000.mp4",
    "file_size": 10485760,
    "mime_type": "video/mp4",
    "duration": 120.5,
    "resolution": "1920x1080",
    "codec": "h264",
    "created_at": "2024-11-18T10:00:00Z",
    "updated_at": "2024-11-18T10:00:00Z"
  }
}
```

#### 6. Stream Video
Stream video content with range support.

**Endpoint:** `GET /videos/{id}/stream`

**Path Parameters:**
- `id` (string, required): Video ID

**Example:**
```bash
curl "http://localhost:8080/api/v1/videos/550e8400-e29b-41d4-a716-446655440000/stream" \
  --output video.mp4
```

**Response:** Binary video stream

#### 7. Delete Video
Delete a video and its associated files.

**Endpoint:** `DELETE /videos/{id}`

**Path Parameters:**
- `id` (string, required): Video ID

**Example:**
```bash
curl -X DELETE "http://localhost:8080/api/v1/videos/550e8400-e29b-41d4-a716-446655440000"
```

**Response:**
```json
{
  "success": true,
  "message": "Video deleted successfully"
}
```

---

### Video Processing

#### 8. Transcode Video
Create a job to transcode video to different format/resolution.

**Endpoint:** `POST /videos/{id}/transcode`

**Path Parameters:**
- `id` (string, required): Video ID

**Request Body:**
```json
{
  "resolution": "720p",
  "format": "mp4",
  "quality": 23
}
```

**Parameters:**
- `resolution` (string, optional): Target resolution (1080p, 720p, 480p, 360p)
- `format` (string, optional): Output format (mp4, webm, avi)
- `quality` (integer, optional): CRF quality value (18-28, lower is better)

**Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/videos/550e8400-e29b-41d4-a716-446655440000/transcode" \
  -H "Content-Type: application/json" \
  -d '{"resolution":"720p","format":"mp4","quality":23}'
```

**Response:**
```json
{
  "success": true,
  "message": "Transcode job created",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "video_id": "550e8400-e29b-41d4-a716-446655440000",
    "operation": "transcode",
    "status": "PENDING",
    "progress": 0,
    "created_at": "2024-11-18T10:05:00Z"
  }
}
```

#### 9. Compress Video
Create a job to compress video file.

**Endpoint:** `POST /videos/{id}/compress`

**Path Parameters:**
- `id` (string, required): Video ID

**Request Body:**
```json
{
  "quality": 23
}
```

**Parameters:**
- `quality` (integer, optional): CRF quality value (18-28, default: 23)

**Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/videos/550e8400-e29b-41d4-a716-446655440000/compress" \
  -H "Content-Type: application/json" \
  -d '{"quality":23}'
```

**Response:**
```json
{
  "success": true,
  "message": "Compress job created",
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "video_id": "550e8400-e29b-41d4-a716-446655440000",
    "operation": "compress",
    "status": "PENDING",
    "progress": 0,
    "created_at": "2024-11-18T10:06:00Z"
  }
}
```

#### 10. Generate Thumbnails
Create a job to generate video thumbnails.

**Endpoint:** `POST /videos/{id}/thumbnail`

**Path Parameters:**
- `id` (string, required): Video ID

**Request Body:**
```json
{
  "count": 4
}
```

**Parameters:**
- `count` (integer, optional): Number of thumbnails to generate (default: 4)

**Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/videos/550e8400-e29b-41d4-a716-446655440000/thumbnail" \
  -H "Content-Type: application/json" \
  -d '{"count":4}'
```

**Response:**
```json
{
  "success": true,
  "message": "Thumbnail job created",
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "video_id": "550e8400-e29b-41d4-a716-446655440000",
    "operation": "thumbnail",
    "status": "PENDING",
    "progress": 0,
    "created_at": "2024-11-18T10:07:00Z"
  }
}
```

---

### Job Management

#### 11. Get Job Status
Check the status of a processing job.

**Endpoint:** `GET /jobs/{id}`

**Path Parameters:**
- `id` (string, required): Job ID

**Example:**
```bash
curl "http://localhost:8080/api/v1/jobs/660e8400-e29b-41d4-a716-446655440001"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "job_id": "660e8400-e29b-41d4-a716-446655440001",
    "status": "COMPLETED",
    "progress": 100,
    "output_path": "./storage/processed/video_transcoded_720p.mp4"
  }
}
```

**Job Statuses:**
- `PENDING`: Job is waiting in queue
- `PROCESSING`: Job is currently being processed
- `COMPLETED`: Job completed successfully
- `FAILED`: Job failed with error

---

## Error Responses

All error responses follow this format:

```json
{
  "success": false,
  "error": "Error message description"
}
```

**Common HTTP Status Codes:**
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request parameters
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

---

## Rate Limiting

Currently no rate limiting is implemented. This will be added in future versions.

---

## Examples

### Complete Workflow Example

```bash
# 1. Upload a video
VIDEO_RESPONSE=$(curl -X POST http://localhost:8080/api/v1/videos/upload \
  -F "video=@sample.mp4")
VIDEO_ID=$(echo $VIDEO_RESPONSE | jq -r '.data.video_id')

# 2. Start transcoding
JOB_RESPONSE=$(curl -X POST "http://localhost:8080/api/v1/videos/$VIDEO_ID/transcode" \
  -H "Content-Type: application/json" \
  -d '{"resolution":"720p","format":"mp4"}')
JOB_ID=$(echo $JOB_RESPONSE | jq -r '.data.id')

# 3. Check job status
curl "http://localhost:8080/api/v1/jobs/$JOB_ID"

# 4. Generate thumbnails
curl -X POST "http://localhost:8080/api/v1/videos/$VIDEO_ID/thumbnail" \
  -H "Content-Type: application/json" \
  -d '{"count":4}'

# 5. Stream the video
curl "http://localhost:8080/api/v1/videos/$VIDEO_ID/stream" --output downloaded.mp4
```

---

## WebSocket Support (Future)

WebSocket support for real-time job progress updates will be added in future versions.

---

## Versioning

API version is included in the URL path (`/api/v1`). Breaking changes will increment the version number.
