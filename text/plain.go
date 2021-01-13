// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text // import "gonum.org/v1/plot/text"

import (
	"math"
	"strings"

	"gonum.org/v1/plot/vg"
)

// Plain is a text/plain handler.
type Plain struct{}

var _ Handler = (*Plain)(nil)

// Lines splits a given block of text into separate lines.
func (hdlr Plain) Lines(txt string) []string {
	txt = strings.TrimRight(txt, "\n")
	return strings.Split(txt, "\n")
}

// Box returns the bounding box of the given non-multiline text where:
//  - width is the horizontal space from the origin.
//  - height is the vertical space above the baseline.
//  - depth is the vertical space below the baseline, a positive number.
func (hdlr Plain) Box(txt string, fnt vg.Font) (width, height, depth vg.Length) {
	ext := fnt.Extents()
	width = fnt.Width(txt)
	height = ext.Ascent
	depth = ext.Descent

	return width, height, depth
}

// Draw renders the given text with the provided style and position
// on the canvas.
func (hdlr Plain) Draw(c vg.Canvas, txt string, sty TextStyle, pt vg.Point) {
	txt = strings.TrimRight(txt, "\n")
	if len(txt) == 0 {
		return
	}

	c.SetColor(sty.Color)

	if sty.Rotation != 0 {
		c.Push()
		c.Rotate(sty.Rotation)
	}

	sin64, cos64 := math.Sincos(sty.Rotation)
	cos := vg.Length(cos64)
	sin := vg.Length(sin64)
	pt.X, pt.Y = pt.Y*sin+pt.X*cos, pt.Y*cos-pt.X*sin

	lines := hdlr.Lines(txt)
	ht := sty.Height(txt)
	pt.Y += ht*vg.Length(sty.YAlign) - sty.Font.Extents().Ascent
	for i, line := range lines {
		xoffs := vg.Length(sty.XAlign) * sty.Font.Width(line)
		n := vg.Length(len(lines) - i)
		c.FillString(sty.Font, pt.Add(vg.Point{X: xoffs, Y: n * sty.Font.Size}), line)
	}

	if sty.Rotation != 0 {
		c.Pop()
	}
}
