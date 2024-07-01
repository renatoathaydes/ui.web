package src

import (
	"bytes"
	"sync"
)

func NewWriter(write_line func(string)) *LogWriter {
	mutex := new(sync.Mutex)
	buf := new(bytes.Buffer)
	return &LogWriter{
		mutex:      mutex,
		write_line: write_line,
		buf:        buf,
	}
}

type LogWriter struct {
	mutex      *sync.Mutex
	write_line func(string)
	buf        *bytes.Buffer
}

func (w *LogWriter) Write(bytes []byte) (int, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	prev_line_pos := 0
	for i, b := range bytes {
		if b == '\n' {
			w.flush(bytes[prev_line_pos:i])
			prev_line_pos = i + 1
		}
	}
	w.buf.Write(bytes[prev_line_pos:])
	return len(bytes), nil
}

func (w *LogWriter) Close() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.write_line(w.buf.String())
	w.buf = nil
}

func (w *LogWriter) flush(bytes []byte) {
	var to_write []byte
	prev_bytes := w.buf.Bytes()
	if len(prev_bytes) == 0 {
		to_write = bytes
	} else {
		to_write = append(to_write, prev_bytes...)
		to_write = append(to_write, bytes...)
		w.buf.Reset()
	}
	w.write_line(string(to_write))
}
