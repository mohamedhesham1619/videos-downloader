package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"videos-downloader/internal/models"
	"videos-downloader/internal/utils"

	"github.com/jaypipes/ghw"
)

type Config struct {

	// the file containing the urls of the videos to download (only the file name if it's in the same directory as the program or the full path if it's in another)
	UrlsFile string

	// the path to the download directory (the default is the directory where the program is executed)
	DownloadPath string

	// if true, the downloader will use fast mode (copies streams directly without re-encoding)
	IsFastMode bool

	// Encoder specifies the FFmpeg video encoder to use.
	// For normal mode:
	// - Uses GPU encoder if available
	// - Falls back to CPU encoding (libx264) if GPU encoder is not available
	// For fast mode (fast flag is used):
	// - Ignores the encoder and copy the streams directly without re-encoding
	Encoder string
}

func New() *Config {

	urlsFlag := flag.String("urls", "urls.txt", "only the file name if it's in the same directory as the program or the full path if it's in another directory")

	pathFlag := flag.String("path", "", "path to the download directory (the default is the current directory)")

	fastFlag := flag.Bool("fast", false, "Fast mode: copies streams directly without re-encoding. Much faster but clips may start slightly early or have frozen frames at the beginning")

	flag.Parse()

	gpu := ""
	encoder := ""

	if !*fastFlag {
		// get the GPU information
		if gpuInfo, err := ghw.GPU(); err == nil {

			gpuVendorName := gpuInfo.GraphicsCards[0].DeviceInfo.Vendor.Name

			// Check the GPU vendor name and set the GPU variable accordingly
			switch {

			case strings.Contains(strings.ToLower(gpuVendorName), "nvidia"):
				gpu = models.NvidiaGPU
			case strings.Contains(strings.ToLower(gpuVendorName), "amd") || strings.Contains(strings.ToLower(gpuVendorName), "advanced micro devices"):
				gpu = models.AMDGPU
			case strings.Contains(strings.ToLower(gpuVendorName), "intel"):
				gpu = models.IntelGPU
			}

			// Check if the GPU encoder is working
			if gpu != "" {
				encoder = models.GPUEncoders[gpu]
				if utils.TestGpuEncoder(encoder) {
					fmt.Println("Using GPU encoder:", encoder)
				} else {
					fmt.Printf("GPU encoder %s is not working. Falling back to CPU encoder %s\n", encoder, models.CPUEncoder)

					encoder = models.CPUEncoder
				}
			}
		} else {
			fmt.Println("Error getting GPU info:", err)
		}
	}
	cfg := &Config{
		UrlsFile:     *urlsFlag,
		DownloadPath: *pathFlag,
		Encoder:      encoder,
		IsFastMode:   *fastFlag,
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
