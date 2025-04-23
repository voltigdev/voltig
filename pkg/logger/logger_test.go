package logger

import (
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"

	"voltig/pkg/logger/testdata"
)

func TestLogger(t *testing.T) {
	buf := testdata.NewBufferWriter()
	cfg := Config{
		Level:      log.InfoLevel,
		TimeFormat: "15:04:05",
		Output:     buf,
		Prefix:     "test",
		ShowCaller: false,
	}
	Configure(cfg)

	Info("info message", "key", "value")
	Warn("warn message")
	Error("error message", "err", "fail")

	output := buf.String()

	assert.Contains(t, output, "info message")
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
	assert.Contains(t, output, "key=value")
	assert.Contains(t, output, "err=fail")
}
