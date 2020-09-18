package main

import (
	"fmt"
	"strings"
)

type imageDetails struct {
	SizePixel uint
	SizeFile int64
	Quality int
	FileName string
	Properties string
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
