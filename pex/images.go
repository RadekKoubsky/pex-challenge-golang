package pex

import (
	"bufio"
	"context"
	"fmt"
	"github.com/EdlinOrg/prominentcolor"
	"image"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

type DownloadedImage struct {
	url string
	img image.Image
}

type ProcessedImage struct {
	url         string
	mostKColors []prominentcolor.ColorItem
}

func DownloadImages(imagesPath string, downloadedImages chan DownloadedImage) {
	urls := make(chan string)
	go readFromFile(imagesPath, urls)

	var wg sync.WaitGroup
	const goroutines = 10
	wg.Add(goroutines)
	// limit number of workers to not cause too many open connections when downloading images
	for i := 0; i < goroutines; i++ {
		go func() {
			for url := range urls {
				DownloadImage(url, downloadedImages)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func readFromFile(imagesPath string, urls chan string) {
	imagesFile, err := os.Open(imagesPath)
	if err != nil {
		log.Fatalln("Failed to open file", imagesFile)
	}
	defer imagesFile.Close()

	scanner := bufio.NewScanner(imagesFile)
	for scanner.Scan() {
		urls <- scanner.Text()
	}
	close(urls)
}

func DownloadImage(url string, downloadedImages chan DownloadedImage) {
	resp, getError := http.Get(url)
	if getError != nil {
		log.Println("Failed to download image", getError)
		return
	}
	defer resp.Body.Close()

	decodedImg, _, decodeErr := image.Decode(resp.Body)
	if decodeErr != nil {
		log.Println("Failed to decode image", url, decodeErr)
		return
	} else {
		downloadedImages <- DownloadedImage{
			url: url,
			img: decodedImg}
	}
}

func FindMostKColors(ctx context.Context, downloadedImages chan DownloadedImage, processedImages chan ProcessedImage,
	mostKColors int) {
	for {
		select {
		case <-ctx.Done():
			log.Println("FindMostKColors done")
			return
		case downloadedImage := <-downloadedImages:
			go FindMostKColor(downloadedImage, processedImages, mostKColors)
		}
	}
}

func FindMostKColor(downloadedImage DownloadedImage, processedImages chan ProcessedImage, mostKColors int) {
	colors, err := prominentcolor.KmeansWithAll(mostKColors, downloadedImage.img, prominentcolor.ArgumentDefault,
		prominentcolor.DefaultSize, prominentcolor.GetDefaultMasks())
	if err != nil {
		log.Println("Failed to process image", downloadedImage.url, err)
		return
	} else {
		processedImages <- ProcessedImage{
			url:         downloadedImage.url,
			mostKColors: colors,
		}
	}
}

func WriteProcessedImagesToFile(ctx context.Context, resultPath string, processedImages chan ProcessedImage, counter *int64) {
	file, err := os.Create(resultPath)
	if err != nil {
		log.Fatalln("Could not create file", err)
	}
	defer file.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("WriteProcessedImagesToFile done")
			return
		case processedImage := <-processedImages:
			_, err := file.WriteString(fmt.Sprintf("%v,%v\n", processedImage.url, printColors(processedImage.mostKColors)))
			if err != nil {
				log.Println("Failed to write a line to result file", err)
			}
			atomic.AddInt64(counter, 1)
		}
	}
}

func printColors(colors []prominentcolor.ColorItem) string {
	colorsHex := make([]string, 0, len(colors))
	for _, color := range colors {
		colorsHex = append(colorsHex, fmt.Sprintf("#%v", color.AsString()))
	}
	return strings.Join(colorsHex, ",")
}
