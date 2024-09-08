package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	
	videoURL := "https://www.facebook.com/reel/8176615212451714"

	downloadVideo(videoURL)
}

func downloadVideo(videoURL string){
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
