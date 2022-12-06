package log

import (
	"os"
	"testing"
)

func Test_write(t *testing.T) {
	logger := lokiLogger{
		BufferSize: 5,
	}
	logger.Write([]byte("{\"some\": \"foo\"}"))
	logger.Write([]byte("test2"))
	logger.Write([]byte("test3"))
	logger.Write([]byte("test4"))
	logger.Write([]byte("test5"))
	logger.Write([]byte("test6"))
}

func Test_debug(t *testing.T) {
	_ = os.Setenv("GF_PLUGIN_LOGGER", "loki")
	_ = os.Setenv("GF_PLUGIN_LOGGER_BUFFER", "5")

	Debug("foo1")
	Debug("foo2")
	Debug("foo3")
	Debug("foo4")
	Debug("foo5")
	Debug("foo6")
}
