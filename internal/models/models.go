package models

type VideoRequest struct {
	Url           string
	IsClip        bool
	ClipTimeRange string // should be in the format HH:MM:SS-HH:MM:SS
}
