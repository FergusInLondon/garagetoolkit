package canlog

import (
	"encoding/gob"
	"io"

	"github.com/brutella/can"
)

// Parser is used for parsing persisted binary CAN logs. It contains no exported
// properties.
type Parser struct {
	gobDecoder *gob.Decoder
}

// DecodedMessage contains the time at which a message *was processed for logging*,
// the raw bytes of the frame, and the parsed frame.
type DecodedMessage struct {
	*Message
	Frame can.Frame
}

// NewParser returns a configured Parser, complete with underlying decoder.
func NewParser(reader io.Reader) *Parser {
	return &Parser{
		gobDecoder: gob.NewDecoder(reader),
	}
}

// Iterate wraps around the underlying decoder and provides a DecodedMessage
// per log item.
func (p *Parser) Iterate() (*DecodedMessage, error) {
	decoded := &DecodedMessage{}
	if err := p.gobDecoder.Decode(&decoded.Message); err != nil {
		return nil, err
	}

	if err := can.Unmarshal(decoded.Raw, &decoded.Frame); err != nil {
		return nil, err
	}

	return decoded, nil
}
