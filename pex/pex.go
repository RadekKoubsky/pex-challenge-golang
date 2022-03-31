package pex

import (
	"context"
	"github.com/EdlinOrg/prominentcolor"
	"image"
)

type DownloadedImage struct {
	url string
	img image.Image
}

type ProcessedImage struct {
	url         string
	mostKColors []prominentcolor.ColorItem
}

type Pex struct {
	ImageReader ImagesReader
	MostKColors MostKColors
	ImageWriter ImagesWriter
}

func (pex Pex) ProcessImages(ctx context.Context, counter *int64) {
	go pex.ImageReader.DownloadImages()
	go pex.MostKColors.FindMostKColors(ctx)
	go pex.ImageWriter.WriteProcessedImagesToFile(ctx, counter)
}
