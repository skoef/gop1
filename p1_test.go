package gop1

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadData(t *testing.T) {
	t.Parallel()

	testdata, err := os.ReadFile("testdata/parser/output0")
	require.NoError(t, err)

	// create p1 with fake io.reader
	p1 := P1{
		serialDevice: bytes.NewReader(testdata),
		Incoming:     make(chan *Telegram),
	}

	go p1.readData()

	telegrams := make([]*Telegram, 0)
	for telegram := range p1.Incoming {
		telegrams = append(telegrams, telegram)
	}

	assert.Len(t, telegrams, 1)
	assert.Len(t, telegrams[0].Objects, 35)
}
