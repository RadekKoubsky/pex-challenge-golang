package main

import (
	"context"
	"log"
	"rkoubsky.com/pex-challenge-golang/pex"
	"sync/atomic"
	"time"
)

func main() {
	config := pex.Config{
		ImagesPath:  "images.txt",
		ResultPath:  "prevalent_colors.csv",
		MostKColors: 3,
	}
	downloadedImages := make(chan pex.DownloadedImage)
	processedImages := make(chan pex.ProcessedImage)
	ctx, cancelFunc := context.WithCancel(context.Background())
	var counter int64
	start := time.Now()
	pex.ProcessImages(ctx, downloadedImages, processedImages, config, &counter)
	wait(&counter)
	elapsed := time.Now().Sub(start)
	log.Printf("Processing images took %s\n", elapsed)
	cancelFunc()
}

func wait(counter *int64) {
	for atomic.LoadInt64(counter) < 1000 {
		time.Sleep(time.Second * 1)
	}
}
