# Videos Downloader

A command-line tool written in Go for downloading videos and clips from various online platforms using [yt-dlp](https://github.com/yt-dlp/yt-dlp) and [ffmpeg](https://github.com/FFmpeg/FFmpeg).

## Features
- Download full videos or clips with specified timestamps

- Process downloads concurrently for faster performance

- Supports downloading from [1000+ sites](https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md)



## Clip Processing Modes
### When downloading clips, there are two modes available:
- **Normal Mode** (Default):
  - Re-encodes the video using GPU or CPU encoding.
  - More precise clip cutting.
  - Slower processing due to re-encoding time.
  - The tool automatically:
    - Detects available GPU.
    - Uses appropriate encoder:
      - NVIDIA: h264_nvenc
      - AMD: h264_amf
      - Intel: h264_qsv
    - Tests GPU encoder compatibility.
    - Falls back to CPU encoder (libx264) if no GPU detected or the GPU encoder is not working


- **Fast Mode** (Using `-fast` flag):
  - Copies the video stream directly without re-encoding.
  - Faster processing.
  - May not cut clips as precisely as normal mode:
    - Clips can start a few seconds before the specified start time.
    - Clips can have frozen frames at the beginning.


## Installation

### Option 1 - Download Latest Release (Windows only)
1. Download the [latest release](https://drive.google.com/file/d/1hPOFOAebPspRruxmSvxdyORNsBoGDuJn/view?usp=drive_link).
   - The release contains the following files:
     - `downloader.exe`: Main executable
     - `ffmpeg.exe`: FFmpeg binary for clip processing
     - `yt-dlp.exe`: yt-dlp binary for downloading videos
     - `urls.txt`: File to add your video URLs for download
2. Extract the ZIP file.
3. Add video URLs to `urls.txt`.
4. Run `downloader.exe` from the command line or double-click it.

### Option 2 - Build from Source

#### Prerequisites
- [Go](https://go.dev/doc/install) 1.22 or later
- [Git](https://git-scm.com/downloads)
- `ffmpeg` and `yt-dlp` (see below for installation per OS)

### 1. Clone Repository
```bash
git clone https://github.com/mohamedhesham1619/videos-downloader
cd videos-downloader 
```

### 2. Install Dependencies
#### Windows
1. Download [yt-dlp.exe](https://github.com/yt-dlp/yt-dlp/releases) and place it in the `release` folder.

2. Download [ffmpeg](https://github.com/BtbN/FFmpeg-Builds/releases) 
    - Download `ffmpeg-master-latest-win64-gpl.zip` 
    - Extract it and copy `ffmpeg.exe` from the `bin` folder to the `release` folder.

#### Linux

``` bash
# Ubuntu/Debian
sudo apt update && sudo apt install -y ffmpeg yt-dlp

# Arch Linux
sudo pacman -S ffmpeg yt-dlp

# Fedora
sudo dnf install ffmpeg yt-dlp
```

#### macOS
```bash
# Using Homebrew
brew install ffmpeg yt-dlp
```

### 3. Build the Project
```bash
# Windows
go build -o release/downloader.exe cmd/main.go

# Linux/macOS
go build -o release/downloader cmd/main.go
```
## Usage

### Setting Up URLs

Create a file named `urls.txt` in the same directory as the executable (e.g., in the release folder if running from there).

Add video URLs you want to download, one per line. You can also specify clips by adding timestamps.

**Supported formats:**

- Full video:  
  ```
  https://youtube.com/watch?v=example
  ```
- Clip from a video (format: `<url> <start_time>-<end_time>`, timestamps in `HH:MM:SS-HH:MM:SS`):  
  ```
  https://youtube.com/watch?v=example 00:01:30-00:02:45
  ```


### Example `urls.txt` File
```plaintext
https://youtube.com/watch?v=example
https://facebook.com/video/example 00:05:00-00:10:00
https://twitter.com/video/example 00:00:30-00:01:00
https://tiktok.com/@user/video/example 
```

### Running the Program
```powershell
# Basic usage (downloads to ./downloads)
./downloader.exe

# Custom download path
./downloader.exe -path "D:\Videos"

# Fast mode (no re-encoding)
./downloader.exe -fast
```

## Demo

[![Video Demo](https://img.youtube.com/vi/lSFwxTx_bD4/maxresdefault.jpg)](https://youtu.be/lSFwxTx_bD4 "Videos Downloader Demo")
