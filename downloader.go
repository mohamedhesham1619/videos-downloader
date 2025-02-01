package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
)

// if the user provides a path flag, the downloaded videos will be saved in that directory. Otherwise, they will be saved in the current directory.
var pathFlag = flag.String("path", "", "path to the download directory")

// the file containing the urls of the videos to download
var urlsFlag = flag.String("urls", "urls.txt", "path to the file containing the urls of the videos to download")

func main() {
	flag.Parse()

	urls, err := readUrlsFromFile(*urlsFlag)

	if err != nil {
		fmt.Println("couldn't extract urls from the file")
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(urls))

	// start downloading videos concurrently
	for i := 0; i < len(urls); i++ {
		go func() {
			videoRequest := parseVideoRequest(urls[i])
			downloadVideo(videoRequest)
			wg.Done()
		}()
	}

	// ensure all goroutines complete
	wg.Wait()
}

// extract urls from a text file into a slice of strings
func readUrlsFromFile(fileName string) ([]string, error) {
	content, err := os.ReadFile(fileName)

	if err != nil {
		return []string{}, err
	}

	// split the content of the file into a slice using newline as separator
	urls := strings.Split(string(content), "\r\n")

	// remove the empty values from the slice
	urls = slices.DeleteFunc(urls, func(value string) bool { return value == "" })

	return urls, nil

}

type videoRequest struct {
	url          string
	isClip       bool   // default is false
	clipDuration string // default is ""
}

// parse the video request from the line in the file
// if the clip duration is provided, only the clip will be downloaded
// the clip duration should be in the format HH:MM:SS-HH:MM:SS
// example: https://www.youtube.com/watch?v=video_id 00:00:00-00:01:10
func parseVideoRequest(line string) videoRequest {

	// split the line by spaces
	parts := strings.Fields(line)

	req := videoRequest{
		url: parts[0],
	}

	if len(parts) > 1 {
		req.isClip = true
		req.clipDuration = parts[1]
	}

	return req
}

// prepare the command to download the whole video
func buildFullDownloadCommand(req videoRequest) *exec.Cmd {

	// this will download the video to the current directory with the title as the file name
	downloadPath := "%(title)s.%(ext)s"

	// if the user provides a path flag, the downloaded videos will be saved in that directory
	if *pathFlag != "" {

		// because the path flag is provided by the user, it may contain backslashes
		downloadPath = strings.ReplaceAll(*pathFlag, `\`, "/") + "/" + downloadPath

	}

	cmd := exec.Command("./yt-dlp", "-f", "b", req.url, "-o", downloadPath)

	return cmd
}

// prepare the command to download the clip of the video
func buildClipDownloadCommand(req videoRequest) *exec.Cmd {

	// Get both the URL and the title with the extension
	cmd := exec.Command("./yt-dlp",
		"-f", "b",
		"--print", "%(title)s.%(ext)s\n%(url)s",
		"--no-download",
		req.url,
	)

	output, err := cmd.Output()

	if err != nil {
		fmt.Println("Error getting video URL and title:", err)
	}

	// Split output into lines
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	if len(lines) < 2 {

		fmt.Println("expected both URL and title but got:", lines)
	}

	videoTitle := lines[0]
	videoURL := lines[1]

	downloadPath := *pathFlag + videoTitle

	clipDuration := strings.Split(req.clipDuration, "-")
	clipStart := clipDuration[0]
	clipEnd := clipDuration[1]

	ffmpegCmd := exec.Command(
		"./ffmpeg", "-i", videoURL,
		"-ss", clipStart, // Set the clip start and end time
		"-to", clipEnd,
		"-c", "copy", // Copy without re-encoding (fast)
		downloadPath,
	)
	
	return ffmpegCmd
}

// get the download command for the video request and run it
func downloadVideo(req videoRequest) {

	var command *exec.Cmd

	if req.isClip {
		command = buildClipDownloadCommand(req)
	} else {
		command = buildFullDownloadCommand(req)
	}

	// Run the command
	fmt.Println("Downloading video:", req.url)
	err := command.Run()

	if err != nil {
		fmt.Println("Error downloading video:", err)
		return
	}

	fmt.Println("Download complete.")
}
