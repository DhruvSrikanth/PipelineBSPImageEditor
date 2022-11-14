package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"proj2/png"
	"proj2/task"
	"strings"
)

type ImageTask = task.ImageTask

func RunPipeline(config Config) {
	dataDirs := strings.Split(config.DataDirs, "+")
	outputPath := "../data/out/%s_%s"
	inputPath := "../data/in/%s/%s"

	// Create a wait group
	imageTasksDone := make(chan bool, config.ThreadCount)
	imageTasksChan := make(chan ImageTask)

	for i := 0; i < config.ThreadCount; i++ {
		go imageWorker(imageTasksChan, imageTasksDone, config)
	}
	// Create the channel
	imageTaskGenerator(dataDirs, inputPath, outputPath, imageTasksChan)

	// Wait for all the go routines to finish
	for i := 0; i < config.ThreadCount; i++ {
		<-imageTasksDone
	}

}

func imageTaskGenerator(dataDirs []string, inputPath, outputPath string, imageTaskChan chan<- ImageTask) {
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

			// Create the Image Task
			ImageTask := ImageTask{
				Image:      img,
				OutputPath: outPath,
				Effects:    task.Effects,
			}

			// Send the Image Task to the channel
			imageTaskChan <- ImageTask
		}
	}
	close(imageTaskChan)
}

func imageWorker(imageTasksChan <-chan ImageTask, imageTasksDone chan<- bool, config Config) {
	// Process the tasks
	for {
		imageTask, additionalWork := <-imageTasksChan
		if additionalWork {
			// Process the task
			processImageTask(imageTask, config)
		} else {
			// No more work to be done
			imageTasksDone <- true
			return
		}
	}
}

func processImageTask(imageTask ImageTask, config Config) {
	// Apply the effects
	effectDone := make(chan bool, config.ThreadCount)
	for _, effect := range imageTask.Effects {
		// Process the effect
		processEffect(effect, imageTask.Image, effectDone, config)

		// Wait for the effect to finish
		for i := 0; i < config.ThreadCount; i++ {
			<-effectDone
		}

		// Get the valid image
		imageTask.Image.Swap()
	}
	close(effectDone)

	imageTask.Image.Swap()

	// Save the output file
	err := imageTask.Image.Save(imageTask.OutputPath)
	if err != nil {
		panic(err)
	}
}

func processEffect(effect string, img *png.Image, effectDone chan<- bool, config Config) {
	// Define the work per thread
	workPerThread := img.Bounds.Max.Y / config.ThreadCount
	for threadIdx := 0; threadIdx < config.ThreadCount-1; threadIdx++ {
		startIdx := threadIdx * workPerThread
		endIdx := startIdx + workPerThread

		// Apply the effect
		go applyEffect(effect, img, startIdx, endIdx, effectDone)
	}

	// Apply the effect for the last thread
	startIdx := (config.ThreadCount - 1) * workPerThread
	endIdx := img.Bounds.Max.Y

	// Apply the effect
	go applyEffect(effect, img, startIdx, endIdx, effectDone)
}

func applyEffect(effect string, img *png.Image, startY, endY int, effectDone chan<- bool) {
	// Apply the effect
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
	effectDone <- true
}
