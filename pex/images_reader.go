package pex

import (
	"bufio"
	"image"
	"log"
	"net/http"
	"os"
	"sync"
)

type ImagesReader struct {
	DownloadedImages chan DownloadedImage
	ImagesPath       string
}

func (imagesReader ImagesReader) DownloadImages() {
	urls := make(chan string)
	go imagesReader.readFromFile(urls)

	var wg sync.WaitGroup
	const goroutines = 10
	wg.Add(goroutines)
	// limit number of workers to not cause too many open connections when downloading images
	for i := 0; i < goroutines; i++ {
		go func() {
			for url := range urls {
				imagesReader.DownloadImage(url)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func (imagesReader ImagesReader) readFromFile(urls chan string) {
	imagesFile, err := os.Open(imagesReader.ImagesPath)
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

func (imagesReader ImagesReader) DownloadImage(url string) {
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
		imagesReader.DownloadedImages <- DownloadedImage{
			url: url,
			img: decodedImg}
	}
}
