package images

import (
	"fmt"
	"strings"

	"github.com/status-im/image_resizer/common"
)

type Details struct {
	SizePixel  uint
	SizeFile   int64
	Quality    int
	FileName   string
	Properties string
}

func GetSourceName(imageName string) string {
	return common.ImageDir + imageName + ".jpg"
}

func MakeDetails(imageName string, size uint, quality int, properties string) Details {
	return Details{
		SizePixel:  size,
		Quality:    quality,
		Properties: properties,
		FileName:   makeOutputName(imageName, size, quality, properties),
	}
}

func makeOutputName(imageName string, size uint, i int, properties string) string {
	if properties != "" {
		properties = "_" + strings.ReplaceAll(properties, " ", "-")
	}
	return fmt.Sprintf(common.ImageDir+"%s_s-%d_q-%d%s.jpg", imageName, size, i, properties)
}
