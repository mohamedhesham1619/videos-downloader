package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"videos-downloader/internal/config"
	"videos-downloader/internal/downloader"
	"videos-downloader/internal/utils"

	"github.com/fatih/color"
)

func main() {
	cfg := config.New()
	downloader := downloader.New(cfg)

	urls, err := utils.ReadUrlsFromFile("urls.txt")

	if err != nil {
		log.Fatal("Error reading urls from file \n", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(urls))

	// start downloading videos concurrently
	for _, url := range urls {
		go func() {
			videoRequest := utils.ParseVideoRequest(url)
			err := downloader.Download(videoRequest)

			if err != nil {
				fmt.Printf("%s (%s): %s\n\n", color.RedString("Error downloading video"), url, color.RedString(err.Error()))
			}
			wg.Done()
		}()
	}

	// ensure all goroutines complete
	wg.Wait()

	// Wait for user input before exiting
	fmt.Print("\nAll downloads completed. Press Enter to exit...\n")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
