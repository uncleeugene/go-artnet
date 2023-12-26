package packet

import (
	"github.com/uncleeugene/go-artnet/packet/code"
)

var _ ArtNetPacket = &ArtDiagDataPacket{}

// ArtDiagDataPacket contains an ArtDiagData Packet.
//
// ArtDiagData is a general purpose packet that allows a node or controller to send
// diagnostics data for display. The ArtPoll packet sent by controllers defines the
// destination to which these messages should be sent.
//
// Packet Strategy:
//
//	Controller -  Receive:            Application Specific
//	              Unicast Transmit:   As defined by ArtPoll
//	              Broadcast Transmit: As defined by ArtPoll
//	Node -        Receive:            No Action
//	              Unicast Transmit:   As defined by ArtPoll
//	              Broadcast Transmit: As defined by ArtPoll
//	MediaServer - Receive:            No Action
//	              Unicast Transmit:   As defined by ArtPoll
//	              Broadcast Transmit: As defined by ArtPoll
type ArtDiagDataPacket struct {
	// Inherit the Header header
	Header

	// Filler1
	_ byte

	// Priority contains the lowest priority of diagnostics message that should be sent
	Priority code.PriorityCode

	// Filler2
	_ [2]byte

	// Length indicates the length of the data
	Length uint16

	// Data is an ASCII string, null terminated. Max length is 512 bytes including the null terminator
	Data [512]byte
}

// NewArtDiagDataPacket returns an ArtNetPacket with the correct OpCode
func NewArtDiagDataPacket() *ArtDiagDataPacket {
	return &ArtDiagDataPacket{}
}

// MarshalBinary marshals an ArtDiagDataPacket into a byte slice.
func (p *ArtDiagDataPacket) MarshalBinary() ([]byte, error) {
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtDiagDataPacket.
func (p *ArtDiagDataPacket) UnmarshalBinary(b []byte) error {
	return unmarshalPacket(p, b)
}

// validate is used to validate the Packet.
func (p *ArtDiagDataPacket) validate() error {
	if err := p.Header.validate(); err != nil {
		return err
	}
	if p.OpCode != code.OpDiagData {
		return errInvalidOpCode
	}
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ArtDiagDataPacket) finish() {
	p.Header.finish()
}
