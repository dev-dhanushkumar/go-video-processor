#!/bin/bash

# FFmpeg installation script for different platforms

echo "Installing FFmpeg..."

# Detect OS
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    if command -v apt-get &> /dev/null; then
        # Debian/Ubuntu
        sudo apt-get update
        sudo apt-get install -y ffmpeg
    elif command -v yum &> /dev/null; then
        # CentOS/RHEL
        sudo yum install -y epel-release
        sudo yum install -y ffmpeg
    else
        echo "Unsupported Linux distribution"
        exit 1
    fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    if command -v brew &> /dev/null; then
        brew install ffmpeg
    else
        echo "Homebrew not found. Please install Homebrew first."
        exit 1
    fi
else
    echo "Unsupported operating system"
    exit 1
fi

# Verify installation
if command -v ffmpeg &> /dev/null && command -v ffprobe &> /dev/null; then
    echo "FFmpeg installed successfully!"
    ffmpeg -version
else
    echo "FFmpeg installation failed"
    exit 1
fi
