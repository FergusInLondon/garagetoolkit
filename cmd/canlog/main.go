package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"go.fergus.london/garagetoolkit/pkg/canlog"
	"go.fergus.london/garagetoolkit/pkg/upload"
)

const (
	defaultWatchdogInterval = "10"
	bucketName       = "canlog"
	bucketLocation   = "us-east-1"
)

var (
	inboundSignal = make(chan os.Signal)
	logger        *canlog.Logger
	ctx           context.Context
	ctxCancel     context.CancelFunc
)

func getArgument(idx int, def string) string {
	if len(os.Args) < (idx + 1) {
		return def
	}

	return os.Args[idx]
}

func getEnvironmentVariable(key, def) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
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
	<-ctx.Done()
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

	watchdogInterval, err := strconv.Atoi(
		getEnvironmentVariable("WATCHDOG_INTERVAL", defaultWatchdogInterval)
	)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After((watchdogInterval / 2) * time.Second):
			daemon.SdNotify(false, daemon.SdNotifyWatchdog)
		}
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

	if err := handleUploads(); err != nil {
		log.Fatal("uploading failed", err)
	}
}

//
func handleUploads() error {
	defer ctxCancel()

	uploadConfig = upload.GetConfig()
	uploadConfig.BucketName = getEnvironmentVariable("BUCKET_NAME", bucketName)
	uploadConfig.BucketLocation = getEnvironmentVariable("BUCKET_LOCATION", bucketLocation)

	uploader, err := upload.NewDirectoryUploader(logsDirectory, uploadConfig)
	if err != nil {
		return err
	}

	return uploader.Upload()
}
