// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type field struct {
	c, r int
	fn   func(x, y float64) plotter.XY
}

func (f field) Dims() (c, r int)           { return f.c, f.r }
func (f field) Vector(c, r int) plotter.XY { return f.fn(f.X(c), f.Y(r)) }
func (f field) X(c int) float64 {
	if c < 0 || c >= f.c {
		panic("column index out of range")
	}
	return float64(c - f.c/2)
}
func (f field) Y(r int) float64 {
	if r < 0 || r >= f.r {
		panic("row index out of range")
	}
	return float64(r - f.r/2)
}

func ExampleField() {
	f := plotter.NewField(field{
		r: 17, c: 19,
		fn: func(x, y float64) plotter.XY {
			return plotter.XY{
				X: y,
				Y: -x,
			}
		},
	})
	f.LineStyle.Width = 0.2

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Vector field"

	p.X.Tick.Marker = integerTicks{}
	p.Y.Tick.Marker = integerTicks{}

	p.Add(f)

	img := vgimg.New(250, 175)
	dc := draw.New(img)

	p.Draw(dc)
	w, err := os.Create("testdata/field.png")
	if err != nil {
		log.Panic(err)
	}
	png := vgimg.PngCanvas{Canvas: img}
	if _, err = png.WriteTo(w); err != nil {
		log.Panic(err)
	}
}

func TestField(t *testing.T) {
	cmpimg.CheckPlot(ExampleField, t, "field.png")
}

func ExampleField_colors() {
	f := plotter.NewField(field{
		r: 5, c: 9,
		fn: func(x, y float64) plotter.XY {
			return plotter.XY{
				X: -0.75*x + y,
				Y: -0.75*y - x,
			}
		},
	})

	pal := moreland.ExtendedBlackBody()
	pal.SetMin(0)
	pal.SetMax(1.1) // Use 1.1 to make highest magnitude vectors visible on white.

	// Provide a DrawGlyph function to render a custom
	// vector instead of the default monochrome arrow.
	f.DrawGlyph = func(c vg.Canvas, sty draw.LineStyle, v plotter.XY) {
		c.Push()
		defer c.Pop()
		mag := math.Hypot(v.X, v.Y)
		var pa vg.Path
		if mag == 0 {
			// Draw a black dot for zero vectors.
			c.SetColor(color.Black)
			pa.Move(vg.Point{X: sty.Width})
			pa.Arc(vg.Point{}, sty.Width, 0, 2*math.Pi)
			pa.Close()
			c.Fill(pa)
			return
		}
		// Choose a color from the palette for the magnitude.
		col, err := pal.At(mag)
		if err != nil {
			panic(err)
		}
		c.SetColor(col)
		pa.Move(vg.Point{})
		pa.Line(vg.Point{X: 1, Y: 0})
		pa.Close()
		c.Stroke(pa)
	}

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Vortex"

	p.X.Tick.Marker = integerTicks{}
	p.Y.Tick.Marker = integerTicks{}

	p.Add(f)

	img := vgimg.New(250, 175)
	dc := draw.New(img)

	p.Draw(dc)
	w, err := os.Create("testdata/color_field.png")
	if err != nil {
		log.Panic(err)
	}
	png := vgimg.PngCanvas{Canvas: img}
	if _, err = png.WriteTo(w); err != nil {
		log.Panic(err)
	}
}

func TestFieldColors(t *testing.T) {
	cmpimg.CheckPlot(ExampleField_colors, t, "color_field.png")
}

func ExampleField_gophers() {
	file, err := os.Open("testdata/gopher_running.png")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	gopher, err := png.Decode(file)
	if err != nil {
		log.Panic(err)
	}

	f := plotter.NewField(field{
		r: 5, c: 9,
		fn: func(x, y float64) plotter.XY {
			return plotter.XY{
				X: -0.75*x + y,
				Y: -0.75*y - x,
			}
		},
	})

	// Provide a DrawGlyph function to render a custom
	// vector glyph instead of the default arrow.
	f.DrawGlyph = func(c vg.Canvas, _ draw.LineStyle, v plotter.XY) {
		// The canvas is unscaled if the vector has a zero
		// magnitude, so return in that case.
		if math.Hypot(v.X, v.Y) == 0 {
			return
		}
		// Vector glyphs are scaled to half unit length by the
		// plotter, so scale the gopher to twice unit size so
		// it fits the cell, and center on the cell.
		c.Translate(vg.Point{X: -1, Y: -1})
		c.DrawImage(vg.Rectangle{Max: vg.Point{X: 2, Y: 2}}, gopher)
	}

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Gopher vortex"

	p.X.Tick.Marker = integerTicks{}
	p.Y.Tick.Marker = integerTicks{}

	p.Add(f)

	img := vgimg.New(250, 175)
	dc := draw.New(img)

	p.Draw(dc)
	w, err := os.Create("testdata/gopher_field.png")
	if err != nil {
		log.Panic(err)
	}
	png := vgimg.PngCanvas{Canvas: img}
	if _, err = png.WriteTo(w); err != nil {
		log.Panic(err)
	}
}

func TestFieldGophers(t *testing.T) {
	cmpimg.CheckPlot(ExampleField_gophers, t, "gopher_field.png")
}

func TestFieldDims(t *testing.T) {
	for _, test := range []struct {
		rows int
		cols int
	}{
		{rows: 1, cols: 2},
		{rows: 2, cols: 1},
		{rows: 2, cols: 2},
	} {
		func() {
			defer func() {
				r := recover()
				if r != nil {
					t.Errorf("unexpected panic for rows=%d cols=%d: %v", test.rows, test.cols, r)
				}
			}()

			f := plotter.NewField(field{
				r: test.rows, c: test.cols,
				fn: func(x, y float64) plotter.XY {
					return plotter.XY{
						X: y,
						Y: -x,
					}
				},
			})

			p, err := plot.New()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			p.Add(f)

			img := vgimg.New(250, 175)
			dc := draw.New(img)

			p.Draw(dc)
		}()
	}
}
