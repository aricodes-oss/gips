package gips

import (
	"bytes"
	_ "embed"
	"testing"
)

//go:embed record.ips
var record []byte

func testBuf(b []byte, t *testing.T) {
	patch, err := NewPatchFrom(bytes.NewReader(b))
	if err != nil {
		t.Fail()
	}

	t.Logf("Read %s (%v)", patch.Contents(), patch.Contents())
}

func TestStandardRecord(t *testing.T) {
	testBuf(record, t)
}

//go:embed rle.ips
var rleRecord []byte

func TestRLERecord(t *testing.T) {
	testBuf(rleRecord, t)
}

func BenchmarkPatchFrom(b *testing.B) {
	buf := bytes.NewReader(record)

	for i := 0; i < b.N; i++ {
		NewPatchFrom(buf)
		buf.Reset(record)
	}
}
