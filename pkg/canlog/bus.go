package canlog

import (
	"encoding/gob"
	"io"
	"log"
	"time"

	"github.com/brutella/can"
)

// Message contains a raw CAN frame, as well as a timestamp corresponding to
// when it was *processed for logging*.
type Message struct {
	Time time.Time
	Raw  []byte
}

// Logger contains the CAN bus connection, access to the underlying log file,
// as well as a gob encoder for space-efficient persistence. Note that none of
// these properties require exportation.
type Logger struct {
	busConnection *can.Bus
	logFile       io.WriteCloser
	gobEncoder    *gob.Encoder
}

// NewLogger provides a Logger complete with configured CAN Bus interface,
// and associated logging subscriber.
func NewLogger(canInterface string, logFile io.WriteCloser) (*Logger, error) {
	bus, err := can.NewBusForInterfaceWithName(canInterface)
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		busConnection: bus,
		logFile:       logFile,
		gobEncoder:    gob.NewEncoder(logFile),
	}

	logger.busConnection.SubscribeFunc(logger.log)
	return logger, nil
}

// Run connects to the CAN bus, awaiting messages. This is a blocking function.
func (l *Logger) Run() error {
	return l.busConnection.ConnectAndPublish()
}

// Stop disconnects from the CAN bus, releasing any calls to `Run()`.
func (l *Logger) Stop() error {
	return l.busConnection.Disconnect()
}

func (l *Logger) log(frame can.Frame) {
	frameBytes, err := can.Marshal(frame)
	if err != nil {
		log.Println("skipping CAN frame - error encountered", err)
	}

	l.gobEncoder.Encode(&Message{
		Time: time.Now(),
		Raw:  frameBytes,
	})
}
