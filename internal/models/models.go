package models

type VideoRequest struct {
	Url           string
	IsClip        bool
	ClipTimeRange string // should be in the format HH:MM:SS-HH:MM:SS
}

const (
	// GPU vendors
	NvidiaGPU = "nvidia"
	AMDGPU    = "amd"
	IntelGPU  = "intel"

	// Default CPU encoder
	CPUEncoder = "libx264"
)

// GPUEncoders maps GPU names to their corresponding encoder names
var GPUEncoders = map[string]string{
	NvidiaGPU: "h264_nvenc",
	AMDGPU:    "h264_amf",
	IntelGPU:  "h264_qsv",
}
