package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
	"videos-downloader/internal/models"

	"github.com/fatih/color"
)

// ReadUrlsFromFile reads URLs from a file and returns them as a slice of strings.
// It ignores empty lines and trims whitespace from each line.
func ReadUrlsFromFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return []string{}, fmt.Errorf("couldn't open the file: %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var urls []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line != "" {
			urls = append(urls, line)
		}
	}
	return urls, nil
}

// create a video request object from a line of text
// the line should be in one of these formats:
// - url only: "https://www.youtube.com/watch?v=video_id"	(downloads full video)
// - url and time range: "https://www.youtube.com/watch?v=video_id 00:00:00-00:02:00"	(downloads clip)
// for clip download, the time range must be in the format HH:MM:SS-HH:MM:SS
func ParseVideoRequest(line string) models.VideoRequest {

	// split the line by spaces
	parts := strings.Fields(line)

	req := models.VideoRequest{
		Url: parts[0],
	}

	// if the line contains a time range, add it to the request
	if len(parts) > 1 {
		req.IsClip = true
		req.ClipTimeRange = parts[1]
	}

	return req
}

// parse clip timing info
// for ffmpeg to accurately extract the needed clip, it needs the start time and clip duration in seconds
func ParseClipDuration(timeRange string) (start string, duration string, err error) {
	// Split the range into start and end times
	parts := strings.Split(timeRange, "-")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid time range format. Expected HH:MM:SS-HH:MM:SS")
	}

	startTime := parts[0]
	endTime := parts[1]

	// Parse times to calculate duration
	layout := "15:04:05"

	t1, err := time.Parse(layout, startTime)
	if err != nil {
		return "", "", fmt.Errorf("invalid start time: %v", err)
	}

	t2, err := time.Parse(layout, endTime)
	if err != nil {
		return "", "", fmt.Errorf("invalid end time: %v", err)
	}

	// Calculate duration in seconds
	durationSeconds := int(t2.Sub(t1).Seconds())

	// Convert duration to string
	duration = strconv.Itoa(durationSeconds)

	return startTime, duration, nil
}

// sanitize the filename to remove or replace characters that are problematic in filenames
func SanitizeFilename(filename string) string {

	replacements := map[rune]rune{
		'/':  '-',
		'\\': '-',
		':':  '-',
		'*':  '-',
		'?':  '-',
		'"':  '-',
		'<':  '-',
		'>':  '-',
		'|':  '-',
	}

	sanitized := []rune{}
	for _, r := range filename {
		if replaced, exists := replacements[r]; exists {
			sanitized = append(sanitized, replaced)
		} else if unicode.IsPrint(r) {
			sanitized = append(sanitized, r)
		}
	}

	return string(sanitized)
}

// Returns the right command to be executed based on the OS
// For Windows, the executable needs to be in the same directory as the Go binary and the command needs to be prefixed with "./"
//
// For other OSes, the command can be executed directly
func GetCommand(command string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("./%s", command)
	}
	return command
}

// Test if the GPU encoder is working
// If the command runs successfully and doesn't return an error, the encoder is working
func TestGpuEncoder(encoder string) bool {
	testCmd := exec.Command(
		GetCommand("ffmpeg"),
		"-hide_banner",
		"-loglevel", "error",
		"-f", "lavfi",
		"-i", "testsrc=duration=1",
		"-c:v", encoder,
		"-frames:v", "10",
		"-f", "null",
		"-",
	)
	return testCmd.Run() == nil
}

// FormatDuration formats a duration in seconds to a human-readable string (e.g., "2m 30s")
func FormatDuration(seconds int) string {
	m := seconds / 60
	s := seconds % 60
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// Formats a user-friendly message for clip downloads
func FormatClipDownloadMessage(timeRange string) string {
	_, durationStr, err := ParseClipDuration(timeRange)
	if err != nil {
		return ""
	}

	durationSecs, _ := strconv.Atoi(durationStr)
	startTime, endTime := strings.Split(timeRange, "-")[0], strings.Split(timeRange, "-")[1]

	return fmt.Sprintf("Downloading clip: %s duration (from %s to %s)",
		color.CyanString(FormatDuration(durationSecs)),
		color.YellowString(startTime),
		color.YellowString(endTime))
}
