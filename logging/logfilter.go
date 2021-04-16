package logging

import (
	"io"
	"log"

	"github.com/hashicorp/logutils"
)

// SetLogFilter will override the minimum log level.
func SetLogFilter(minLevel string, writer io.Writer) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "NOTICE", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(minLevel),
		Writer:   writer,
	}
	log.SetOutput(filter)
}
