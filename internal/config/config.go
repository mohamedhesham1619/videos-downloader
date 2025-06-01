package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"videos-downloader/internal/models"
	"videos-downloader/internal/utils"

	"github.com/fatih/color"
	"github.com/jaypipes/ghw"
)

type Config struct {

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

	pathFlag := flag.String("path", "", "path to the download directory (the default is the current directory)")

	fastFlag := flag.Bool("fast", false, "Fast mode: copies streams directly without re-encoding. Much faster but clips may start slightly early or have frozen frames at the beginning")

	flag.Parse()

	encoder := ""

	// If the fast flag is not used, select the encoder based on the GPU. If the GPU is not detected or the GPU encoder is not working, the CPU encoder will be used.

	// If the fast flag is used, the encoder will be ignored and the streams will be copied directly without re-encoding.
	if !*fastFlag {
		encoder = selectEncoder()
	} else {
		color.Cyan("Fast mode enabled: copying streams directly without re-encoding.\n")
	}

	cfg := &Config{
		DownloadPath: *pathFlag,
		Encoder:      encoder,
		IsFastMode:   *fastFlag,
	}

	// if the user provides a path flag, the downloaded videos will be saved in that directory. Otherwise, they will be saved in the "downloads" folder in the current directory.
	if cfg.DownloadPath == "" {
		err := os.MkdirAll("downloads", os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating downloads directory: %v Will use the current directory instead.\n\n", err)
			cfg.DownloadPath = ""
		} else {
			cfg.DownloadPath = "downloads"
		}

	}

	return cfg
}

func detectGpu() (string, error) {

	gpuInfo, err := ghw.GPU()

	if err != nil {
		return "", fmt.Errorf("error getting GPU info: %v", err)
	}

	gpu := ""
	gpuVendorName := strings.ToLower(gpuInfo.GraphicsCards[0].DeviceInfo.Vendor.Name)

	// Check the GPU vendor name and set the GPU variable accordingly
	switch {

	case strings.Contains(gpuVendorName, "nvidia"):
		gpu = models.NvidiaGPU
	case strings.Contains(gpuVendorName, "amd") || strings.Contains(gpuVendorName, "advanced micro devices"):
		gpu = models.AMDGPU
	case strings.Contains(gpuVendorName, "intel"):
		gpu = models.IntelGPU
	}

	return gpu, nil

}

// If the fast flag is not used, the encoder will be selected based on the GPU. If the GPU is not detected or the GPU encoder is not working, the CPU encoder will be used.
func selectEncoder() string {
	gpu, err := detectGpu()
	if err != nil || gpu == "" || !utils.TestGpuEncoder(models.GPUEncoders[gpu]) {
		color.Cyan("Could not use GPU encoder. Falling back to CPU encoder:", models.CPUEncoder, "\n")
		return models.CPUEncoder
	}

	color.Cyan("Using %s encoder: %s\n", gpu, models.GPUEncoders[gpu])
	return models.GPUEncoders[gpu]
}
