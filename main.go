package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"image/jpeg"
	"log"
	"os"
	"strconv"
)

var (
	imageDir = "images/"
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
		64,
		128,
		256,
		512,
	}

	fSizes = make(map[string]map[int]map[uint]int64)
)

func main() {
	for _, imageName := range images {
		if fSizes[imageName] == nil {
			fSizes[imageName] = make(map[int]map[uint]int64)
		}

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

				fi, _ := out.Stat()
				if fSizes[imageName][i*10] == nil {
					fSizes[imageName][i*10] = make(map[uint]int64)
				}
				fSizes[imageName][i*10][size] = fi.Size()
			}
		}
	}

	// ---

	rmc := makeReadMeContent()

	rm, err := os.Create("README.md")
	if err != nil {
		log.Fatal(err)
	}
	defer rm.Close()

	rm.WriteString(rmc)
}

func makeReadMeContent() string {
	var txt string

	for _, imageName := range images {
		fn := imageDir + imageName + ".jpg"

		txt += "## " + imageName + "\n\n"
		txt += fmt.Sprintf("![Original %s image](%s)\n\n", imageName, fn)

		txt += "| Image | Size (px^2) | Image Quality (%) | Size (bytes) |\n"
		txt += "| :---: | ----------: | ----------------: | -----------: |\n"

		for _, size := range sizes {
			for i := 1; i < 11; i++ {
				rfn := fmt.Sprintf(imageDir + "%s_s-%d_q-%d.jpg", imageName, size, i*10)

				txt += fmt.Sprintf("| ![](%s) | %d | %d | %s |\n", rfn, size, i*10, format(fSizes[imageName][i*10][size]))
			}
		}

		txt += "\n"
	}

	return txt
}

func format(n int64) string {
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}