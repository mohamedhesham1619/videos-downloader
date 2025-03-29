package main

import (
	"fmt"
	"sync"
	"videos-downloader/internal/config"
	"videos-downloader/internal/downloader"
	"videos-downloader/internal/utils"
)

func main() {
	cfg := config.New()
	downloader := downloader.New(cfg)

	urls, err := utils.ReadUrlsFromFile(cfg.UrlsFile)

	if err != nil {
		fmt.Println(err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(urls))

	// start downloading videos concurrently
	for _, url := range urls {
		go func() {
			videoRequest := utils.ParseVideoRequest(url)
			err := downloader.Download(videoRequest)

			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}()
	}

	// ensure all goroutines complete
	wg.Wait()
}