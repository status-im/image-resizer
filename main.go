package main

var (
	images = []string{
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
	imgDs := make(map[string][]imageDetails)

	for _, imageName := range images {
		img := getImage(getSourceImageName(imageName))
		croppedImg := cropImage(img)

		for _, size := range sizes {
			for i := 1; i < 11; i++ {

				ri := resizeSquareImage(size, croppedImg)
				id := makeImageDetails(imageName, size, i*10, "")
				outputImage(ri, &id)
				imgDs[imageName] = append(imgDs[imageName], id)

				precci := circleCropImage(ri, int(size))
				precid := makeImageDetails(imageName, size, i*10, "pre-render circle crop")
				outputImage(precci, &precid)
				imgDs[imageName] = append(imgDs[imageName], precid)

				li := getImage(id.FileName)
				postcci := circleCropImage(li, int(size))
				postcid := makeImageDetails(imageName, size, i*10, "post-render circle crop")
				outputImage(postcci, &postcid)
				imgDs[imageName] = append(imgDs[imageName], postcid)
			}
		}
	}

	makeReadMe(imgDs)
}
