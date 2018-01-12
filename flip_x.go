package main

import (
	"image"
)

type flipX struct {

}

func (c *flipX) Execute(img image.Image) image.Image {
	mirrorImg := image.NewNRGBA(img.Bounds())
	my := img.Bounds().Max.Y
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		my--
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			mirrorImg.Set(x, my, img.At(x, y))
		}
	}
	return mirrorImg
}

func (c *flipX) Explanation() string {
	return "Flip the image on the x axis"
}

