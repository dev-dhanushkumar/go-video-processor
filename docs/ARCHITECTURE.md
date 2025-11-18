# Architecture Documentation

## System Overview

The Video Processing System is designed as a scalable, microservices-ready backend application built with Go. It follows clean architecture principles with clear separation of concerns.

## Architecture Layers

### 1. API Layer (`internal/api`)
- **Handlers**: Process HTTP requests and responses
- **Middleware**: Cross-cutting concerns (logging, CORS, authentication)
- **Routes**: Define API endpoints and routing logic

### 2. Service Layer (`internal/services`)
- **VideoService**: Video upload, retrieval, and management
- **ProcessingService**: Video processing job creation and execution
- **StorageService**: Storage management and metrics

### 3. Repository Layer (`internal/repository`)
- **VideoRepository**: Video metadata persistence
- **JobRepository**: Processing job persistence
- Abstracts database operations with interfaces

### 4. Processing Layer (`internal/processor`)
- **FFmpegProcessor**: Core FFmpeg wrapper
- **Transcoder**: Video format/codec conversion
- **Compressor**: Video compression algorithms
- **Thumbnail**: Thumbnail generation

### 5. Queue Layer (`internal/queue`)
- **JobQueue**: Channel-based job queue
- **WorkerPool**: Concurrent job processing

### 6. Models Layer (`internal/models`)
- Data structures and DTOs
- Request/Response models

### 7. Utils Layer (`internal/utils`)
- File operations
- Video validation
- Logging utilities

## Data Flow

```
Client Request
    ↓
API Handler
    ↓
Service Layer
    ↓
Repository Layer ←→ Database
    ↓
Worker Pool
    ↓
FFmpeg Processor
    ↓
File System
```

## Asynchronous Processing

### Job Queue Architecture

```
Upload Video → Create Job → Add to Queue
                                ↓
                         Worker Pool (4 workers)
                                ↓
                         Process Video
                                ↓
                         Update Job Status
```

### Worker Pool Pattern

- Fixed number of workers (configurable)
- Channel-based job distribution
- Graceful shutdown support
- Progress tracking

## Database Schema

### Videos Table
- Stores video metadata
- Indexed by creation date
- Foreign key constraints to jobs

### Processing Jobs Table
- Tracks processing job state
- Stores job parameters as JSON
- Linked to video records

## Storage Organization

```
storage/
├── uploads/          # Original uploaded videos
├── processed/        # Transcoded/compressed videos
└── thumbnails/       # Generated thumbnails
    └── {video_id}/   # Thumbnails per video
```

## Scalability Considerations

### Current Architecture
- Single server deployment
- Channel-based job queue
- Local file storage

### Future Scalability Path

1. **Distributed Processing**
   - Redis-based job queue
   - Multiple worker nodes
   - Distributed lock management

2. **Cloud Storage**
   - S3/GCS integration
   - CDN for video delivery
   - Object storage for processed files

3. **Database Scaling**
   - Read replicas
   - Connection pooling
   - Query optimization

4. **Horizontal Scaling**
   - Load balancer
   - Stateless API servers
   - Shared job queue

## Security Architecture

### Current Implementation
- File validation (size, type, MIME)
- Filename sanitization
- Path traversal prevention

### Future Enhancements
- JWT authentication
- Rate limiting per user
- Video encryption at rest
- Signed URLs for streaming

## Monitoring & Observability

### Current Metrics
- Total videos
- Job statistics (active, completed, failed)
- Storage usage
- Worker pool status

### Future Monitoring
- Prometheus metrics
- Distributed tracing
- Error tracking (Sentry)
- Performance APM

## Error Handling

### Strategy
1. Input validation at API layer
2. Error propagation through layers
3. Structured error logging
4. Graceful degradation
5. Automatic job retries

### Recovery
- Failed job cleanup
- Orphaned file cleanup
- Database transaction rollback

## Performance Optimization

### Current Optimizations
- Concurrent video processing
- Streaming video responses
- Connection pooling
- Efficient file I/O

### Future Optimizations
- Video chunking for upload
- Progressive processing
- Caching layer (Redis)
- Background cleanup jobs

## Technology Decisions

### Why Go?
- Excellent concurrency primitives
- Fast compilation and execution
- Strong standard library
- Low memory footprint

### Why FFmpeg?
- Industry standard for video processing
- Supports all major formats
- Hardware acceleration support
- Battle-tested and reliable

### Why PostgreSQL?
- ACID compliance
- JSON support for job parameters
- Strong indexing capabilities
- Reliability and data integrity

### Why Gin Framework?
- Fast HTTP router
- Middleware support
- JSON validation
- Well-documented

## Deployment Architecture

### Development
```
Developer Machine
├── Go Application
├── PostgreSQL (Docker)
└── Local Storage
```

### Production
```
Load Balancer
    ↓
API Servers (N instances)
    ↓
PostgreSQL Cluster
    ↓
Shared Storage (NFS/S3)
```

## API Versioning Strategy

- URL-based versioning (`/api/v1`)
- Backward compatibility within major version
- Deprecation warnings for old endpoints
- Migration guides for breaking changes

## Testing Strategy

### Unit Tests
- Service layer logic
- Repository operations
- Utility functions

### Integration Tests
- API endpoint testing
- Database operations
- FFmpeg integration

### Load Tests
- Concurrent uploads
- Processing throughput
- Queue performance

## Configuration Management

### Environment Variables
- Server configuration
- Database credentials
- Storage paths
- Processing parameters

### Secrets Management
- Database passwords
- API keys (future)
- Encryption keys (future)

## Disaster Recovery

### Backup Strategy
- Database regular backups
- Video file backups
- Configuration backups

### Recovery Procedures
- Database restore
- File system recovery
- Job queue rebuild

## Future Architectural Improvements

1. **Microservices Split**
   - Upload service
   - Processing service
   - Streaming service

2. **Event-Driven Architecture**
   - Event bus (Kafka/RabbitMQ)
   - Event sourcing
   - CQRS pattern

3. **Service Mesh**
   - Istio/Linkerd
   - Service discovery
   - Circuit breakers

4. **Containerization**
   - Kubernetes deployment
   - Auto-scaling
   - Health checks

## API Gateway Pattern

Future implementation may include:
- API Gateway (Kong/Traefik)
- Rate limiting
- Authentication/Authorization
- Request aggregation

## Conclusion

The current architecture provides a solid foundation for a video processing system with clear paths for scaling and enhancement. The modular design allows for incremental improvements without major rewrites.
