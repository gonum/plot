// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

// vecimg implements the vg.Canvas interface
// using the draw2d package as a backend to output
// raster images.
package vecimg

import (
	"bufio"
	"code.google.com/p/draw2d/draw2d"
	"code.google.com/p/plotinum/vg"
	"fmt"
	"go/build"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

const (
	// dpi is the number of dots per inch.
	dpi = 96

	// importString is the current package import string.
	importString = "code.google.com/p/plotinum/vg/vecimg"
)

type Canvas struct {
	gc    draw2d.GraphicContext
	img   image.Image
	color color.Color

	// width is the current line width.
	width vg.Length
}

// New returns a new image canvas with the size specified.
// rounded up to the nearest pixel.
func New(width, height vg.Length) (*Canvas, error) {
	pkg, err := build.Import(importString, "", build.FindOnly)
	if err != nil {
		return nil, err
	}
	draw2d.SetFontFolder(filepath.Join(pkg.Dir, "fonts"))

	w := width.Inches() * dpi
	h := height.Inches() * dpi
	img := image.NewRGBA(image.Rect(0, 0, int(w+0.5), int(h+0.5)))
	gc := draw2d.NewGraphicContext(img)
	gc.SetDPI(dpi)
	gc.Scale(1, -1)
	gc.Translate(0, -h)
	c := &Canvas{gc: gc, img: img}
	vg.Initialize(c)
	return c, nil
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
	c.color = clr
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
	c.gc.Save()
}

func (c *Canvas) Pop() {
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
				comp.Start, comp.Start - comp.Finish)

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
	c.gc.Translate(x.Dots(c), (y + font.Extents().Ascent).Dots(c))
	c.gc.Scale(1, -1)
	c.gc.DrawImage(c.textImage(font, str))
	c.gc.Restore()
}

func (c *Canvas) textImage(font vg.Font, str string) *image.RGBA {
	w := font.Width(str).Dots(c)
	h := font.Extents().Height.Dots(c)
	img := image.NewRGBA(image.Rect(0, 0, int(w+0.5), int(h+0.5)))
	gc := draw2d.NewGraphicContext(img)

	gc.SetDPI(int(c.DPI()))
	gc.SetFillColor(c.color)
	data, ok := fontMap[font.Name()]
	if !ok {
		panic(fmt.Sprintf("Font name %s is unknown", font.Name()))
	}

	gc.SetFontData(data)
	gc.SetFontSize(font.Size.Points())
	gc.MoveTo(0, h+font.Extents().Descent.Dots(c))
	gc.FillString(str)

	return img
}

var (
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

func (c *Canvas) SavePNG(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	err = png.Encode(b, c.img)
	if err != nil {
		return err
	}
	return b.Flush()
}
