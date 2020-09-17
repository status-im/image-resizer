package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"image/jpeg"
	"log"
	"os"
)

var (
	imageDir = "images/"
	images = []string{
		"goat",
		"psychedelic",
		"rainbow",
		"romanian-flag",
		"tormund",
		"woman",
	}

	sizes = []uint{
		64,
		256,
		512,
	}
)

func main() {
	for _, imageName := range images {
		file, err := os.Open(imageDir + imageName + ".jpg")
		if err != nil {
			log.Fatal(err)
		}

		// decode jpeg into imageName.Image
		img, err := jpeg.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()

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

		for i := 1; i < 11; i++ {
			for _, size := range sizes {
				n := fmt.Sprintf(imageDir + "%s_s-%d_q-%d.jpg", imageName, size, i*10)
				m := resize.Resize(size, 0, croppedImg, resize.Bilinear)
				out, err := os.Create(n)
				if err != nil {
					log.Fatal(err)
				}
				defer out.Close()

				o := new(jpeg.Options)
				o.Quality = i * 10

				// write new imageName to file
				jpeg.Encode(out, m, o)
			}
		}
	}
}