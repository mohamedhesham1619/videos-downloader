package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func main() {

	//videoURL := "https://www.facebook.com/reel/8176615212451714"

	//downloadVideo(videoURL)
	urls, err := extractUrls("urls.txt")

	if err != nil {
		fmt.Println("couldn't extract urls from the file")
		return
	}
	fmt.Println(urls)
}

// extract urls from a text file into a slice of strings
func extractUrls(fileName string) ([]string, error){
	content, err := os.ReadFile(fileName)

	if err != nil {
		return []string{}, err
	}

	// split the content of the file into a slice using newline as separator
	urls := strings.Split(string(content), "\r\n")
	
	// remove the empty values from the slice
	urls = slices.DeleteFunc(urls, func(value string) bool {return value == ""})

	return urls, nil
	
}

func downloadVideo(videoURL string) {
	cmd := exec.Command("./yt-dlp", "-o", "H:\\قرآن\\%(title)s.%(ext)s", videoURL)

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
