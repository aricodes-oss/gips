package ips

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

type Record struct {
	Data   []byte
	Size   uint16
	Offset uint32
	RLE    bool
}

// Shorthand, since these will be appearing so frequently
var buint16 = binary.BigEndian.Uint16
var buint32 = binary.BigEndian.Uint32
var rfull = io.ReadFull

func (r *Record) Read(input *bufio.Reader) *Record {
	offset := make([]byte, 3)
	size := make([]byte, 2)
	rfull(input, offset)
	rfull(input, size)

	r.Offset = buint32(append([]byte{0x00}, offset...))
	r.Size = buint16(size)

	// 0-size records are RLE encoded, the next 2 bytes tell us
	// how many times to repeat the next single byte value
	if r.Size == 0 {
		r.RLE = true

		size := make([]byte, 2)
		rfull(input, size)
		r.Size = buint16(size)

		r.Data = make([]byte, 1)
	} else {
		r.Data = make([]byte, r.Size)
	}

	rfull(input, r.Data)

	return r
}

func (r *Record) End() int {
	return int(uint32(r.Size) + r.Offset)
}

func (r *Record) Patch(input []byte) []byte {
	if r.RLE {
		for address := r.Offset; address < r.Offset+uint32(r.Size); address++ {
			input[address] = r.Data[0]
		}

		return input
	}

	for idx, val := range r.Data {
		input[int(r.Offset)+idx] = val
	}

	return input
}
