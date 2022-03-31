package main

import (
	"context"
	"log"
	"rkoubsky.com/pex-challenge-golang/pex"
	"sync/atomic"
	"time"
)

func main() {
	downloadedImages := make(chan pex.DownloadedImage)
	processedImages := make(chan pex.ProcessedImage)
	ctx, cancelFunc := context.WithCancel(context.Background())
	const mostK = 3
	pexChallenge := pex.Pex{
		pex.ImagesReader{downloadedImages, "images.txt"},
		pex.MostKColors{downloadedImages, processedImages, mostK},
		pex.ImagesWriter{"prevalent_colors.csv", processedImages},
	}
	var counter int64
	start := time.Now()
	go pexChallenge.ProcessImages(ctx, &counter)
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
