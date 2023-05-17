package gips

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Patch struct {
	Data   []byte
	Offset uint32
	Size   uint16

	IsRLE   bool
	RLESize uint16
}

func (p *Patch) Contents() []byte {
	if p.IsRLE {
		return bytes.Repeat(p.Data, int(p.RLESize))
	}
	return p.Data
}

func NewPatchFrom(buf io.Reader) (*Patch, error) {
	p := &Patch{}

	prelude := make([]byte, 5)
	_, err := io.ReadFull(buf, prelude)
	if err != nil {
		return nil, err
	}

	// IPS patch files use a uint24 to store offsets
	p.Offset = binary.BigEndian.Uint32(append([]byte{0x00}, prelude[:3]...))
	p.Size = binary.BigEndian.Uint16(prelude[3:])
	p.Data = make([]byte, p.Size)
	p.IsRLE = p.Size == 0
	if p.IsRLE {
		rleSizeBuf := make([]byte, 2)
		_, err := io.ReadFull(buf, rleSizeBuf)
		if err != nil {
			return nil, err
		}

		p.RLESize = binary.BigEndian.Uint16(rleSizeBuf)
		p.Data = make([]byte, 1)
	}

	_, err = io.ReadFull(buf, p.Data)
	if err != nil {
		return nil, err
	}

	return p, nil
}
