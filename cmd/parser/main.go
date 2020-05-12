package main

import (
	"fmt"
	"os"

	"go.fergus.london/canlog/pkg/logger"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage:\n\t%s <filename>\n\n", os.Args[0])
		fmt.Printf("Parameters:\n\t<filename> - binary log to parse\n\n")
		return
	}

	logFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	messageDecoder := logger.NewCanMessageDecoder(logFile)

	for {
		entry, err := messageDecoder.Iterate()
		if entry == nil || err != nil {
			fmt.Println("error encountered", err)
			break
		}

		fmt.Println(entry)
	}

	logFile.Close()
}
