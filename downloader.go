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
var pathFlage = flag.String("path", "", "path to the download directory")

func main() {
	flag.Parse()

	urls, err := extractUrls("urls.txt")

	if err != nil {
		fmt.Println("couldn't extract urls from the file")
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(urls))

	// start downloading videos concurrently
	for i := 0; i < len(urls); i++ {
		go func() {
			downloadVideo(urls[i])
			wg.Done()
		}()
	}

	// ensure all goroutines complete
	wg.Wait()
}

// extract urls from a text file into a slice of strings
func extractUrls(fileName string) ([]string, error) {
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

func downloadVideo(videoURL string) {

	downloadPath := "%(title)s.%(ext)s"
	
	if *pathFlage != "" {
		downloadPath = strings.ReplaceAll(*pathFlage, `\`, "/") + "/" + downloadPath
	}
	
	cmd := exec.Command("./yt-dlp", "-o", downloadPath, videoURL)

	// Set the output to the terminal for progress and errors
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	fmt.Println("Downloading video:", videoURL)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error downloading video:", err)
		return
	}

	fmt.Println("Download complete.")
}
