package logui

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/fatih/color"
)

type FmtTime = func(t time.Time) string

type LogUiHandler struct {
	level      slog.Level
	out        io.Writer
	group      string
	groupColor *color.Color
	mux        *sync.Mutex
	attrs      []slog.Attr
	attrs_buf  []byte
	fmt_time   FmtTime
}

// New creates a new LogUiHandler.
func New(level slog.Level, out io.Writer, group string, groupColor *color.Color) *LogUiHandler {
	return &LogUiHandler{level, out, group, groupColor, &sync.Mutex{}, []slog.Attr{}, []byte{}, fmtTime}
}

// New creates a new LogUiHandler with a time formatter.
func NewWithTimeFormatter(level slog.Level, out io.Writer, group string, fmt_time FmtTime) *LogUiHandler {
	return &LogUiHandler{level, out, group, nil, &sync.Mutex{}, []slog.Attr{}, []byte{}, fmt_time}
}

// Enabled implements slog.Handler.
func (l *LogUiHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= l.level
}

// WithAttrs implements slog.Handler.
func (l *LogUiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	res := NewWithTimeFormatter(l.level, l.out, l.group, l.fmt_time)
	res.groupColor = l.groupColor
	res.attrs = make([]slog.Attr, len(l.attrs)+len(attrs))
	res.attrs = append(res.attrs, l.attrs...)
	res.attrs = append(res.attrs, attrs...)
	res.attrs_buf = make([]byte, len(l.attrs_buf))
	copy(res.attrs_buf, l.attrs_buf)
	res.attrs_buf = appendAttrs(res.attrs_buf, attrs, len(l.attrs_buf) > 0)
	return res
}

// WithGroup implements slog.Handler.
func (l *LogUiHandler) WithGroup(name string) slog.Handler {
	res := New(l.level, l.out, name, l.groupColor)
	res.attrs = l.attrs
	return res
}

// Handle implements slog.Handler.
func (l *LogUiHandler) Handle(ctx context.Context, rec slog.Record) error {
	buf := make([]byte, 0, 512)
	buf = centerString(buf, l.group, l.groupColor, 6)
	b := bytes.NewBuffer(buf)
	levelColor(rec.Level).Fprintf(b, "[%s]", rec.Level.String())
	buf = b.Bytes()
	buf = fmt.Appendf(buf, " %s %s", l.fmt_time(rec.Time), rec.Message)
	if len(l.attrs_buf) > 0 || rec.NumAttrs() > 0 {
		buf = append(buf, " | "...)
	}
	buf = append(buf, l.attrs_buf...)
	prepend_comma := len(l.attrs) > 0 && rec.NumAttrs() > 0
	lastIndex := rec.NumAttrs() - 1
	index := 0
	rec.Attrs(func(a slog.Attr) bool {
		if prepend_comma {
			buf = append(buf, ", "...)
			prepend_comma = false
		}
		buf = appendAttr(buf, a, index == lastIndex)
		index++
		return true
	})
	if len(l.attrs_buf) > 0 || rec.NumAttrs() > 0 {
		buf = append(buf, " |\n"...)
	} else {
		buf = append(buf, '\n')
	}
	l.mux.Lock()
	defer l.mux.Unlock()
	_, err := l.out.Write(buf)
	return err
}

func appendAttrs(buf []byte, attrs []slog.Attr, prepend_comma bool) []byte {
	if len(attrs) == 0 {
		return buf
	}
	if prepend_comma {
		buf = append(buf, ", "...)
	} else {
		buf = append(buf, '{')
	}
	last := len(attrs) - 1
	for i, a := range attrs {
		buf = appendAttr(buf, a, i == last)
	}
	buf = append(buf, '}')
	return buf
}

func appendAttr(buf []byte, a slog.Attr, last bool) []byte {
	buf = append(buf, []byte(a.Key)...)
	buf = append(buf, '=')
	buf = appendValue(buf, &a.Value)
	if !last {
		buf = append(buf, ", "...)
	}
	return buf
}

func appendValue(buf []byte, v *slog.Value) []byte {
	switch v.Kind() {
	case slog.KindBool:
		buf = fmt.Appendf(buf, "%t", v.Bool())
	case slog.KindString:
		buf = fmt.Appendf(buf, "\"%s\"", v.String())
	case slog.KindTime:
		buf = fmt.Appendf(buf, "%s", v.Time().Format("2006-01-02T15:04:05.9999"))
	default:
		buf = fmt.Appendf(buf, "%s", v)
	}
	return buf
}

func fmtTime(t time.Time) string {
	return t.Format("15:04:05.9999")
}

func centerString(buf []byte, s string, c *color.Color, width int) []byte {
	w := bytes.NewBuffer(buf)
	if len(s) >= width {
		if c != nil {
			c.Fprintf(w, "%s", s)
		} else {
			w.WriteString(s)
		}
	} else {
		a := (width - len(s)) / 2
		b := a + len(s)
		if width%2 != 0 {
			a++
		}
		if c != nil {
			c.Fprintf(w, "%[1]*[3]s%[2]*[4]s", b, a, s, "")
		} else {
			fmt.Fprintf(w, "%[1]*[3]s%[2]*[4]s", b, a, s, "")
		}
	}
	return w.Bytes()
}
