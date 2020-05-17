// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"fmt"
	"math"

	"github.com/go-latex/latex/drawtex"
	"github.com/go-latex/latex/font/ttf"
	"github.com/go-latex/latex/mtex"
	"github.com/go-latex/latex/tex"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Latex parses, formats and renders LaTeX.
type Latex struct{}

// Box returns the bounding box of the given text where:
//  - width is the horizontal space from the origin.
//  - height is the vertical space above the baseline.
//  - depth is the vertical space below the baseline, a negative number.
func (hdlr Latex) Box(txt string, fnt vg.Font, dpi float64) (width, height, depth vg.Length) {
	cnv := drawtex.New()
	box, err := mtex.Parse(txt, fnt.Size.Points(), 72, ttf.New(cnv))
	if err != nil {
		panic(fmt.Errorf("could not parse math expression: %w", err))
	}

	var sh tex.Ship
	sh.Call(0, 0, box.(tex.Tree))

	width = vg.Length(box.Width())
	height = vg.Length(box.Height())
	depth = vg.Length(box.Depth())

	return width, height, -depth
}

// Draw renders the given text with the provided style, position and
// dots-per-inch on the canvas.
func (hdlr Latex) Draw(c *draw.Canvas, txt string, sty draw.TextStyle, pt vg.Point, dpi float64) {
	cnv := drawtex.New()
	box, err := mtex.Parse(txt, sty.Font.Size.Points(), 72, ttf.New(cnv))
	if err != nil {
		panic(fmt.Errorf("could not parse math expression: %w", err))
	}

	var sh tex.Ship
	sh.Call(0, 0, box.(tex.Tree))

	w := box.Width()
	h := box.Height()
	d := box.Depth()

	o := latex{
		cnv: c,
		sty: sty,
		pt:  pt,
	}
	err = o.Render(w/72, math.Ceil(h+math.Max(d, 0))/72, dpi, cnv)
	if err != nil {
		panic(fmt.Errorf("could not render math expression: %w", err))
	}
}

type latex struct {
	cnv *draw.Canvas
	sty draw.TextStyle
	pt  vg.Point
}

func (r *latex) Render(width, height, dpi float64, c *drawtex.Canvas) error {
	//var (
	//	w = width * dpi
	//	h = height * dpi
	//)

	r.cnv.SetColor(r.sty.Color)
	r.cnv = &draw.Canvas{
		Canvas: r.cnv.Canvas,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: r.pt.X, Y: r.pt.Y},
			Max: vg.Point{X: r.pt.X + vg.Length(width), Y: r.pt.Y + vg.Length(height)},
		},
	}

	for _, op := range c.Ops() {
		switch op := op.(type) {
		case drawtex.GlyphOp:
			r.drawGlyph(dpi, op)
		case drawtex.RectOp:
			r.drawRect(dpi, op)
		default:
			panic(fmt.Errorf("unknown drawtex op %T", op))
		}
	}

	return nil
}

func (r *latex) drawGlyph(dpi float64, op drawtex.GlyphOp) {
	fnt := r.sty.Font
	fnt.Size = vg.Length(op.Glyph.Size)

	x := vg.Length(op.X * dpi / 72)
	y := vg.Length(op.Y * dpi / 72)
	r.cnv.FillString(fnt, vg.Point{X: x, Y: y}, op.Glyph.Symbol)
}

func (r *latex) drawRect(dpi float64, op drawtex.RectOp) {
	x1 := vg.Length(op.X1 * dpi / 72)
	x2 := vg.Length(op.X2 * dpi / 72)
	y1 := vg.Length(op.Y1 * dpi / 72)
	y2 := vg.Length(op.Y2 * dpi / 72)
	pts := []vg.Point{
		{X: x1, Y: y1},
		{X: x2, Y: y1},
		{X: x2, Y: y2},
		{X: x2, Y: y2},
		{X: x1, Y: y1},
	}

	r.cnv.FillPolygon(r.sty.Color, pts)
}

var (
	_ draw.TextHandler = (*Latex)(nil)
	_ mtex.Renderer    = (*latex)(nil)
)
