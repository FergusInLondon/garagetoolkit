package main

import (
	"log"
	"os"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"go.fergus.london/canlog/pkg/canbus"
	"go.fergus.london/canlog/pkg/logger"
)

func getArgument(idx int, def string) string {
	if len(os.Args) < (idx + 1) {
		return def
	}

	return os.Args[idx]
}

func main() {
	interfaceName := getArgument(1, "can0")
	logFilename := getArgument(2, "canlog.bin")

	log.Printf("attempting to connect to interface '%s', logging to '%s'.",
		interfaceName, logFilename)

	if err := canbus.Connect(interfaceName); err != nil {
		log.Fatalf("unable to connect to can interface: %v", err)
	}

	logFile, err := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("unable to open log file: %v", err)
	}

	logListener := logger.NewCanMessageLogger(logFile)
	canbus.RegisterListener(logListener)

	go func() {
		<-time.After(5 * time.Second)
		log.Println(canbus.Stop())
	}()

	var (
		canlogError    error
		listenComplete = make(chan struct{})
	)

	go func(done chan struct{}) {
		defer close(done)
		canlogError = canbus.Run()
	}(listenComplete)

	go systemd_monitor()

	<-listenComplete
	logFile.Close()
	log.Fatal(canlogError)
}

// notify systemd that we're ready
// notify systemd that we're still alive.
func systemd_monitor() {
	daemon.SdNotify(false, daemon.SdNotifyReady)
}
