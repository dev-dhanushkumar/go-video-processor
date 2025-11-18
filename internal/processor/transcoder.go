package processor

import (
	"context"
	"fmt"
	"os/exec"
)

type TranscodeOptions struct {
	OutputPath string
	VideoCodec string
	AudioCodec string
	Resolution string // e.g., "1280x720"
	Bitrate    string // e.g., "2M"
	Preset     string // e.g., "medium", "fast", "slow"
	CRF        int    // Constant Rate Factor (0-51, lower is better quality)
	FrameRate  float64
}

// Transcode converts video to different format/codec
func (p *FFmpegProcessor) Transcode(ctx context.Context, inputPath string, opts TranscodeOptions) error {
	args := []string{
		"-i", inputPath,
		"-y", // Overwrite output file
	}

	// Video codec
	if opts.VideoCodec != "" {
		args = append(args, "-c:v", opts.VideoCodec)
	}

	// Audio codec
	if opts.AudioCodec != "" {
		args = append(args, "-c:a", opts.AudioCodec)
	}

	// Resolution
	if opts.Resolution != "" {
		args = append(args, "-s", opts.Resolution)
	}

	// Bitrate
	if opts.Bitrate != "" {
		args = append(args, "-b:v", opts.Bitrate)
	}

	// Preset
	if opts.Preset != "" {
		args = append(args, "-preset", opts.Preset)
	}

	// CRF
	if opts.CRF > 0 {
		args = append(args, "-crf", fmt.Sprintf("%d", opts.CRF))
	}

	// Frame rate
	if opts.FrameRate > 0 {
		args = append(args, "-r", fmt.Sprintf("%.2f", opts.FrameRate))
	}

	args = append(args, opts.OutputPath)

	cmd := exec.CommandContext(ctx, p.ffmpegPath, args...)

	// Get command output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("transcode failed: %w, output: %s", err, string(output))
	}

	return nil
}

// GetResolutionString converts resolution name to dimension string
func GetResolutionString(resolution string) string {
	resolutions := map[string]string{
		"1080p": "1920x1080",
		"720p":  "1280x720",
		"480p":  "854x480",
		"360p":  "640x360",
		"240p":  "426x240",
	}

	if res, ok := resolutions[resolution]; ok {
		return res
	}
	return resolution
}
