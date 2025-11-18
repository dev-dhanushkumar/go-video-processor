package processor

import (
	"context"
	"fmt"
	"os/exec"
)

type CompressOptions struct {
	OutputPath string
	Quality    int    // CRF value (18-28 recommended, lower is better)
	Preset     string // ultrafast, superfast, veryfast, faster, fast, medium, slow, slower, veryslow
	MaxBitrate string // e.g., "1M", "2M"
}

// Compress reduces video file size while maintaining quality
func (p *FFmpegProcessor) Compress(ctx context.Context, inputPath string, opts CompressOptions) error {
	// Default values
	if opts.Quality == 0 {
		opts.Quality = 23 // Good balance between quality and size
	}
	if opts.Preset == "" {
		opts.Preset = "medium"
	}

	args := []string{
		"-i", inputPath,
		"-c:v", "libx264",
		"-crf", fmt.Sprintf("%d", opts.Quality),
		"-preset", opts.Preset,
		"-c:a", "aac",
		"-b:a", "128k", // Audio bitrate
	}

	// Add max bitrate if specified
	if opts.MaxBitrate != "" {
		args = append(args, "-maxrate", opts.MaxBitrate, "-bufsize", opts.MaxBitrate)
	}

	args = append(args, "-y", opts.OutputPath)

	cmd := exec.CommandContext(ctx, p.ffmpegPath, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compression failed: %w, output: %s", err, string(output))
	}

	return nil
}

// CompressWithTargetSize compresses video to approximate target file size
func (p *FFmpegProcessor) CompressWithTargetSize(ctx context.Context, inputPath string, outputPath string, targetSizeMB float64) error {
	// Get video duration first
	metadata, err := p.GetMetadata(ctx, inputPath)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	// Calculate target bitrate
	// Formula: bitrate (kbps) = (target size in MB * 8192) / duration in seconds
	// Subtract audio bitrate (128 kbps)
	audioBitrate := 128.0 // kbps
	targetBitrate := (targetSizeMB * 8192.0 / metadata.Duration) - audioBitrate

	if targetBitrate <= 0 {
		return fmt.Errorf("target size too small for video duration")
	}

	args := []string{
		"-i", inputPath,
		"-c:v", "libx264",
		"-b:v", fmt.Sprintf("%.0fk", targetBitrate),
		"-c:a", "aac",
		"-b:a", "128k",
		"-preset", "medium",
		"-y", outputPath,
	}

	cmd := exec.CommandContext(ctx, p.ffmpegPath, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compression failed: %w, output: %s", err, string(output))
	}

	return nil
}
