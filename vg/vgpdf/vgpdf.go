// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vgpdf implements the vg.Canvas interface
// using gopdf (bitbucket.org/zombiezen/gopdf/pdf).
package vgpdf // import "gonum.org/v1/plot/vg/vgpdf"

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"

	pdf "github.com/jung-kurt/gofpdf"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/fonts"
)

// DPI is the nominal resolution of drawing in PDF.
const DPI = 72

// Canvas implements the vg.Canvas interface,
// drawing to a PDF.
type Canvas struct {
	pdf         *pdf.Fpdf
	w, h        vg.Length
	lineVisible bool

	dpi   int
	nimgs int // images counter
	ctx   []context
	fonts map[string]struct{}
}

type context struct {
	fillc color.Color // fill color
	drawc color.Color // draw color
	linew vg.Length   // line width
}

// New creates a new PDF Canvas.
func New(w, h vg.Length) *Canvas {
	cfg := pdf.InitType{
		UnitStr: "pt",
		Size:    pdf.SizeType{Wd: w.Points(), Ht: h.Points()},
	}
	c := &Canvas{
		pdf:         pdf.NewCustom(&cfg),
		w:           w,
		h:           h,
		lineVisible: true,
		dpi:         DPI,
		ctx:         make([]context, 1),
		fonts:       make(map[string]struct{}),
	}
	vg.Initialize(c)
	c.pdf.SetMargins(0, 0, 0)
	c.pdf.AddPage()
	c.Push()
	c.Translate(vg.Point{0, h})
	c.Scale(1, -1)

	return c
}

func (c *Canvas) DPI() float64 {
	return float64(c.dpi)
}

func (c *Canvas) cur() *context {
	return &c.ctx[len(c.ctx)-1]
}

func (c *Canvas) Size() (w, h vg.Length) {
	return c.w, c.h
}

func (c *Canvas) SetLineWidth(w vg.Length) {
	c.cur().linew = w
	lw := c.unit(w)
	c.pdf.SetLineWidth(lw)
	c.lineVisible = w > 0
}

func (c *Canvas) SetLineDash(dashes []vg.Length, offs vg.Length) {
	ds := make([]float64, len(dashes))
	for i, d := range dashes {
		ds[i] = c.unit(d)
	}
	c.pdf.SetDashPattern(ds, c.unit(offs))
}

func (c *Canvas) SetColor(clr color.Color) {
	if clr == nil {
		clr = color.Black
	}
	c.cur().drawc = clr
	c.cur().fillc = clr
	c.pdf.SetFillColor(rgb(clr))
	c.pdf.SetDrawColor(rgb(clr))
	c.pdf.SetTextColor(rgb(clr))
}

func (c *Canvas) Rotate(r float64) {
	c.pdf.TransformRotate(-r*180/math.Pi, 0, 0)
}

func (c *Canvas) Translate(pt vg.Point) {
	xp, yp := c.pdfPoint(pt)
	c.pdf.TransformTranslate(xp, yp)
}

func (c *Canvas) Scale(x float64, y float64) {
	c.pdf.TransformScale(x*100, y*100, 0, 0)
}

func (c *Canvas) Push() {
	c.ctx = append(c.ctx, *c.cur())
	c.pdf.TransformBegin()
}

func (c *Canvas) Pop() {
	c.pdf.TransformEnd()
	c.ctx = c.ctx[:len(c.ctx)-1]
}

func (c *Canvas) Stroke(p vg.Path) {
	if c.lineVisible {
		c.pdfPath(p, "D")
	}
}

func (c *Canvas) Fill(p vg.Path) {
	c.pdfPath(p, "F")
}

func (c *Canvas) FillString(fnt vg.Font, pt vg.Point, str string) {
	c.font(fnt, pt)
	c.pdf.SetFont(fnt.Name(), "", c.unit(fnt.Size))

	c.Push()
	defer c.Pop()
	c.Translate(pt)
	// go-fpdf uses the top left corner as origin.
	c.Scale(1, -1)
	left, top, right, bottom := c.sbounds(fnt, str)
	w := right - left
	h := bottom - top
	margin := c.pdf.GetCellMargin()

	c.pdf.MoveTo(-left-margin, top)
	c.pdf.CellFormat(w, h, str, "", 0, "BL", false, 0, "")
}

func (c *Canvas) sbounds(fnt vg.Font, txt string) (left, top, right, bottom float64) {
	_, h := c.pdf.GetFontSize()
	d := c.pdf.GetFontDesc("", "")
	if d.Ascent == 0 {
		// not defined (standard font?), use average of 81%
		top = 0.81 * h
	} else {
		top = -float64(d.Ascent) * h / float64(d.Ascent-d.Descent)
	}
	return 0, top, c.pdf.GetStringWidth(txt), top + h
}

// whitespace determines if a rune is whitespace
func whitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

// DrawImage implements the vg.Canvas.DrawImage method.
func (c *Canvas) DrawImage(rect vg.Rectangle, img image.Image) {
	opts := pdf.ImageOptions{ImageType: "png", ReadDpi: true}
	name := c.imageName()

	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		log.Panicf("error encoding image to PNG: %v", err)
	}
	c.pdf.RegisterImageOptionsReader(name, opts, buf)

	xp, yp := c.pdfPoint(rect.Min)
	wp, hp := c.pdfPoint(rect.Size())

	c.pdf.ImageOptions(name, xp, yp, wp, hp, false, opts, 0, "")
}

// font registers a font and a size with the PDF canvas.
func (c *Canvas) font(fnt vg.Font, pt vg.Point) {
	if _, ok := c.fonts[fnt.Name()]; ok {
		return
	}
	if n, ok := vg.FontMap[fnt.Name()]; ok {
		json, err := fonts.Asset(n + ".json")
		if err != nil {
			log.Panicf("vgpdf: could not load JSON description for TTF font %q: %v", n+".json", err)
		}

		zfile, err := fonts.Asset(n + ".z")
		if err != nil {
			log.Panicf("vgpdf: could not load zip file for TTF font: %q: %v", n+".z", err)
		}

		c.fonts[fnt.Name()] = struct{}{}
		c.pdf.AddFontFromBytes(fnt.Name(), "", json, zfile)
		return
	}
	log.Panicf("vgpdf: not implemented")
}

// pdfPath processes a vg.Path and applies it to the canvas.
func (c *Canvas) pdfPath(path vg.Path, style string) {
	var (
		xp float64
		yp float64
	)
	for _, comp := range path {
		switch comp.Type {
		case vg.MoveComp:
			xp, yp = c.pdfPoint(comp.Pos)
			c.pdf.MoveTo(xp, yp)
		case vg.LineComp:
			c.pdf.LineTo(c.pdfPoint(comp.Pos))
		case vg.ArcComp:
			// FIXME(sbinet): use c.pdf.ArcTo
			c.arc(comp, style)
		case vg.CloseComp:
			c.pdf.LineTo(xp, yp)
			c.pdf.ClosePath()
		default:
			panic(fmt.Sprintf("Unknown path component type: %d\n", comp.Type))
		}
	}
	c.pdf.DrawPath(style)
	return
}

// Approximate a circular arc using multiple
// cubic Bézier curves, one for each π/2 segment.
//
// This is from:
// 	http://hansmuller-flex.blogspot.com/2011/04/approximating-circular-arc-with-cubic.html
func (c *Canvas) arc(comp vg.PathComp, style string) {
	x0 := comp.Pos.X + comp.Radius*vg.Length(math.Cos(comp.Start))
	y0 := comp.Pos.Y + comp.Radius*vg.Length(math.Sin(comp.Start))
	c.pdf.LineTo(c.pdfPointXY(x0, y0))

	a1 := comp.Start
	end := a1 + comp.Angle
	sign := 1.0
	if end < a1 {
		sign = -1.0
	}
	left := math.Abs(comp.Angle)

	// Square root of the machine epsilon for IEEE 64-bit floating
	// point values.  This is the equality threshold recommended
	// in Numerical Recipes, if I recall correctly—it's small enough.
	const epsilon = 1.4901161193847656e-08

	for left > epsilon {
		a2 := a1 + sign*math.Min(math.Pi/2, left)
		c.partialArc(comp.Pos.X, comp.Pos.Y, comp.Radius, a1, a2, style)
		left -= math.Abs(a2 - a1)
		a1 = a2
	}
}

// Approximate a circular arc of fewer than π/2
// radians with cubic Bézier curve.
func (c *Canvas) partialArc(x, y, r vg.Length, a1, a2 float64, style string) {
	a := (a2 - a1) / 2
	x4 := r * vg.Length(math.Cos(a))
	y4 := r * vg.Length(math.Sin(a))
	x1 := x4
	y1 := -y4

	const k = 0.5522847498 // some magic constant
	f := k * vg.Length(math.Tan(a))
	x2 := x1 + f*y4
	y2 := y1 + f*x4
	x3 := x2
	y3 := -y2

	// Rotate and translate points into position.
	ar := a + a1
	sinar := vg.Length(math.Sin(ar))
	cosar := vg.Length(math.Cos(ar))
	x2r := x2*cosar - y2*sinar + x
	y2r := x2*sinar + y2*cosar + y
	x3r := x3*cosar - y3*sinar + x
	y3r := x3*sinar + y3*cosar + y
	x4 = r*vg.Length(math.Cos(a2)) + x
	y4 = r*vg.Length(math.Sin(a2)) + y
	c.pdf.Curve(c.unit(x2r), c.unit(y2r), c.unit(x3r), c.unit(y3r), c.unit(x4), c.unit(y4), style)
}

func (c *Canvas) pdfPointXY(x, y vg.Length) (float64, float64) {
	return c.unit(x), c.unit(y)
}

func (c *Canvas) pdfPoint(pt vg.Point) (float64, float64) {
	return c.unit(pt.X), c.unit(pt.Y)
}

// unit returns a fpdf.Unit, converted from a vg.Length.
func (c *Canvas) unit(l vg.Length) float64 {
	return l.Dots(c.DPI())
}

// imageName generates a unique image name for this PDF canvas
func (c *Canvas) imageName() string {
	c.nimgs++
	return fmt.Sprintf("image_%03d.png", c.nimgs)
}

// WriterCounter implements the io.Writer interface, and counts
// the total number of bytes written.
type writerCounter struct {
	io.Writer
	n int64
}

func (w *writerCounter) Write(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	w.n += int64(n)
	return n, err
}

// WriteTo writes the Canvas to an io.Writer.
// After calling Write, the canvas is closed
// and may no longer be used for drawing.
func (c *Canvas) WriteTo(w io.Writer) (int64, error) {
	c.Pop()
	c.pdf.Close()
	wc := writerCounter{Writer: w}
	b := bufio.NewWriter(&wc)
	if err := c.pdf.Output(b); err != nil {
		return wc.n, err
	}
	err := b.Flush()
	return wc.n, err
}

const (
	c255 = 255.0 / 65535.0
)

// rgb converts a Go color into a gofpdf 3-tuple int
func rgb(c color.Color) (int, int, int) {
	if c == nil {
		c = color.Black
	}
	r, g, b, _ := c.RGBA()
	return int(float64(r) * c255), int(float64(g) * c255), int(float64(b) * c255)
}
