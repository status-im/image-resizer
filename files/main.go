package files

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"strconv"

	"github.com/status-im/image_resizer/common"
	"github.com/status-im/image_resizer/images"
)

func GetImage(fileName string) image.Image {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	return img
}

func RenderImage(img image.Image, imgDetail *images.Details) {
	out, err := os.Create(imgDetail.FileName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	o := new(jpeg.Options)
	o.Quality = imgDetail.Quality

	jpeg.Encode(out, img, o)

	fi, _ := out.Stat()
	imgDetail.SizeFile = fi.Size()
}

func MakeReadMe(imageList []string, imgDs map[string][]images.Details) {
	var txt string

	for _, imageName := range imageList {
		txt += "## " + imageName + "\n\n"
		txt += fmt.Sprintf("![Original %s image](%s)\n\n", imageName, common.ImageDir+imageName+".jpg")

		txt += "| Image | Properties | Size (px^2) | Image Quality (%) | Size (bytes) |\n"
		txt += "| :---: | ---------- | ----------: | ----------------: | -----------: |\n"

		for _, id := range imgDs[imageName] {
			txt += fmt.Sprintf("| ![](%s) | %s | %d | %d | %s |\n", id.FileName, id.Properties, id.SizePixel, id.Quality, format(id.SizeFile))
		}

		txt += "\n"
	}

	rm, err := os.Create("README.md")
	if err != nil {
		log.Fatal(err)
	}
	defer rm.Close()

	rm.WriteString(txt)
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
