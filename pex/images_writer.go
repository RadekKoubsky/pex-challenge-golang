package pex

import (
	"context"
	"fmt"
	"github.com/EdlinOrg/prominentcolor"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

type ImagesWriter struct {
	ResultPath      string
	ProcessedImages chan ProcessedImage
}

func (imagesWriter ImagesWriter) WriteProcessedImagesToFile(ctx context.Context, counter *int64) {
	file, err := os.Create(imagesWriter.ResultPath)
	if err != nil {
		log.Fatalln("Could not create file", err)
	}
	defer file.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("WriteProcessedImagesToFile done")
			return
		case processedImage := <-imagesWriter.ProcessedImages:
			_, err := file.WriteString(fmt.Sprintf("%v,%v\n", processedImage.url, imagesWriter.printColors(processedImage.mostKColors)))
			if err != nil {
				log.Println("Failed to write a line to result file", err)
			}
			atomic.AddInt64(counter, 1)
		}
	}
}

func (imagesWriter ImagesWriter) printColors(colors []prominentcolor.ColorItem) string {
	colorsHex := make([]string, 0, len(colors))
	for _, color := range colors {
		colorsHex = append(colorsHex, fmt.Sprintf("#%v", color.AsString()))
	}
	return strings.Join(colorsHex, ",")
}
