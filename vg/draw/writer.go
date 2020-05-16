// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw

import (
	"gonum.org/v1/plot/vg"
)

const defaultBufSize = 1024

// Writer implements buffering for a vg.Writer.
type Writer struct {
	buf *[]vg.Op
	n   int
	w   vg.Writer
}

func NewWriter(w vg.Writer) *Writer {
	return newWriterSize(w, defaultBufSize)
}

func newWriterSize(w vg.Writer, size int) *Writer {
	b, ok := w.(*Writer)
	if ok && len(*b.buf) >= size {
		return b
	}
	if size <= 0 {
		size = defaultBufSize
	}
	buf := make([]vg.Op, size)
	return &Writer{
		buf: &buf,
		w:   w,
	}
}

func (w *Writer) Size() int      { return len(*w.buf) }
func (w *Writer) Available() int { return len(*w.buf) - w.n }

func (w *Writer) Flush() {
	if w.n == 0 {
		return
	}
	w.w.Write((*w.buf)[0:w.n])
	w.n = 0
	return
}

func (w *Writer) Write(ops []vg.Op) {
	for len(ops) > w.Available() {
		switch {
		case w.n == 0:
			// large write, empty buffer.
			// write directly from ops to avoid copy
			w.w.Write(ops)
		default:
			n := copy((*w.buf)[w.n:], ops)
			w.n += n
			w.Flush()
			ops = ops[n:]
		}
	}
	n := copy((*w.buf)[w.n:], ops)
	w.n += n
}

func (w *Writer) WriteOp(op vg.Op) {
	if w.Available() <= 0 {
		w.Flush()
	}
	(*w.buf)[w.n] = op
	w.n++
}

var (
	_ vg.Writer = (*Writer)(nil)
)
