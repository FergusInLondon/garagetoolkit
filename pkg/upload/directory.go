package upload

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
)

const (
	workerCount = 3
	lookupHost  = "google.com"
)

var (
	NoActiveNetworkInterface = errors.New("no active network connection available")
)

type DirectoryUploader struct {
	directory string
}

func NewDirectoryUploader(dirName string, config *Configuration) (*DirectoryUploader, error) {
	// Initiate uploading package here
	if !hasActiveNetworkInterface() {
		return nil, NoActiveNetworkInterface
	}

	if err := Connect(config); err != nil {
		return nil, err
	}

	return &DirectoryUploader{dirName}, nil
}

func (du *DirectoryUploader) Upload() error {
	uploadChan := make(chan string, workerCount*2)
	wg := &sync.WaitGroup{}

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go uploadWorker(wg, uploadChan)
	}

	files, err := ioutil.ReadDir(du.directory)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			uploadChan <- du.directory + "/" + file.Name()
		}
	}

	wg.Wait()
	return nil
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

		nBytes, err := Upload(f, fileInfo)
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

func hasActiveNetworkInterface() bool {
	_, err := net.LookupIP(lookupHost)
	return err == nil
}
