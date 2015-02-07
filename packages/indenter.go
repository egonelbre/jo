package packages

import (
	"bytes"
	"io"
)

// Writer indents each line of its input.
type indenter struct {
	w      io.Writer
	add    bool
	indent []byte
}

func newIndenter(w io.Writer, indent []byte) io.Writer {
	return &indenter{w, true, indent}
}

func (w *indenter) Write(buf []byte) (n int, err error) {
	for len(buf) > 0 {
		if w.add {
			_, err = w.w.Write(w.indent)
			if err != nil {
				return n, err
			}
			w.add = false
		}

		p := bytes.IndexRune(buf, '\n')
		if p < 0 {
			x, err := w.w.Write(buf)
			n += x
			return n, err
		}

		x, err := w.w.Write(buf[:p+1])
		n += x
		if err != nil {
			return n, err
		}
		buf = buf[p+1:]
		w.add = true
	}

	return n, nil
}
