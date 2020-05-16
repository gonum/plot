// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vg

type Writer interface {
	Write(ops []Op)
}

type Reader interface {
	Read(ops []Op)
}

func WriterFrom(c Canvas) Writer {
	switch c := c.(type) {
	case Writer:
		return c
	}
	return &writer{c: c}
}

type writer struct {
	c Canvas
}

func (w writer) Write(ops []Op) {
	for _, op := range ops {
		w.WriteOp(op)
	}
}

func (w writer) WriteOp(op Op) {
	switch op := op.(type) {
	case LineWidth:
		w.c.SetLineWidth(op.Width)
	case LineDash:
		w.c.SetLineDash(op.Pattern, op.Offset)
	case Color:
		w.c.SetColor(op.Color)
	case Rotate:
		w.c.Rotate(op.Radians)
	case Translate:
		w.c.Translate(op.Point)
	case Scale:
		w.c.Scale(op.X, op.Y)
	case Push:
		w.c.Push()
	case Pop:
		w.c.Pop()
	case Stroke:
		w.c.Stroke(op.Path)
	case Fill:
		w.c.Fill(op.Path)
	case FillString:
		w.c.FillString(op.Font, op.Point, op.Text)
	case DrawImage:
		w.c.DrawImage(op.Rect, op.Image)
	}
}
