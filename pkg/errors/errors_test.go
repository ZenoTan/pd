package errors

import (
	"bytes"
	"fmt"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"strings"
	"testing"
)

// testingWriter is a WriteSyncer that writes to the the messages.
type testingWriter struct {
	messages []string
}

func newTestingWriter() *testingWriter {
	return &testingWriter{}
}

func (w *testingWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	p = bytes.TrimRight(p, "\n")
	m := fmt.Sprintf("%s", p)
	w.messages = append(w.messages, m)
	return n, nil
}

func (w *testingWriter) Sync() error {
	return nil
}

type verifyLogger struct {
	*zap.Logger
	w *testingWriter
}

func (logger *verifyLogger) Contain(t *testing.T, s string) {
	if logger.w.messages == nil {
		t.Error()
	}
	msg := logger.w.messages[len(logger.w.messages) - 1]
	IsContain(t, msg, s)
}

func newZapTestLogger(cfg *log.Config, opts ...zap.Option) verifyLogger {
	// TestingWriter is used to write to memory.
	// Used in the verify logger.
	writer := newTestingWriter()
	lg, _, _ := log.InitLoggerWithWriteSyncer(cfg, writer, opts...)

	return verifyLogger{
		Logger: lg,
		w:      writer,
	}
}

func IsContain(t *testing.T, s1 string, s2 string) {
	if !strings.Contains(s1, s2) {
		t.Error()
	}
}

func TestError(t *testing.T) {
	conf := &log.Config{Level: "debug", File: log.FileLogConfig{}, DisableTimestamp: true}
	lg := newZapTestLogger(conf)
	log.ReplaceGlobals(lg.Logger, nil)

	rfc := `[error="[format:ErrFormatParseHistoryIndex] parse history index error"]`
	log.Error(ErrFormatParseHistoryIndex.MessageTemplate(), zap.Error(ErrFormatParseHistoryIndex.FastGenByArgs()))
	lg.Contain(t, rfc)
	rfc = `[error="[internal:InternalStoreNotFound] store id 1 not found"]`
	log.Error(InternalStoreNotFound.MessageTemplate(), zap.Error(InternalStoreNotFound.FastGenByArgs(1)))
	lg.Contain(t, rfc)
}