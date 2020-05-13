package main

import (
	"log"
	"os"
	"sync"

	"go.fergus.london/garagetoolkit/pkg/upload"
)

const poolCount = 5

var uploadConfiguration = &upload.Configuration{ /* Take from Environment? */ }

func main() {
	if err := upload.Connect(uploadConfiguration); err != nil {
		panic(err)
	}

	var (
		fileNames  = getFilesToUpload()
		wg         = &sync.WaitGroup{}
		uploadChan = make(chan string, len(fileNames))
	)

	defer func() {
		close(uploadChan)
		wg.Wait()
	}()

	wg.Add(poolCount)
	for i := 0; i < poolCount; i++ {
		go uploadWorker(wg, uploadChan)
	}

	for _, file := range fileNames {
		uploadChan <- file
	}
}

func getFilesToUpload() []string {
	// Parse Directory for files.
	return []string{}
}

func uploadWorker(wg *sync.WaitGroup, inbound chan string) {
	defer wg.Done()
	for file := range inbound {
		f, err := os.Open(file)
		if err != nil {
			log.Println("error encountered opening file. skipping.", file, err)
			continue
		}

		fileInfo, err := f.Stat()
		if err != nil {
			log.Println("invalid file info, skipping.", file, err)
			f.Close()
			continue
		}

		nBytes, err := upload.Upload(f, fileInfo)
		if err != nil {
			log.Println("failed to upload file, skipping.", file, err)
			f.Close()
			continue
		}

		log.Println("uploaded file", nBytes, file)
		f.Close()

		if err := os.Remove(file); err != nil {
			log.Println("unable to remove file", file, err)
		}
	}
}
