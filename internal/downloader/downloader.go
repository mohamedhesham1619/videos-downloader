package downloader

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"videos-downloader/internal/config"
	"videos-downloader/internal/models"
	"videos-downloader/internal/utils"

	"github.com/fatih/color"
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
	fmt.Printf("Downloading video: %s\n\n", video.Url)
	output, err := command.CombinedOutput()

	if err != nil {
		return fmt.Errorf("%s%v%s", color.RedString("error downloading ("), video.Url, color.RedString("): "+err.Error()+"\nOutput: "+string(output)))
	}

	fmt.Printf("%s %s\n\n", color.GreenString("Download completed:"), video.Url)
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

// prepare the command to download a clip of the video
func (d *Downloader) buildClipDownloadCommand(req models.VideoRequest) (*exec.Cmd, error) {

	// Get both the direct URL and the title with the extension
	cmd := exec.Command("./yt-dlp",
		"-f", "bv*+ba/b/best",
		"--print", "%(title).244s\n%(urls)s",
		"--encoding", "utf-8",
		"--no-playlist",
		"--no-download",
		"--no-warnings",
		req.Url,
	)

	// Run the command and get the output
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("error getting (%v) info: %v\noutput:%v", req.Url, err, string(output))
	}

	// Split output into lines
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	if len(lines) < 2 {
		return nil, fmt.Errorf("expected both URL and title but got:%v", lines)
	}

	videoTitle := utils.SanitizeFilename(lines[0]) + ".mp4"
	downloadPath := filepath.Join(d.Config.DownloadPath, videoTitle)

	clipStart, clipDuration, err := utils.ParseClipDuration(req.ClipTimeRange)

	if err != nil {
		return nil, fmt.Errorf("error parsing clip duration: %v", err)
	}

	var ffmpegCmd *exec.Cmd

	// Multiple URLs (separate video and audio)
	if len(lines) > 2 {
		videoUrl := lines[1]
		audioUrl := lines[2]

		ffmpegCmd = exec.Command(
			"./ffmpeg",
			"-ss", clipStart, // Seek position for video
			"-i", videoUrl,
			"-ss", clipStart, // Seek position for audio
			"-i", audioUrl,
			"-t", clipDuration,
		)

		// If fast mode is enabled, copy the streams without re-encoding
		// Otherwise, use the encoder
		if d.Config.IsFastMode {
			ffmpegCmd.Args = append(ffmpegCmd.Args,
				"-c", "copy",
				downloadPath)
		} else {
			ffmpegCmd.Args = append(ffmpegCmd.Args,
				"-c:v", d.Config.Encoder,
				"-q:a", "0", // Highest audio quality
				downloadPath)
		}
	} else { // Single URL (combined format)

		url := lines[1]

		ffmpegCmd = exec.Command(
			"./ffmpeg",
			"-ss", clipStart, // Seek position
			"-i", url,
			"-t", clipDuration,
		)
		if d.Config.IsFastMode {

			ffmpegCmd.Args = append(ffmpegCmd.Args,
				"-c", "copy", // Copy the stream without re-encoding
				downloadPath,
			)
		} else {
			ffmpegCmd.Args = append(ffmpegCmd.Args,
				"-c:v", d.Config.Encoder,
				"-q:a", "0", // Highest audio quality
				downloadPath,
			)
		}
	}

	return ffmpegCmd, nil
}
