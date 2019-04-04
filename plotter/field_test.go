// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"image/png"
	"log"
	"os"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
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
		panic("index out of range")
	}
	return float64(c - f.c/2)
}
func (f field) Y(r int) float64 {
	if r < 0 || r >= f.r {
		panic("index out of range")
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
	f.DrawGlyph = func(c vg.Canvas, v plotter.XY) {
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func TestFieldGophers(t *testing.T) {
	cmpimg.CheckPlot(ExampleField_gophers, t, "gopher_field.png")
}
