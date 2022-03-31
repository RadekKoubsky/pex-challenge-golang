package pex

import "context"

type Config struct {
	ImagesPath  string
	ResultPath  string
	MostKColors int
}

func ProcessImages(ctx context.Context, downloadedImages chan DownloadedImage, processedImages chan ProcessedImage,
	config Config, counter *int64) {
	go DownloadImages(config.ImagesPath, downloadedImages)
	go FindMostKColors(ctx, downloadedImages, processedImages, config.MostKColors)
	go WriteProcessedImagesToFile(ctx, config.ResultPath, processedImages, counter)
}
