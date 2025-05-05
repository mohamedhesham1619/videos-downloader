package config

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Config struct {
	UrlsFile string // the file containing the urls of the videos to download (only the file name if it's in the same directory as the program or the full path if it's in another)

	DownloadPath string // the path to the download directory (the default is the directory where the program is executed)
}

func New() *Config {

	urlsFlag := flag.String("urls", "urls.txt", "only the file name if it's in the same directory as the program or the full path if it's in another directory")

	pathFlag := flag.String("path", "", "path to the download directory (the default is the current directory)")

	flag.Parse()

	cfg := &Config{
		UrlsFile:     *urlsFlag,
		DownloadPath: *pathFlag,
	}

	// if the user provides a path flag, the downloaded videos will be saved in that directory. Otherwise, they will be saved in the "downloads" folder in the current directory.
	if cfg.DownloadPath == "" {
		err := os.MkdirAll("downloads", os.ModePerm)
		if err != nil {
			log.Fatal(fmt.Errorf("error creating downloads directory: %v", err))
		}
		cfg.DownloadPath = "downloads"
	}

	return cfg
}
