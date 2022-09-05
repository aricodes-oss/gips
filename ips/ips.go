package ips

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

type ManyRecords []*Record

type PatchFile struct {
	Records ManyRecords
	Reader  *bufio.Reader
}

func (p *PatchFile) Patch(input []byte) []byte {
	max := p.MaxAddress()
	if len(input) < max {
		difference := max - len(input)
		input = append(input, make([]byte, difference)...)
	}

	for _, record := range p.Records {
		input = record.Patch(input)
	}

	return input
}

func (p *PatchFile) MaxAddress() (result int) {
	for _, record := range p.Records {
		end := record.End()
		if end > result {
			result = end
		}
	}

	return
}

func (p *PatchFile) ValidateHeader() error {
	header := make([]byte, 5)
	io.ReadFull(p.Reader, header)

	if string(header) != "PATCH" {
		return errors.New("Patch file is not an IPS patch, exiting")
	}

	return nil
}

func (p *PatchFile) LoadRecords() error {
	var (
		peeked []byte
		err    error
	)

	// Pre-allocate a large-ish buffer so that we're not doing constant
	// append operations in a loop - roughly 64kb
	p.Records = make(ManyRecords, 1024*8)

	// The odds of us getting an IPS patch that is literally
	// PATCHEOF should be basically zero, but somebody somewhere
	// has definitely made a patch that bad so now I must support it
	peeked, err = p.Reader.Peek(3)
	if err != nil {
		return err
	}

	idx := 0

	for string(peeked) != "EOF" {
		next := &Record{}
		next.Read(p.Reader)

		if idx >= len(p.Records) {
			p.Records = append(p.Records, next)
		} else {
			p.Records[idx] = next
		}

		idx++

		peeked, err = p.Reader.Peek(3)
		if err != nil {
			return err
		}
	}

	// Re-allocate and trim the slice down so we're not holding a bunch of 0 values
	p.Records = p.Records[:idx]

	return nil
}

func LoadPatchFile(filename string) (result *PatchFile) {
	result = &PatchFile{}

	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	result.Reader = bufio.NewReader(bytes.NewReader(data))
	result.ValidateHeader()

	return
}
