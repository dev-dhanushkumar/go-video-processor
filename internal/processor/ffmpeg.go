package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type FFmpegProcessor struct {
	ffmpegPath  string
	ffprobePath string
}

func NewFFmpegProcessor(ffmpegPath, ffprobePath string) *FFmpegProcessor {
	return &FFmpegProcessor{
		ffmpegPath:  ffmpegPath,
		ffprobePath: ffprobePath,
	}
}

// VideoMetadata contains metadata extracted from video file
type VideoMetadata struct {
	Duration   float64
	Resolution string
	Codec      string
	Width      int
	Height     int
	Bitrate    int64
	FrameRate  float64
}

// GetMetadata extracts metadata from video file using ffprobe
func (p *FFmpegProcessor) GetMetadata(ctx context.Context, inputPath string) (*VideoMetadata, error) {
	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		inputPath,
	}

	cmd := exec.CommandContext(ctx, p.ffprobePath, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get video metadata: %w", err)
	}

	var result struct {
		Streams []struct {
			CodecName  string `json:"codec_name"`
			CodecType  string `json:"codec_type"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			RFrameRate string `json:"r_frame_rate"`
		} `json:"streams"`
		Format struct {
			Duration string `json:"duration"`
			BitRate  string `json:"bit_rate"`
		} `json:"format"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	metadata := &VideoMetadata{}

	// Extract video stream info
	for _, stream := range result.Streams {
		if stream.CodecType == "video" {
			metadata.Codec = stream.CodecName
			metadata.Width = stream.Width
			metadata.Height = stream.Height
			metadata.Resolution = fmt.Sprintf("%dx%d", stream.Width, stream.Height)

			// Parse frame rate
			if stream.RFrameRate != "" {
				parts := strings.Split(stream.RFrameRate, "/")
				if len(parts) == 2 {
					num, _ := strconv.ParseFloat(parts[0], 64)
					den, _ := strconv.ParseFloat(parts[1], 64)
					if den != 0 {
						metadata.FrameRate = num / den
					}
				}
			}
			break
		}
	}

	// Parse duration
	if result.Format.Duration != "" {
		duration, err := strconv.ParseFloat(result.Format.Duration, 64)
		if err == nil {
			metadata.Duration = duration
		}
	}

	// Parse bitrate
	if result.Format.BitRate != "" {
		bitrate, err := strconv.ParseInt(result.Format.BitRate, 10, 64)
		if err == nil {
			metadata.Bitrate = bitrate
		}
	}

	return metadata, nil
}

// ValidateFFmpeg checks if FFmpeg and FFprobe are available
func (p *FFmpegProcessor) ValidateFFmpeg() error {
	// Check ffmpeg
	if err := exec.Command(p.ffmpegPath, "-version").Run(); err != nil {
		return fmt.Errorf("ffmpeg not found at %s: %w", p.ffmpegPath, err)
	}

	// Check ffprobe
	if err := exec.Command(p.ffprobePath, "-version").Run(); err != nil {
		return fmt.Errorf("ffprobe not found at %s: %w", p.ffprobePath, err)
	}

	return nil
}
