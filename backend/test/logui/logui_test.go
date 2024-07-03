package logui

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"ui.web/server/src/logui"
)

func TestLogSimple(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	logger := slog.New(logui.NewWithTimeFormatter(slog.LevelError, buf, "fe", fmtTime))
	logger.Info("hello world")
	require.Equal(t, "  fe   [INFO] T - hello world\n", buf.String())
}

func TestLogWithAttr(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	logger := slog.New(logui.NewWithTimeFormatter(slog.LevelError, buf, "be", fmtTime).WithAttrs([]slog.Attr {
		{Key: "proc", Value: slog.StringValue("node")},
	}))
	logger.Warn("foo")
	require.Equal(t, "  be   [WARN] T {proc=\"node\"} - foo\n", buf.String())
}


func fmtTime(t time.Time) string {
	return "T"
}
