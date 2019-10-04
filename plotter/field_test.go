// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/plotter"
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

func TestField(t *testing.T) {
	cmpimg.CheckPlot(ExampleField, t, "field.png")
}

func TestFieldColors(t *testing.T) {
	cmpimg.CheckPlot(ExampleField_colors, t, "color_field.png")
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
