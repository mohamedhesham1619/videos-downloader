package downloader

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"videos-downloader/internal/config"
	"videos-downloader/internal/models"
	"videos-downloader/internal/utils"
)

type Downloader struct {
	Config *config.Config
}

func New(cfg *config.Config) *Downloader {
	return &Downloader{
		Config: cfg,
	}
}

func (d *Downloader) Download(video models.VideoRequest) error {
	var command *exec.Cmd
	var err error

	if video.IsClip {
		command, err = d.buildClipDownloadCommand(video)

		if err != nil {
			return err
		}
	} else {
		command = d.buildFullDownloadCommand(video)
	}

	// Run the command
	fmt.Println("Downloading video:", video.Url)

	output, err := command.CombinedOutput()

	if err != nil {
		return fmt.Errorf("error downloading video: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("Download complete.")
	return nil
}

// prepare the command to download the whole video
func (d *Downloader) buildFullDownloadCommand(req models.VideoRequest) *exec.Cmd {

	// yt-dlp output template: "%(title).244s.%(ext)s"
	// - %(title)s: video title from metadata
	// - .244s: limits title to 244 characters to avoid path length issues on Windows
	// - %(ext)s: file extension based on selected format
	downloadPath := filepath.Join(d.Config.DownloadPath, "%(title).244s.%(ext)s")

	cmd := exec.Command("./yt-dlp",
		"-f", "b",
		"-o", downloadPath,
		req.Url)

	return cmd
}

// prepare the command to download the clip of the video
func (d *Downloader) buildClipDownloadCommand(req models.VideoRequest) (*exec.Cmd, error) {

	// Get both the URL and the title with the extension
	cmd := exec.Command("./yt-dlp",
		"-f", "b[ext=mp4]/bv*[ext=mp4]+ba[ext=m4a]/bv*+ba/b",
		"--print", "%(title).244s.%(ext)s\n%(url)s",
		"--encoding", "utf-8",
		"--no-download",
		"--windows-filenames",
		"--output-na-placeholder", "_",
		"--no-warnings", // Reduce noise in output
		req.Url,
	)

	// Run the command and get the output
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("error getting video info: %v\noutput:%v", err, string(output))
	}

	// Split output into lines
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	if len(lines) < 2 {
		return nil, fmt.Errorf("expected both URL and title but got:%v", lines)
	}

	videoTitle := utils.SanitizeFilename(lines[0])
	videoURL := lines[1]

	downloadPath := filepath.Join(d.Config.DownloadPath, videoTitle)

	clipStart, clipDuration, err := utils.ParseClipDuration(req.ClipTimeRange)

	if err != nil {
		return nil, fmt.Errorf("error parsing clip duration: %v", err)
	}

	ffmpegCmd := exec.Command(
		"./ffmpeg",
		"-ss", clipStart,
		"-i", videoURL,
		"-t", clipDuration,

		// Copy without re-encoding (the clip may start few seconds earlier than the specified time or the first few seconds in the video can be frozed. Remove this flag to fix these issues but it will increase the cpu usage and slow down the download process)
		"-c", "copy",

		downloadPath,
	)

	return ffmpegCmd, nil
}
