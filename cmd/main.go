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
				// if the download fails and fast mode is not enabled, try to download the video without re-encoding
				if !cfg.IsFastMode {
					fmt.Printf("%s failed to download video(%s)\n Trying again without re-encoding\n", color.RedString("Error:"), url)
					cfg.IsFastMode = true
					err = downloader.Download(videoRequest)
					if err != nil {
						fmt.Printf("%s %v\n", color.RedString("Failed to download video even in fast mode:"), err)
					} else {
						fmt.Printf("%s Downloaded successfully in fast mode.\n", color.GreenString("Success:"))
					}
				}
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
