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

	log.Printf("connected to minio storage bucket")

	var (
		fileNames  = getFilesToUpload()
		wg         = &sync.WaitGroup{}
		uploadChan = make(chan string, len(fileNames))
	)

	defer func() {
		close(uploadChan)
		log.Printf("waiting for uploads to complete.")
		wg.Wait()
	}()

	wg.Add(poolCount)
	for i := 0; i < poolCount; i++ {
		log.Printf("starting upload worker (%d)", i)
		go uploadWorker(wg, uploadChan)
	}

	log.Printf("got %d files to process", len(fileNames))
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
		log.Printf("[ %s ] processing file", file)

		f, err := os.Open(file)
		if err != nil {
			log.Printf("[ %s ] unable to open file - %v", file, err)
			continue
		}

		fileInfo, err := f.Stat()
		if err != nil {
			log.Printf("[ %s ] invalid file info - %v", file, err)
			f.Close()
			continue
		}

		nBytes, err := upload.Upload(f, fileInfo)
		if err != nil {
			log.Printf("[ %s ] failed to upload file - %v", file, err)
			f.Close()
			continue
		}

		log.Printf("[ %s ] uploaded file successfully (%d bytes)", file, nBytes)
		f.Close()

		if err := os.Remove(file); err != nil {
			log.Println("unable to remove file", file, err)
			continue
		}

		log.Printf("[ %s ] successfully uploaded file and removed local version", file)
	}
}
