// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

// vgimg implements the vg.Canvas interface using
// draw2d (code.google.com/p/draw2d/draw2d)
// as a backend to output raster images.
package vgimg

import (
	"bufio"
	"code.google.com/p/draw2d/draw2d"
	"code.google.com/p/go.image/tiff"
	"code.google.com/p/plotinum/vg"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
)

// dpi is the number of dots per inch.
const dpi = 96

// Canvas implements the vg.Canvas interface,
// drawing to an image.Image using draw2d.
type Canvas struct {
	gc    draw2d.GraphicContext
	img   image.Image
	w, h  vg.Length
	color []color.Color

	// width is the current line width.
	width vg.Length
}

// New returns a new image canvas with
// the size specified  rounded up to the
// nearest pixel.
func New(width, height vg.Length) *Canvas {
	w := width.Inches() * dpi
	h := height.Inches() * dpi
	img := image.NewRGBA(image.Rect(0, 0, int(w+0.5), int(h+0.5)))

	return NewImage(img)
}

// NewImage returns a new image canvas
// that draws to the given image.  The
// minimum point of the given image
// should probably be 0,0.
func NewImage(img draw.Image) *Canvas {
	w := float64(img.Bounds().Max.X - img.Bounds().Min.X)
	h := float64(img.Bounds().Max.Y - img.Bounds().Min.Y)
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)
	gc := draw2d.NewGraphicContext(img)
	gc.SetDPI(dpi)
	gc.Scale(1, -1)
	gc.Translate(0, -h)
	c := &Canvas{
		gc:    gc,
		img:   img,
		w:     vg.Inches(w / dpi),
		h:     vg.Inches(h / dpi),
		color: []color.Color{color.Black},
	}
	vg.Initialize(c)
	return c
}

func (c *Canvas) Size() (w, h vg.Length) {
	return c.w, c.h
}

func (c *Canvas) SetLineWidth(w vg.Length) {
	c.width = w
	c.gc.SetLineWidth(w.Dots(c))
}

func (c *Canvas) SetLineDash(ds []vg.Length, offs vg.Length) {
	dashes := make([]float64, len(ds))
	for i, d := range ds {
		dashes[i] = d.Dots(c)
	}
	c.gc.SetLineDash(dashes, offs.Dots(c))
}

func (c *Canvas) SetColor(clr color.Color) {
	if clr == nil {
		clr = color.Black
	}
	c.gc.SetFillColor(clr)
	c.gc.SetStrokeColor(clr)
	c.color[len(c.color)-1] = clr
}

func (c *Canvas) Rotate(t float64) {
	c.gc.Rotate(t)
}

func (c *Canvas) Translate(x, y vg.Length) {
	c.gc.Translate(x.Dots(c), y.Dots(c))
}

func (c *Canvas) Scale(x, y float64) {
	c.gc.Scale(x, y)
}

func (c *Canvas) Push() {
	c.color = append(c.color, c.color[len(c.color)-1])
	c.gc.Save()
}

func (c *Canvas) Pop() {
	c.color = c.color[:len(c.color)-1]
	c.gc.Restore()
}

func (c *Canvas) Stroke(p vg.Path) {
	if c.width == 0 {
		return
	}
	c.outline(p)
	c.gc.Stroke()
}

func (c *Canvas) Fill(p vg.Path) {
	c.outline(p)
	c.gc.Fill()
}

func (c *Canvas) outline(p vg.Path) {
	c.gc.BeginPath()
	for _, comp := range p {
		switch comp.Type {
		case vg.MoveComp:
			c.gc.MoveTo(comp.X.Dots(c), comp.Y.Dots(c))

		case vg.LineComp:
			c.gc.LineTo(comp.X.Dots(c), comp.Y.Dots(c))

		case vg.ArcComp:
			c.gc.ArcTo(comp.X.Dots(c), comp.Y.Dots(c),
				comp.Radius.Dots(c), comp.Radius.Dots(c),
				comp.Start, comp.Angle)

		case vg.CloseComp:
			c.gc.Close()

		default:
			panic(fmt.Sprintf("Unknown path component: %d", comp.Type))
		}
	}
}

func (c *Canvas) DPI() float64 {
	return float64(c.gc.GetDPI())
}

func (c *Canvas) FillString(font vg.Font, x, y vg.Length, str string) {
	c.gc.Save()
	defer c.gc.Restore()

	data, ok := fontMap[font.Name()]
	if !ok {
		panic(fmt.Sprintf("Font name %s is unknown", font.Name()))
	}
	if !registeredFont[font.Name()] {
		draw2d.RegisterFont(data, font.Font())
		registeredFont[font.Name()] = true
	}
	c.gc.SetFontData(data)
	c.gc.Translate(x.Dots(c), y.Dots(c))
	c.gc.Scale(1, -1)
	c.gc.FillString(str)
}

var (
	// RegisteredFont contains the set of font names
	// that have already been registered with draw2d.
	registeredFont = map[string]bool{}

	// FontMap contains a mapping from vg's font
	// names to draw2d.FontData for the corresponding
	// font.  This is needed to register the  fonts with
	// draw2d.
	fontMap = map[string]draw2d.FontData{
		"Courier": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilyMono,
			Style:  draw2d.FontStyleNormal,
		},
		"Courier-Bold": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilyMono,
			Style:  draw2d.FontStyleBold,
		},
		"Courier-Oblique": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilyMono,
			Style:  draw2d.FontStyleItalic,
		},
		"Courier-BoldOblique": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilyMono,
			Style:  draw2d.FontStyleItalic | draw2d.FontStyleBold,
		},
		"Helvetica": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleNormal,
		},
		"Helvetica-Bold": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleBold,
		},
		"Helvetica-Oblique": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleItalic,
		},
		"Helvetica-BoldOblique": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleItalic | draw2d.FontStyleBold,
		},
		"Times-Roman": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySerif,
			Style:  draw2d.FontStyleNormal,
		},
		"Times-Bold": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySerif,
			Style:  draw2d.FontStyleBold,
		},
		"Times-Italic": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySerif,
			Style:  draw2d.FontStyleItalic,
		},
		"Times-BoldItalic": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySerif,
			Style:  draw2d.FontStyleItalic | draw2d.FontStyleBold,
		},
	}
)

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

// A JpegCanvas is an image canvas with a WriteTo method
// that writes a jpeg image.
type JpegCanvas struct {
	*Canvas
}

// WriteTo implements the io.WriterTo interface, writing a jpeg image.
func (c JpegCanvas) WriteTo(w io.Writer) (int64, error) {
	wc := writerCounter{Writer: w}
	b := bufio.NewWriter(&wc)
	if err := jpeg.Encode(b, c.img, nil); err != nil {
		return wc.n, err
	}
	err := b.Flush()
	return wc.n, err
}

// A PngCanvas is an image canvas with a WriteTo method that
// writes a png image.
type PngCanvas struct {
	*Canvas
}

// WriteTo implements the io.WriterTo interface, writing a png image.
func (c PngCanvas) WriteTo(w io.Writer) (int64, error) {
	wc := writerCounter{Writer: w}
	b := bufio.NewWriter(&wc)
	if err := png.Encode(b, c.img); err != nil {
		return wc.n, err
	}
	err := b.Flush()
	return wc.n, err
}

// A TiffCanvas is an image canvas with a WriteTo method that
// writes a tiff image.
type TiffCanvas struct {
	*Canvas
}

// WriteTo implements the io.WriterTo interface, writing a tiff image.
func (c TiffCanvas) WriteTo(w io.Writer) (int64, error) {
	wc := writerCounter{Writer: w}
	b := bufio.NewWriter(&wc)
	if err := tiff.Encode(b, c.img, nil); err != nil {
		return wc.n, err
	}
	err := b.Flush()
	return wc.n, err
}
