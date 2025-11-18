# Video Processing System

A scalable backend system built with Golang and FFmpeg for efficient video processing, storage, and retrieval.

## Features

- **Video Upload & Management**: Upload, retrieve, stream, and delete videos
- **Video Processing**:
  - Transcoding to multiple formats (MP4, WebM, AVI, etc.)
  - Resolution conversion (1080p, 720p, 480p, 360p)
  - Video compression with quality control
  - Automatic thumbnail generation
- **Asynchronous Processing**: Job queue system with worker pools
- **RESTful API**: Clean and intuitive API endpoints
- **Storage Management**: Organized file storage with cleanup
- **Monitoring**: Health checks and system metrics

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Video Processing**: FFmpeg
- **Job Queue**: Channel-based worker pools

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12+
- FFmpeg and FFprobe installed
- Git

### Installing FFmpeg

#### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install ffmpeg
```

#### macOS
```bash
brew install ffmpeg
```

#### Verify Installation
```bash
ffmpeg -version
ffprobe -version
```

## Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd go-video-processor
```

### 2. Setup Environment

```bash
# Copy environment template
cp .env.example .env

# Edit .env with your configuration
nano .env
```

### 3. Install Dependencies

```bash
make install-deps
```

### 4. Setup Database

```bash
# Create database
createdb video_processing

# Run migrations
make migrate-up
```

### 5. Run the Application

```bash
# Development mode
make run

# Or build and run
make build
./bin/video-processor
```

The server will start on `http://localhost:8080`

## API Documentation

### Health & Monitoring

#### Health Check
```bash
GET /api/v1/health
```

#### System Metrics
```bash
GET /api/v1/metrics
```

### Video Management

#### Upload Video
```bash
POST /api/v1/videos/upload
Content-Type: multipart/form-data

# cURL example
curl -X POST http://localhost:8080/api/v1/videos/upload \
  -F "video=@/path/to/video.mp4"
```

#### List Videos
```bash
GET /api/v1/videos?limit=10&offset=0
```

#### Get Video Details
```bash
GET /api/v1/videos/{video_id}
```

#### Stream Video
```bash
GET /api/v1/videos/{video_id}/stream
```

#### Delete Video
```bash
DELETE /api/v1/videos/{video_id}
```

### Video Processing

#### Transcode Video
```bash
POST /api/v1/videos/{video_id}/transcode
Content-Type: application/json

{
  "resolution": "720p",
  "format": "mp4",
  "quality": 23
}
```

#### Compress Video
```bash
POST /api/v1/videos/{video_id}/compress
Content-Type: application/json

{
  "quality": 23
}
```

#### Generate Thumbnails
```bash
POST /api/v1/videos/{video_id}/thumbnail
Content-Type: application/json

{
  "count": 4
}
```

### Job Management

#### Get Job Status
```bash
GET /api/v1/jobs/{job_id}
```

## Configuration

All configuration is done through environment variables. See `.env.example` for all available options.

Key configurations:
- `SERVER_PORT`: API server port (default: 8080)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`: Database connection
- `PROCESSING_WORKER_COUNT`: Number of concurrent workers (default: 4)
- `STORAGE_MAX_FILE_SIZE`: Maximum upload file size in bytes (default: 500MB)

## Project Structure

```
.
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/            # HTTP handlers and routes
│   ├── models/         # Data models
│   ├── services/       # Business logic
│   ├── repository/     # Database layer
│   ├── processor/      # FFmpeg integration
│   ├── queue/          # Job queue system
│   └── utils/          # Utilities
├── pkg/
│   └── config/         # Configuration
├── storage/            # File storage
├── migrations/         # Database migrations
├── docs/              # Documentation
└── scripts/           # Utility scripts
```

## Development

### Running Tests
```bash
make test
```

### Code Coverage
```bash
make test-coverage
```

### Linting
```bash
make lint
```

### Format Code
```bash
make format
```

## Docker

### Build Docker Image
```bash
make docker-build
```

### Run with Docker
```bash
make docker-run
```

### Docker Compose
```bash
# Start all services
make docker-compose-up

# Stop all services
make docker-compose-down
```

## Database Migrations

### Apply Migrations
```bash
make migrate-up
```

### Rollback Migrations
```bash
make migrate-down
```

## Performance Tuning

- Adjust `PROCESSING_WORKER_COUNT` based on CPU cores
- Set `PROCESSING_MAX_CONCURRENT_JOBS` to control queue size
- Configure `STORAGE_MAX_FILE_SIZE` based on your needs
- Use SSD storage for better I/O performance

## Monitoring

Access system metrics at:
```
GET /api/v1/metrics
```

Returns:
- Total videos
- Active/completed/failed jobs
- Storage usage
- Active workers
- Queue depth

## Troubleshooting

### FFmpeg not found
Ensure FFmpeg is installed and the path is correctly set in `.env`:
```
FFMPEG_PATH=/usr/bin/ffmpeg
FFPROBE_PATH=/usr/bin/ffprobe
```

### Database connection failed
Verify PostgreSQL is running and credentials are correct:
```bash
psql -h localhost -U postgres -d video_processing
```

### Out of memory during processing
Reduce `PROCESSING_WORKER_COUNT` or upgrade server resources.

## License

MIT License

## Author

Dhanush Kumar M  
Email: dev.dhanushkumar@gmail.com

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
