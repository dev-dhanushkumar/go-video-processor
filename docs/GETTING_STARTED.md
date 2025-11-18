# Getting Started Guide

Welcome to the Video Processing System! This guide will help you get the application up and running.

## Prerequisites Checklist

Before you begin, ensure you have the following installed:

- [ ] Go 1.21 or higher
- [ ] PostgreSQL 12 or higher
- [ ] FFmpeg and FFprobe
- [ ] Git
- [ ] Make (optional, but recommended)

### Verify Prerequisites

```bash
# Check Go version
go version

# Check PostgreSQL
psql --version

# Check FFmpeg
ffmpeg -version
ffprobe -version

# Check Make
make --version
```

## Step-by-Step Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd go-video-processor
```

### 2. Install Go Dependencies

```bash
# Using Make
make install-deps

# Or manually
go mod download
go mod tidy
```

### 3. Setup Database

#### Create Database
```bash
# Connect to PostgreSQL
sudo -u postgres psql

# Create database and user
CREATE DATABASE video_processing;
CREATE USER video_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE video_processing TO video_user;
\q
```

#### Run Migrations
```bash
# Set environment variables
export DB_HOST=localhost
export DB_USER=video_user
export DB_PASSWORD=your_password
export DB_NAME=video_processing

# Run migrations
make migrate-up
```

Alternatively, run migrations manually:
```bash
psql -U video_user -d video_processing -f migrations/001_create_videos_table.up.sql
psql -U video_user -d video_processing -f migrations/002_create_jobs_table.up.sql
```

### 4. Configure Environment

```bash
# Copy environment template
cp .env.example .env

# Edit configuration
nano .env
```

Update the following critical settings:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=video_user
DB_PASSWORD=your_password
DB_NAME=video_processing
```

### 5. Create Storage Directories

The directories should already exist, but verify:
```bash
ls -la storage/
```

You should see:
- `uploads/`
- `processed/`
- `thumbnails/`

### 6. Build the Application

```bash
# Using Make
make build

# Or manually
go build -o bin/video-processor cmd/server/main.go
```

### 7. Run the Application

```bash
# Using Make
make run

# Or run the binary
./bin/video-processor
```

You should see output like:
```
{"timestamp":"2024-11-18T10:00:00Z","level":"INFO","message":"Starting Video Processing System..."}
{"timestamp":"2024-11-18T10:00:00Z","level":"INFO","message":"FFmpeg validated successfully"}
{"timestamp":"2024-11-18T10:00:00Z","level":"INFO","message":"Worker pool started with 4 workers"}
{"timestamp":"2024-11-18T10:00:00Z","level":"INFO","message":"Server listening on localhost:8080"}
```

## Verify Installation

### 1. Check Health Endpoint

```bash
curl http://localhost:8080/api/v1/health
```

Expected response:
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2024-11-18T10:00:00Z",
    "version": "1.0.0"
  }
}
```

### 2. Check Metrics

```bash
curl http://localhost:8080/api/v1/metrics
```

### 3. Test Video Upload

```bash
# Create a test video (or use your own)
curl -X POST http://localhost:8080/api/v1/videos/upload \
  -F "video=@/path/to/test-video.mp4"
```

## Common Issues and Solutions

### Issue: "Failed to connect to database"

**Solution:**
1. Verify PostgreSQL is running: `sudo systemctl status postgresql`
2. Check database credentials in `.env`
3. Ensure database exists: `psql -l | grep video_processing`

### Issue: "FFmpeg validation failed"

**Solution:**
1. Install FFmpeg: `sudo apt-get install ffmpeg` (Ubuntu/Debian)
2. Verify installation: `ffmpeg -version`
3. Update paths in `.env` if FFmpeg is in a custom location

### Issue: "Permission denied" for storage directories

**Solution:**
```bash
chmod -R 755 storage/
```

### Issue: Port 8080 already in use

**Solution:**
Change the port in `.env`:
```env
SERVER_PORT=8081
```

### Issue: "No such file or directory" for migrations

**Solution:**
Ensure you're in the project root directory when running `make migrate-up`

## Development Workflow

### Run in Development Mode

```bash
# Install air for live reload (optional)
go install github.com/air-verse/air@latest

# Run with auto-reload
air
```

### Run Tests

```bash
make test
```

### Check Code Coverage

```bash
make test-coverage
open coverage.html
```

### Format Code

```bash
make format
```

### Run Linter

```bash
# Install golangci-lint first
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
make lint
```

## Using Docker (Alternative Setup)

If you prefer using Docker:

### Build Docker Image

```bash
make docker-build
```

### Run with Docker Compose

```bash
# Start all services (app + PostgreSQL)
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

## Next Steps

1. **Read the API Documentation**: See `docs/API.md` for detailed API usage
2. **Understand Architecture**: Review `docs/ARCHITECTURE.md`
3. **Explore Examples**: Try the example workflows in the API documentation
4. **Configure Workers**: Adjust `PROCESSING_WORKER_COUNT` based on your CPU cores
5. **Setup Monitoring**: Consider adding Prometheus/Grafana for production

## Getting Help

- Check the README.md for general information
- Review API documentation in `docs/API.md`
- Look at example requests in this guide
- Check GitHub issues for common problems

## Production Deployment Checklist

Before deploying to production:

- [ ] Change database password
- [ ] Set up proper logging (file-based)
- [ ] Configure CORS properly
- [ ] Set up HTTPS/TLS
- [ ] Enable authentication
- [ ] Configure backup strategy
- [ ] Set up monitoring and alerts
- [ ] Use environment-specific configs
- [ ] Set appropriate worker count
- [ ] Configure storage limits
- [ ] Set up log rotation
- [ ] Enable rate limiting (future feature)

## Performance Tuning

### Optimize for Your Hardware

```env
# For 8-core CPU
PROCESSING_WORKER_COUNT=6

# For 4-core CPU
PROCESSING_WORKER_COUNT=3

# Adjust based on available memory
PROCESSING_MAX_CONCURRENT_JOBS=10
```

### Storage Optimization

- Use SSD for storage directories
- Set up separate partition for video files
- Configure regular cleanup jobs

### Database Optimization

- Enable connection pooling (already configured)
- Create indexes on frequently queried columns
- Set up read replicas for scaling

## Troubleshooting Logs

View application logs:
```bash
# If running with systemd
journalctl -u video-processor -f

# If running in terminal
# Logs will appear in stdout
```

## Support

For issues and questions:
- Email: dev.dhanushkumar@gmail.com
- GitHub Issues: [Create an issue]

---

**Congratulations!** Your video processing system is now ready to use! 🎉
