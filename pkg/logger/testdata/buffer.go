package testdata

import (
	"bytes"
)

// BufferWriter is a simple io.Writer that captures written bytes for testing.
type BufferWriter struct {
	buf bytes.Buffer
}

func (b *BufferWriter) Write(p []byte) (n int, err error) {
	return b.buf.Write(p)
}

func (b *BufferWriter) String() string {
	return b.buf.String()
}

func NewBufferWriter() *BufferWriter {
	return &BufferWriter{}
}
