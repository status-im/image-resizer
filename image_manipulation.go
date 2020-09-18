package main

import (
	"image"
	"image/draw"
	"log"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

func cropImage(img image.Image) image.Image {
	var sl int
	if img.Bounds().Max.X > img.Bounds().Max.Y {
		sl = img.Bounds().Max.Y
	} else {
		sl = img.Bounds().Max.X
	}

	croppedImg, err := cutter.Crop(img, cutter.Config{
		Width: sl,
		Height: sl,
		Mode: cutter.Centered,
	})
	if err != nil {
		log.Fatal(err)
	}

	return croppedImg
}

func resizeSquareImage(size uint, img image.Image) image.Image {
	return resize.Resize(size, 0, img, resize.Bilinear)
}

func circleCropImage(img image.Image, size int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.DrawMask(
		dst,
		dst.Bounds(),
		img,
		image.Point{},
		&circle{
			image.Point{X: size / 2, Y: size / 2},
			size / 2,
		},
		image.Point{},
		draw.Over,
	)
	return dst
}
