package main

import (
	"fmt"
	"strings"

	"github.com/status-im/image_resizer/common"
)

type imageDetails struct {
	SizePixel uint
	SizeFile int64
	Quality int
	FileName string
	Properties string
}

func getSourceImageName(imageName string) string {
	return common.ImageDir + imageName + ".jpg"
}

func makeOutputImageName(imageName string, size uint, i int, properties string) string {
	if properties != "" {
		properties = "_" + strings.ReplaceAll(properties, " ", "-")
	}
	return fmt.Sprintf(common.ImageDir + "%s_s-%d_q-%d%s.jpg", imageName, size, i, properties)
}

func makeImageDetails(imageName string, size uint, quality int, properties string) imageDetails {
	return imageDetails{
		SizePixel: size,
		Quality:   quality,
		Properties: properties,
		FileName:  makeOutputImageName(imageName, size, quality, properties),
	}
}
