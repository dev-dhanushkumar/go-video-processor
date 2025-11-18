package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Storage    StorageConfig
	Processing ProcessingConfig
	FFmpeg     FFmpegConfig
	Logging    LoggingConfig
}

type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host           string
	Port           int
	User           string
	Password       string
	Database       string
	MaxConnections int
}

type StorageConfig struct {
	UploadDir    string
	ProcessedDir string
	ThumbnailDir string
	MaxFileSize  int64
}

type ProcessingConfig struct {
	WorkerCount       int
	MaxConcurrentJobs int
	TempDir           string
	SupportedFormats  []string
}

type FFmpegConfig struct {
	Path              string
	FFprobePath       string
	DefaultCodec      string
	DefaultAudioCodec string
}

type LoggingConfig struct {
	Level  string
	Format string
	Output string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "localhost"),
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Host:           getEnv("DB_HOST", "localhost"),
			Port:           getEnvAsInt("DB_PORT", 5432),
			User:           getEnv("DB_USER", "postgres"),
			Password:       getEnv("DB_PASSWORD", "password"),
			Database:       getEnv("DB_NAME", "video_processing"),
			MaxConnections: getEnvAsInt("DB_MAX_CONNECTIONS", 25),
		},
		Storage: StorageConfig{
			UploadDir:    getEnv("STORAGE_UPLOAD_DIR", "./storage/uploads"),
			ProcessedDir: getEnv("STORAGE_PROCESSED_DIR", "./storage/processed"),
			ThumbnailDir: getEnv("STORAGE_THUMBNAIL_DIR", "./storage/thumbnails"),
			MaxFileSize:  getEnvAsInt64("STORAGE_MAX_FILE_SIZE", 524288000), // 500MB
		},
		Processing: ProcessingConfig{
			WorkerCount:       getEnvAsInt("PROCESSING_WORKER_COUNT", 4),
			MaxConcurrentJobs: getEnvAsInt("PROCESSING_MAX_CONCURRENT_JOBS", 10),
			TempDir:           getEnv("PROCESSING_TEMP_DIR", "/tmp/video-processing"),
			SupportedFormats:  getEnvAsSlice("PROCESSING_SUPPORTED_FORMATS", []string{"mp4", "avi", "mov", "mkv", "webm"}),
		},
		FFmpeg: FFmpegConfig{
			Path:              getEnv("FFMPEG_PATH", "/usr/bin/ffmpeg"),
			FFprobePath:       getEnv("FFPROBE_PATH", "/usr/bin/ffprobe"),
			DefaultCodec:      getEnv("FFMPEG_DEFAULT_CODEC", "libx264"),
			DefaultAudioCodec: getEnv("FFMPEG_DEFAULT_AUDIO_CODEC", "aac"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
			Output: getEnv("LOG_OUTPUT", "stdout"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}

func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Storage.MaxFileSize <= 0 {
		return fmt.Errorf("invalid max file size: %d", c.Storage.MaxFileSize)
	}

	if c.Processing.WorkerCount < 1 {
		return fmt.Errorf("worker count must be at least 1")
	}

	return nil
}
