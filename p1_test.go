package gop1

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadData(t *testing.T) {
	testdata, err := os.ReadFile("testdata/parser/output0")
	require.NoError(t, err)

	// create p1 with fake io.reader
	p1 := P1{
		serialDevice: bytes.NewReader(testdata),
		Incoming:     make(chan *Telegram),
	}

	go p1.readData()
	telegrams := 0
	for range p1.Incoming {
		telegrams++
	}
	assert.Equal(t, 500, telegrams)
}
