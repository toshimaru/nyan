package helpers

import (
	"bytes"
	"fmt"
	"io"
)

type NumberWriter struct {
	w           io.Writer
	currentLine uint64
	buf         []byte
}

func (w *NumberWriter) Write(p []byte) (n int, err error) {
	// Early return.
	// Can't calculate the line numbers until the line breaks are made, so store them all in a buffer.
	if !bytes.Contains(p, []byte{'\n'}) {
		w.buf = append(w.buf, p...)
		return len(p), nil
	}

	var (
		original = p
		tokenLen uint
	)
	for i, c := range original {
		tokenLen++
		if c != '\n' {
			continue
		}

		token := p[:tokenLen]
		p = original[i+1:]
		tokenLen = 0

		format := "%6d\t%s%s"
		if w.currentLine > 999999 {
			format = "%d\t%s%s"
		}

		_, er := fmt.Fprintf(w.w, format, w.currentLine, string(w.buf), string(token))
		if er != nil {
			return i + 1, er
		}
		w.buf = w.buf[:0]
		w.currentLine++
	}

	if len(p) > 0 {
		w.buf = append(w.buf, p...)
	}
	return len(original), nil
}

func (w *NumberWriter) Flush() error {
	terminalReset := []byte("\u001B[0m")
	if bytes.Compare(w.buf, terminalReset) == 0 {
		// In almost all cases, a control code is passed last to reset the terminal's color code.
		// This is not a printable character and should not be counted as a line, so it is output as is without a line number.
		_, err := fmt.Fprintf(w.w, "%s", string(w.buf))
		return err
	}

	format := "%6d\t%s"
	if w.currentLine > 999999 {
		format = "%d\t%s"
	}
	_, err := fmt.Fprintf(w.w, format, w.currentLine, string(w.buf))
	w.buf = w.buf[:0]
	return err
}
