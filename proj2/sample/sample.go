package main

import (
	"os"
	"proj2/png"
)

func main() {

	/******
		The following code shows you how to work with PNG files in Golang.
	******/

	//Assumes the user specifies a file as the first argument
	filePath := os.Args[1]

	//Loads the png image and returns the image or an error
	pngImg, err := png.Load(filePath)

	if err != nil {
		panic(err)
	}

	//Performs filtering on the image
	pngImg.Grayscale(pngImg.Bounds.Min.Y, pngImg.Bounds.Max.Y)
	//Saves the image to a new file
	err = pngImg.Save("test_gray.png")
	// do the same for sharpen
	pngImg.Sharpen(pngImg.Bounds.Min.Y, pngImg.Bounds.Max.Y)
	err = pngImg.Save("test_sharpen.png")
	// do the same for blur
	pngImg.Blur(pngImg.Bounds.Min.Y, pngImg.Bounds.Max.Y)
	err = pngImg.Save("test_blur.png")
	// do the same for edge detect
	pngImg.EdgeDetect(pngImg.Bounds.Min.Y, pngImg.Bounds.Max.Y)
	// Print every pixel in the image
	err = pngImg.Save("test_edge.png")

	//Checks to see if there were any errors when saving.
	if err != nil {
		panic(err)
	}

}
