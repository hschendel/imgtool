package main

import (
	"image"
)

type flipY struct {

}

func (c *flipY) Execute(img image.Image) image.Image {
	mirrorImg := image.NewNRGBA(img.Bounds())
	mx := img.Bounds().Max.X
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		mx--
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			mirrorImg.Set(mx, y, img.At(x, y))
		}
	}
	return mirrorImg
}

func (c *flipY) Explanation() string {
	return "Flip the image on the y axis"
}

