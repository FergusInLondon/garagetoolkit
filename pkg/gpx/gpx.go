package gpx

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/twpayne/go-gpx"
)

var (
	// TrackPointBufferLength is the size of the buffer for the inbound TrackPoint channel
	TrackPointBufferLength = 16
)

// TrackPoint contains the four specific properties of a GPX Waypoint that we're
// interested in; namely Latitude, Longitude, Elevation, and Speed.
type TrackPoint struct {
	Latitude  float64
	Longitude float64
	Elevation float64
	Speed     float64
}

// Logger is responsible for logging GPS data to a GPX file; it exports no properties.
type Logger struct {
	gpxPoint     chan TrackPoint
	ctx          context.Context
	cancel       context.CancelFunc
	gpxFile      io.WriteCloser
	gpxWaypoints []*gpx.WptType
	journeyBegin time.Time
}

// NewGPXLogger provides a new Logger that writes to an underlying GPX io.WriteCloser
// - i.e generally an .XML file.
func NewGPXLogger(parentCtx context.Context, gpxFile io.WriteCloser) *Logger {
	ctx, cancel := context.WithCancel(parentCtx)
	startTime := time.Now()
	return &Logger{
		gpxPoint:     make(chan TrackPoint, TrackPointBufferLength),
		ctx:          ctx,
		cancel:       cancel,
		gpxWaypoints: make([]*gpx.WptType, 1),
		gpxFile:      gpxFile,
		journeyBegin: startTime,
	}
}

// Start kickstarts the background worker which parses inbound TrackPoint structs
// and converts them to GPX XML ready Waypoint entities.
func (l *Logger) Start() {
	go func() {
		for {
			select {
			case <-l.ctx.Done():
				return
			case point := <-l.gpxPoint:
				l.gpxWaypoints = append(l.gpxWaypoints, &gpx.WptType{
					Lat:   point.Latitude,
					Lon:   point.Longitude,
					Speed: point.Speed,
					Ele:   point.Elevation,
					Time:  time.Now(),
				})
			}
		}
	}()
}

// Log takes a TrackPoint and passes it to the internal background runner which parses
// it to a valid GPX Waypoint Type.
func (l *Logger) Log(point TrackPoint) {
	l.gpxPoint <- point
}

// Stop is responsible for concluding the process of logging GPS data; it stops the
// background worker from processing inbound GPS data, before generating the GPX XML
// document and writing it and closing the file.
func (l *Logger) Stop() error {
	l.cancel()

	gpxDocument := &gpx.GPX{
		Version: "1.0",
		Creator: "CANLog 0.1 - https://fergus.london/",
		Trk: []*gpx.TrkType{{
			Name: fmt.Sprintf("%s - %s", l.journeyBegin.String(), time.Now().String()),
			TrkSeg: []*gpx.TrkSegType{{
				TrkPt: l.gpxWaypoints,
			}},
		}},
	}

	if _, err := l.gpxFile.Write([]byte(xml.Header)); err != nil {
		return err
	}

	if err := gpxDocument.Write(l.gpxFile); err != nil {
		return err
	}

	return l.gpxFile.Close()
}
