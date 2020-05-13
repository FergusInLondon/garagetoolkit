package gpx

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/twpayne/go-gpx"
)

// TrackPoint ...
type TrackPoint struct {
	Latitude  float64
	Longitude float64
	Elevation float64
	Speed     float64
}

// Logger ...
type Logger struct {
	gpxPoint     chan TrackPoint
	ctx          context.Context
	cancel       context.CancelFunc
	gpxFile      io.WriteCloser
	gpxWaypoints []*gpx.WptType
	journeyBegin time.Time
}

// NewGPXLogger ...
func NewGPXLogger(parentCtx context.Context) *Logger {
	ctx, cancel := context.WithCancel(parentCtx)
	startTime := time.Now()
	return &Logger{
		gpxPoint:     make(chan TrackPoint, 10),
		ctx:          ctx,
		cancel:       cancel,
		gpxWaypoints: make([]*gpx.WptType, 1),
		journeyBegin: startTime,
	}
}

// Start ...
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

// Log ...
func (l *Logger) Log(point TrackPoint) {
	l.gpxPoint <- point
}

// Stop ...
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
