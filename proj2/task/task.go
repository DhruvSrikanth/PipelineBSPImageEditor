package task

import (
	"proj2/png"
)

type Image = png.Image

type Task struct {
	InPath  string   `json:"inPath"`
	OutPath string   `json:"outPath"`
	Effects []string `json:"effects"`
}

type ImageTask struct {
	Image      *Image
	OutputPath string
	Effects    []string
}
