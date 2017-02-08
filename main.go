package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/masahide/twitter-stream-aggregator/twitter"
)

func checkFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readTwitterStream(ctx context.Context) chan string {
	urlCh := make(chan string)
	tw, err := twitter.NewTwitter(ctx)
	checkFatal(err)
	resp, err := tw.StreamFilter(`?track=pig`)
	checkFatal(err)
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	return urlCh
}

func downloadWorker(resCh, inCh chan string) {
	for url := range inCh {
		resCh <- url
	}
}
func downloadImage(inCh chan string, maxWorker int) chan string {
	resCh := make(chan string)
	for i := 0; i < maxWorker; i++ {
		go downloadWorker(resCh, inCh)
	}
	return resCh

}

func googleVisionWorker(resCh, inCh chan string) {
	for file := range inCh {
		resCh <- file
	}
}
func googleVision(inCh chan string, maxWorker int) chan string {
	resCh := make(chan string)
	for i := 0; i < maxWorker; i++ {
		go googleVisionWorker(resCh, inCh)
	}
	return resCh

}
func ibmRecognitionWorker(resCh, inCh chan string) {
	for file := range inCh {
		resCh <- file
	}
}
func ibmRecognition(inCh chan string, maxWorker int) chan string {
	resCh := make(chan string)
	for i := 0; i < maxWorker; i++ {
		go ibmRecognitionWorker(resCh, inCh)
	}
	return resCh

}

const (
	maxWorker = 10
)

func main() {

	ctx := context.Background()

	// twitter ストリームを読む
	urlCh := readTwitterStream(ctx)

	// ダウンロードする
	fileCh := downloadImage(urlCh, maxWorker)

	// google vision api
	filterCh := googleVision(fileCh, maxWorker)

	// IBM Watson Visual Recognition
	resultCh := ibmRecognition(filterCh, maxWorker)

	//結果の表示
	for file := range resultCh {
		log.Print(file)
	}
}
