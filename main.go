package main

import (
	"github.com/status-im/image_resizer/files"
	"github.com/status-im/image_resizer/images"
)

var (
	imageList = []string{
		"elephant",
		"frog",
		"goat",
		"mars",
		"psychedelic",
		"rainbow",
		"romanian-flag",
		"tormund",
		"woman",
	}

	sizes = []uint{
		80,
		240,
	}
)

func main() {
	imgDs := make(map[string][]images.Details)

	for _, imageName := range imageList {
		img := files.GetImage(images.GetSourceName(imageName))
		croppedImg := images.Crop(img)

		for _, size := range sizes {
			for i := 1; i < 11; i++ {

				ri := images.ResizeSquare(size, croppedImg)
				id := images.MakeDetails(imageName, size, i*10, "")
				files.RenderImage(ri, &id)
				imgDs[imageName] = append(imgDs[imageName], id)

				precci := images.CropCircle(ri, int(size))
				precid := images.MakeDetails(imageName, size, i*10, "pre-render circle crop")
				files.RenderImage(precci, &precid)
				imgDs[imageName] = append(imgDs[imageName], precid)

				li := files.GetImage(id.FileName)
				postcci := images.CropCircle(li, int(size))
				postcid := images.MakeDetails(imageName, size, i*10, "post-render circle crop")
				files.RenderImage(postcci, &postcid)
				imgDs[imageName] = append(imgDs[imageName], postcid)
			}
		}
	}

	files.MakeReadMe(imageList, imgDs)
}
