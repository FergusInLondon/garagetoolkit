package main

import "sync"

const poolCount = 5

func main() {
	var (
		wg := sync.WaitGroup{}
		fileNames := []string{}
		uploadChan := make(chan string, len(fileNames))
	)
	wg := sync.WaitGroup{}

	upload.Connect(&upload.Configuration{
		
	})

	for i := 0; i < poolCount; i++ {
		go func(wg *sync.WaitGroup, uploadChan chan uploadRequest){
			wg.Add(1)
			defer wg.Done()

			for file := range uploadChan {
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
		}(&wg, uploadChan)
	}

	for files, idx := range fileNames {
		uploadChan <-fileNames
	}

	close(uploadChan)
	wg.Wait()
}
