package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"proj2/png"
	"proj2/task"
	"strings"
)

type Task = task.Task
type Image = png.Image

func RunSequential(config Config) {
	dataDirs := strings.Split(config.DataDirs, "+")
	outputPath := "../data/out/%s_%s"
	inputPath := "../data/in/%s/%s"

	effectsPathFile := "../data/effects.txt"
	effectsFile, err := os.Open(effectsPathFile)
	if err != nil {
		panic(err)
	}
	defer effectsFile.Close()

	// Get the decoder
	reader := json.NewDecoder(effectsFile)

	// Decode the json requests in the effects file
	for {
		// Read the next request from the effects file
		// If there are no more requests, break
		task := Task{}
		err := reader.Decode(&task)
		if err != nil {
			break
		}

		// Process the task
		for _, dataDir := range dataDirs {
			inPath := fmt.Sprintf(inputPath, dataDir, task.InPath)
			outPath := fmt.Sprintf(outputPath, dataDir, task.OutPath)

			// Read the input file
			img, err := png.Load(inPath)
			if err != nil {
				panic(err)
			}

			// Process the effects
			bounds := img.Bounds
			processEffectsSeq(img, task.Effects, bounds.Min.Y, bounds.Max.Y)
			img.Swap()

			// Save the output file
			err = img.Save(outPath)
			if err != nil {
				panic(err)
			}
		}
	}
}

func applyEffect(img *Image, effect string, startY, endY int) {
	switch effect {
	case "G":
		img.Grayscale(startY, endY)
	case "S":
		img.Sharpen(startY, endY)
	case "B":
		img.Blur(startY, endY)
	case "E":
		img.EdgeDetect(startY, endY)
	default:
		panic("Invalid effect")
	}
}

func processEffectsSeq(img *Image, effects []string, startY, endY int) {
	for _, effect := range effects {
		applyEffect(img, effect, startY, endY)
		img.Swap()
	}
}
