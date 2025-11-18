package processor

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

type ThumbnailOptions struct {
	OutputDir  string
	Timestamps []float64 // Timestamps in seconds
	Width      int
	Height     int
	Format     string // jpg, png
}

// GenerateThumbnail creates a single thumbnail at specified timestamp
func (p *FFmpegProcessor) GenerateThumbnail(ctx context.Context, inputPath string, outputPath string, timestamp float64) error {
	args := []string{
		"-i", inputPath,
		"-ss", fmt.Sprintf("%.2f", timestamp),
		"-vframes", "1",
		"-q:v", "2", // Quality (1-31, lower is better)
		"-y",
		outputPath,
	}

	cmd := exec.CommandContext(ctx, p.ffmpegPath, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("thumbnail generation failed: %w, output: %s", err, string(output))
	}

	return nil
}

// GenerateThumbnails creates multiple thumbnails at different timestamps
func (p *FFmpegProcessor) GenerateThumbnails(ctx context.Context, inputPath string, opts ThumbnailOptions) ([]string, error) {
	if len(opts.Timestamps) == 0 {
		return nil, fmt.Errorf("no timestamps provided")
	}

	var thumbnailPaths []string
	format := opts.Format
	if format == "" {
		format = "jpg"
	}

	for i, timestamp := range opts.Timestamps {
		filename := fmt.Sprintf("thumb_%d.%s", i+1, format)
		outputPath := filepath.Join(opts.OutputDir, filename)

		args := []string{
			"-i", inputPath,
			"-ss", fmt.Sprintf("%.2f", timestamp),
			"-vframes", "1",
		}

		// Add size if specified
		if opts.Width > 0 && opts.Height > 0 {
			args = append(args, "-s", fmt.Sprintf("%dx%d", opts.Width, opts.Height))
		}

		args = append(args, "-q:v", "2", "-y", outputPath)

		cmd := exec.CommandContext(ctx, p.ffmpegPath, args...)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return thumbnailPaths, fmt.Errorf("thumbnail %d generation failed: %w, output: %s", i+1, err, string(output))
		}

		thumbnailPaths = append(thumbnailPaths, outputPath)
	}

	return thumbnailPaths, nil
}

// GenerateAutoThumbnails creates thumbnails at evenly distributed intervals
func (p *FFmpegProcessor) GenerateAutoThumbnails(ctx context.Context, inputPath string, outputDir string, count int) ([]string, error) {
	// Get video duration
	metadata, err := p.GetMetadata(ctx, inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	if metadata.Duration <= 0 {
		return nil, fmt.Errorf("invalid video duration")
	}

	// Calculate timestamps
	interval := metadata.Duration / float64(count+1)
	timestamps := make([]float64, count)
	for i := 0; i < count; i++ {
		timestamps[i] = interval * float64(i+1)
	}

	return p.GenerateThumbnails(ctx, inputPath, ThumbnailOptions{
		OutputDir:  outputDir,
		Timestamps: timestamps,
		Format:     "jpg",
	})
}
