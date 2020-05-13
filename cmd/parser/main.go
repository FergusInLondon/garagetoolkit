package main

import (
	"fmt"
	"io"
	"os"

	"go.fergus.london/garagetoolkit/pkg/canlog"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage:\n\t%s <filename>\n\n", os.Args[0])
		fmt.Printf("Parameters:\n\t<filename> - binary log to parse\n\n")
		return
	}

	var (
		logFile io.ReadCloser
		entry   *canlog.DecodedMessage
		err     error
	)

	if logFile, err = os.Open(os.Args[1]); err != nil {
		panic(err)
	}

	messageParser := canlog.NewParser(logFile)
	defer logFile.Close()

	entryCount := 0
	for {
		if entry, err = messageParser.Iterate(); err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}

		fmt.Printf("[ %s ]\tMessage ID: '%d'\t%v\n", entry.Time.String(), entry.Frame.ID, entry.Frame.Data)
		entryCount++
	}

	fmt.Printf("Successfully parsed '%s', with %d entries.\n", os.Args[1], entryCount)
}
