package logger

import (
	"encoding/gob"
	"io"
	"time"

	"github.com/brutella/can"
)

type CANMessageLogger struct {
	logFile    io.Writer
	gobEncoder *gob.Encoder
}

type CANMessageIterator func() CANMessageGob

type CANMessageGob struct {
	Time  time.Time
	Frame can.Frame
}

func NewCanMessageDecoder(reader io.Reader) *CANMessageDecoder {
	return &CANMessageDecoder{
		gobDecoder: gob.NewDecoder(reader),
	}
}

func (d *CANMessageDecoder) Iterate() (*CANMessageGob, error) {
	var canMessage CANMessageGob
	err := d.gobDecoder.Decode(&canMessage)

	return &canMessage, err
}

type CANMessageDecoder struct {
	gobDecoder *gob.Decoder
}

func NewCanMessageLogger(writer io.Writer) *CANMessageLogger {
	return &CANMessageLogger{
		logFile:    writer,
		gobEncoder: gob.NewEncoder(writer),
	}
}

func (logger *CANMessageLogger) Handle(frame can.Frame) {
	logger.gobEncoder.Encode(&CANMessageGob{
		Time:  time.Now(),
		Frame: frame,
	})
}
