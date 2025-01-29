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

// prepare the download path
// if the user provides a path flag, the downloaded videos will be saved in that directory. Otherwise, they will be saved in the current directory.
func buildDownloadPath(basePath string) string {
    downloadPath := "%(title)s.%(ext)s"
    if basePath != "" {
        downloadPath = strings.ReplaceAll(basePath, `\`, "/") + "/" + downloadPath
    }
    return downloadPath
}

// prepare the command to download the video
func buildCommand(req videoRequest) *exec.Cmd {
	downloadPath := buildDownloadPath(*pathFlag)
	cmd := exec.Command("./yt-dlp", "-o", downloadPath, req.url)

	// if the user wants to download a clip of the video, add the clip duration to the command
	if(req.isClip) {
		cmd.Args = append(cmd.Args, "--download-sections", fmt.Sprintf("*%v", req.clipDuration))
	}

	cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

	return cmd
}

func downloadVideo(req videoRequest) {
	command := buildCommand(req)

	// Run the command
	fmt.Println("Downloading video:", req.url)
	err := command.Run()
	if err != nil {
		fmt.Println("Error downloading video:", err)
		return
	}

	fmt.Println("Download complete.")
}
