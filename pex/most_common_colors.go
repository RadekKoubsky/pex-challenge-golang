package pex

import (
	"context"
	"github.com/EdlinOrg/prominentcolor"
	"log"
)

type MostKColors struct {
	DownloadedImages chan DownloadedImage
	ProcessedImages  chan ProcessedImage
	MostK            int
}

func (mostKColors MostKColors) FindMostKColors(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("FindMostKColors done")
			return
		case downloadedImage := <-mostKColors.DownloadedImages:
			go mostKColors.FindMostKColor(downloadedImage)
		}
	}
}

func (mostKColors MostKColors) FindMostKColor(downloadedImage DownloadedImage) {
	colors, err := prominentcolor.KmeansWithAll(mostKColors.MostK, downloadedImage.img, prominentcolor.ArgumentDefault,
		prominentcolor.DefaultSize, prominentcolor.GetDefaultMasks())
	if err != nil {
		log.Println("Failed to process image", downloadedImage.url, err)
		return
	} else {
		mostKColors.ProcessedImages <- ProcessedImage{
			url:         downloadedImage.url,
			mostKColors: colors,
		}
	}
}
