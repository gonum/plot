// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"fmt"
	"math"
	"strings"

	"github.com/go-latex/latex/drawtex"
	"github.com/go-latex/latex/font/ttf"
	"github.com/go-latex/latex/mtex"
	"github.com/go-latex/latex/tex"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Latex parses, formats and renders LaTeX.
type Latex struct {
	// DPI is the dot-per-inch controlling the font resolution used by LaTeX.
	// If zero, the resolution defaults to 72.
	DPI float64
}

var _ draw.TextHandler = (*Latex)(nil)

// Lines splits a given block of text into separate lines.
func (hdlr Latex) Lines(txt string) []string {
	txt = strings.TrimRight(txt, "\n")
	return strings.Split(txt, "\n")
}

// Box returns the bounding box of the given non-multiline text where:
//  - width is the horizontal space from the origin.
//  - height is the vertical space above the baseline.
//  - depth is the vertical space below the baseline, a positive number.
func (hdlr Latex) Box(txt string, fnt vg.Font) (width, height, depth vg.Length) {
	cnv := drawtex.New()
	fnts := &ttf.Fonts{
		Rm:      fnt.Font(),
		Default: fnt.Font(),
		It:      fnt.Font(), // FIXME(sbinet): need a gonum/plot font set
	}
	box, err := mtex.Parse(txt, fnt.Size.Points(), latexDPI, ttf.NewFrom(cnv, fnts))
	if err != nil {
		panic(fmt.Errorf("could not parse math expression: %w", err))
	}

	var sh tex.Ship
	sh.Call(0, 0, box.(tex.Tree))

	width = vg.Length(box.Width())
	height = vg.Length(box.Height())
	depth = vg.Length(box.Depth())

	// Add a bit of space, with a linegap as mtex.Box is returning
	// a very tight bounding box.
	// See gonum/plot#661.
	if depth != 0 {
		var (
			e       = fnt.Extents()
			linegap = e.Height - (e.Ascent + e.Descent)
		)
		depth += linegap
	}

	dpi := vg.Length(hdlr.dpi() / latexDPI)
	return width * dpi, height * dpi, depth * dpi
}

// Draw renders the given text with the provided style and position
// on the canvas.
func (hdlr Latex) Draw(c *draw.Canvas, txt string, sty draw.TextStyle, pt vg.Point) {
	cnv := drawtex.New()
	fnts := &ttf.Fonts{
		Rm:      sty.Font.Font(),
		Default: sty.Font.Font(),
		It:      sty.Font.Font(), // FIXME(sbinet): need a gonum/plot font set
	}
	box, err := mtex.Parse(txt, sty.Font.Size.Points(), latexDPI, ttf.NewFrom(cnv, fnts))
	if err != nil {
		panic(fmt.Errorf("could not parse math expression: %w", err))
	}

	var sh tex.Ship
	sh.Call(0, 0, box.(tex.Tree))

	w := box.Width()
	h := box.Height()
	d := box.Depth()

	dpi := hdlr.dpi() / latexDPI
	o := latex{
		cnv: c,
		sty: sty,
		pt:  pt,
		w:   vg.Length(w * dpi),
		h:   vg.Length((h + d) * dpi),
		cos: 1,
		sin: 0,
	}
	e := sty.Font.Extents()
	o.xoff = vg.Length(sty.XAlign) * o.w
	o.yoff = o.h + o.h*vg.Length(sty.YAlign) - (e.Height - e.Ascent)

	if sty.Rotation != 0 {
		sin64, cos64 := math.Sincos(sty.Rotation)
		o.cos = vg.Length(cos64)
		o.sin = vg.Length(sin64)

		o.cnv.Push()
		defer o.cnv.Pop()
		o.cnv.Rotate(sty.Rotation)
	}

	err = o.Render(w/latexDPI, (h+d)/latexDPI, dpi, cnv)
	if err != nil {
		panic(fmt.Errorf("could not render math expression: %w", err))
	}
}

// latexDPI is the default LaTeX resolution used for computing the LaTeX
// layout of equations and regular text.
// Dimensions are then rescaled to the desired resolution.
const latexDPI = 72.0

func (hdlr Latex) dpi() float64 {
	if hdlr.DPI == 0 {
		return latexDPI
	}
	return hdlr.DPI
}

type latex struct {
	cnv *draw.Canvas
	sty draw.TextStyle
	pt  vg.Point

	w vg.Length
	h vg.Length

	cos vg.Length
	sin vg.Length

	xoff vg.Length
	yoff vg.Length
}

var _ mtex.Renderer = (*latex)(nil)

func (r *latex) Render(width, height, dpi float64, c *drawtex.Canvas) error {
	r.cnv.SetColor(r.sty.Color)

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

	pt := r.pt
	if r.sty.Rotation != 0 {
		pt.X, pt.Y = r.rotate(pt.X, pt.Y)
	}

	pt = pt.Add(vg.Point{
		X: r.xoff + vg.Length(op.X*dpi),
		Y: r.yoff - vg.Length(op.Y*dpi),
	})

	r.cnv.FillString(fnt, pt, op.Glyph.Symbol)
}

func (r *latex) drawRect(dpi float64, op drawtex.RectOp) {
	x1 := r.xoff + vg.Length(op.X1*dpi)
	x2 := r.xoff + vg.Length(op.X2*dpi)
	y1 := r.yoff - vg.Length(op.Y1*dpi)
	y2 := r.yoff - vg.Length(op.Y2*dpi)

	pt := r.pt
	if r.sty.Rotation != 0 {
		pt.X, pt.Y = r.rotate(pt.X, pt.Y)
	}

	pts := []vg.Point{
		vg.Point{X: x1, Y: y1}.Add(pt),
		vg.Point{X: x2, Y: y1}.Add(pt),
		vg.Point{X: x2, Y: y2}.Add(pt),
		vg.Point{X: x1, Y: y2}.Add(pt),
		vg.Point{X: x1, Y: y1}.Add(pt),
	}

	r.cnv.FillPolygon(r.sty.Color, pts)
}

func (r *latex) rotate(x, y vg.Length) (vg.Length, vg.Length) {
	u := x*r.cos + y*r.sin
	v := y*r.cos - x*r.sin
	return u, v
}
