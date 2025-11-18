package utils

import (
	"fmt"
	"mime/multipart"
	"strings"
)

var supportedVideoFormats = []string{
	"mp4", "avi", "mov", "mkv", "webm", "flv", "wmv", "m4v", "3gp",
}

var videoMimeTypes = map[string][]string{
	"mp4":  {"video/mp4"},
	"avi":  {"video/x-msvideo", "video/avi"},
	"mov":  {"video/quicktime"},
	"mkv":  {"video/x-matroska"},
	"webm": {"video/webm"},
	"flv":  {"video/x-flv"},
	"wmv":  {"video/x-ms-wmv"},
	"m4v":  {"video/x-m4v"},
	"3gp":  {"video/3gpp"},
}

// ValidateVideoFile validates if the uploaded file is a valid video
func ValidateVideoFile(file *multipart.FileHeader, maxSize int64) error {
	// Check file size
	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxSize)
	}

	// Check file extension
	ext := GetFileExtension(file.Filename)
	if !IsVideoFormat(ext) {
		return fmt.Errorf("unsupported video format: %s", ext)
	}

	// Check MIME type
	contentType := file.Header.Get("Content-Type")
	if !IsValidVideoMimeType(contentType, ext) {
		return fmt.Errorf("invalid MIME type for video: %s", contentType)
	}

	return nil
}

// IsVideoFormat checks if the format is a supported video format
func IsVideoFormat(format string) bool {
	format = strings.ToLower(format)
	for _, supported := range supportedVideoFormats {
		if format == supported {
			return true
		}
	}
	return false
}

// IsValidVideoMimeType checks if the MIME type matches the expected format
func IsValidVideoMimeType(mimeType, format string) bool {
	format = strings.ToLower(format)
	mimeType = strings.ToLower(mimeType)

	// If MIME type is generic video, allow it
	if strings.HasPrefix(mimeType, "video/") {
		return true
	}

	// Allow application/octet-stream for valid video file extensions
	// This handles cases where MIME type isn't properly detected
	if mimeType == "application/octet-stream" && IsVideoFormat(format) {
		return true
	}

	// Check specific MIME types
	if expectedTypes, ok := videoMimeTypes[format]; ok {
		for _, expectedType := range expectedTypes {
			if mimeType == expectedType {
				return true
			}
		}
	}

	return false
}

// GetSupportedFormats returns list of supported video formats
func GetSupportedFormats() []string {
	return supportedVideoFormats
}

// FormatFileSize formats file size in bytes to human readable format
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}
