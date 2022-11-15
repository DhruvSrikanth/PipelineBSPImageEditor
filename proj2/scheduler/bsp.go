package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"proj2/png"
	"strings"
	"sync"
)

// Shared context for BSP workers
type bspWorkerContext struct {
	mutex         *sync.Mutex
	cond          *sync.Cond
	imageTasks    []ImageTask
	numThreads    int
	effectIdx     int
	taskIdx       int
	threadCounter int
}

// Obtain the shared context for the BSP workers
func NewBSPContext(config Config) *bspWorkerContext {
	dataDirs := strings.Split(config.DataDirs, "+")
	outputPath := "../data/out/%s_%s"
	inputPath := "../data/in/%s/%s"

	// Create the image task pipeline
	imageTasks := make([]ImageTask, 0)
	imageTaskGeneratorBSP(dataDirs, inputPath, outputPath, &imageTasks)

	//Initialize the context
	mutex := &sync.Mutex{}
	cond := sync.NewCond(mutex)
	// Keep track of the effect
	effectIdx := 0
	// Keep track of the task
	taskIdx := 0
	// Keep track of the number of threads
	threadCounter := 0

	return &bspWorkerContext{
		mutex:         mutex,
		cond:          cond,
		imageTasks:    imageTasks,
		numThreads:    config.ThreadCount,
		effectIdx:     effectIdx,
		taskIdx:       taskIdx,
		threadCounter: threadCounter,
	}
}

// Generate tasks for the BSP workers to perform
func imageTaskGeneratorBSP(dataDirs []string, inputPath, outputPath string, imageTasks *[]ImageTask) {

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

			// Add the image task to the pipeline
			*imageTasks = append(*imageTasks, ImageTask)
		}
	}
}

// Perform the work done by a BSP worker
func RunBSPWorker(id int, ctx *bspWorkerContext) {
	for {
		// Determine work per thread for the current task
		workPerThread := ctx.imageTasks[ctx.taskIdx].Image.Bounds.Max.Y / ctx.numThreads

		// Determine the work interval for each Image Task
		var startIdx, endIdx int
		if id+1 == ctx.numThreads {
			startIdx = (ctx.numThreads - 1) * workPerThread
			endIdx = ctx.imageTasks[ctx.taskIdx].Image.Bounds.Max.Y
		} else {
			startIdx = id * workPerThread
			endIdx = startIdx + workPerThread
		}

		// Apply the effect to the image
		applyEffectBSP(ctx.imageTasks[ctx.taskIdx].Effects[ctx.effectIdx], ctx.imageTasks[ctx.taskIdx].Image, startIdx, endIdx)

		// Lock the mutex
		ctx.mutex.Lock()

		// Increment the thread counter
		ctx.threadCounter += 1

		// If all threads are done, perform swap and move on to next effect
		if ctx.threadCounter == ctx.numThreads {
			ctx.imageTasks[ctx.taskIdx].Image.Swap()

			ctx.threadCounter = 0
			ctx.effectIdx += 1

			// If all effects are done, perform swap, save image and move on to next task
			if ctx.effectIdx == len(ctx.imageTasks[ctx.taskIdx].Effects) {
				ctx.imageTasks[ctx.taskIdx].Image.Swap()
				ctx.imageTasks[ctx.taskIdx].Image.Save(ctx.imageTasks[ctx.taskIdx].OutputPath)

				ctx.effectIdx = 0
				ctx.taskIdx += 1
				// If all tasks are done, break
				if ctx.taskIdx == len(ctx.imageTasks) {
					ctx.mutex.Unlock()
					break
				} else {
					// there are still tasks to be done
					// Wake up threads
					ctx.cond.Broadcast()
				}
			} else {
				// there are still effects to be done
				// Wake up threads
				ctx.cond.Broadcast()
			}
		} else {
			// there are still threads to be finish
			// Wait for other threads to finish
			ctx.cond.Wait()
		}

		// Unlock the mutex
		ctx.mutex.Unlock()
	}
}

// Apply the effect to the image for the given interval
func applyEffectBSP(effect string, img *png.Image, startY, endY int) {
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
}
