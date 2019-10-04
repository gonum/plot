// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"fmt"
	"math"
	"testing"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type offsetUnitGrid struct {
	XOffset, YOffset float64

	Data mat.Matrix
}

func (g offsetUnitGrid) Dims() (c, r int)   { r, c = g.Data.Dims(); return c, r }
func (g offsetUnitGrid) Z(c, r int) float64 { return g.Data.At(r, c) }
func (g offsetUnitGrid) X(c int) float64 {
	_, n := g.Data.Dims()
	if c < 0 || c >= n {
		panic("column index out of range")
	}
	return float64(c) + g.XOffset
}
func (g offsetUnitGrid) Y(r int) float64 {
	m, _ := g.Data.Dims()
	if r < 0 || r >= m {
		panic("row index out of range")
	}
	return float64(r) + g.YOffset
}

type integerTicks struct{}

func (integerTicks) Ticks(min, max float64) []plot.Tick {
	var t []plot.Tick
	for i := math.Trunc(min); i <= max; i++ {
		t = append(t, plot.Tick{Value: i, Label: fmt.Sprint(i)})
	}
	return t
}

func TestHeatMap(t *testing.T) {
	cmpimg.CheckPlot(ExampleHeatMap, t, "heatMap.png")
}

func TestHeatMapDims(t *testing.T) {
	pal := palette.Heat(12, 1)

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

			m := offsetUnitGrid{Data: mat.NewDense(test.rows, test.cols, nil)}
			h := plotter.NewHeatMap(m, pal)

			p, err := plot.New()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			p.Add(h)

			img := vgimg.New(250, 175)
			dc := draw.New(img)

			p.Draw(dc)
		}()
	}
}
