package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
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
		80,
		240,
	}
)

type imageDetails struct {
	SizePixel uint
	SizeFile int64
	Quality int
	FileName string
	Properties string
}

func main() {
	imgDs := make(map[string][]imageDetails)

	for _, imageName := range images {
		img := getImage(getSourceImageName(imageName))
		croppedImg := cropImage(img)

		for _, size := range sizes {
			for i := 1; i < 11; i++ {

				ri := resize.Resize(size, 0, croppedImg, resize.Bilinear)
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

func getSourceImageName(imageName string) string {
	return imageDir + imageName + ".jpg"
}

func makeOutputImageName(imageName string, size uint, i int, properties string) string {
	if properties != "" {
		properties = "_" + strings.ReplaceAll(properties, " ", "-")
	}
	return fmt.Sprintf(imageDir + "%s_s-%d_q-%d%s.jpg", imageName, size, i, properties)
}

func makeImageDetails(imageName string, size uint, quality int, properties string) imageDetails {
	return imageDetails{
		SizePixel: size,
		Quality:   quality,
		Properties: properties,
		FileName:  makeOutputImageName(imageName, size, quality, properties),
	}
}

func getImage(fileName string) image.Image {
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

func outputImage(img image.Image, imgDetail *imageDetails) {
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

func makeReadMe(imgDs map[string][]imageDetails) {
	var txt string

	for _, imageName := range images {
		fn := imageDir + imageName + ".jpg"

		txt += "## " + imageName + "\n\n"
		txt += fmt.Sprintf("![Original %s image](%s)\n\n", imageName, fn)

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

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}
