package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"go.fergus.london/canlog/pkg/canlog"
)

var (
	inboundSignal = make(chan os.Signal)
	logger        *canlog.Logger
)

func getArgument(idx int, def string) string {
	if len(os.Args) < (idx + 1) {
		return def
	}

	return os.Args[idx]
}

func main() {
	signal.Notify(inboundSignal, syscall.SIGTERM, syscall.SIGABRT)
	go handleSigTerm()

	logFileName := getArgument(2, "canlog.bin")
	logFile := CreateLogFile(logFileName)
	log.Println("opened log file for writing", logFile.FilePath())

	canbusInterface := getArgument(1, "can0")
	log.Println("connecting to canbus interface", canbusInterface)

	var err error
	logger, err = canlog.NewLogger(canbusInterface, logFile.File)
	if err != nil {
		panic(err)
	}

	defer func() {
		log.Println("stopped listening for canbus frames")
		logFile.Finish()
		log.Println("closed log file")
	}()

	go sdNotifier()
	log.Println(logger.Run())
}

// notify systemd that we're ready
// notify systemd that we're still alive.
//
// For liveness we rely on a watchdog, via systemd. If we fail to respond then
//
//
//
func sdNotifier() {
	daemon.SdNotify(false, daemon.SdNotifyReady)

	for {
		daemon.SdNotify(false, daemon.SdNotifyWatchdog)
		time.Sleep(3 * time.Second)
	}
}

// When systemd shutdowns a process, it sends a SIGTERM signal to it. In the
// event that this signal is ignored, it will usually follow-up with a SIGKILL.
// Additionally, systemd will trigger SIGABRT in the event that the watchdog
// fails to fire.
//
// Here we trap both signals, and use it for clean-up: i.e stopping listeners
// and closing log files.
//
// @see https://www.freedesktop.org/software/systemd/man/systemd.service.html
func handleSigTerm() {
	sig := <-inboundSignal
	log.Println("received signal from operating system", sig.String())
	log.Println("stopping canbus listener", logger.Stop())
}
