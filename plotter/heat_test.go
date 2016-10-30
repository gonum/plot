// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"log"
	"testing"

	"github.com/gonum/matrix/mat64"
	"github.com/gonum/plot"
	"github.com/gonum/plot/internal/cmpimg"
	"github.com/gonum/plot/palette"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/recorder"
)

type offsetUnitGrid struct {
	XOffset, YOffset float64

	Data mat64.Matrix
}

func (g offsetUnitGrid) Dims() (c, r int)   { r, c = g.Data.Dims(); return c, r }
func (g offsetUnitGrid) Z(c, r int) float64 { return g.Data.At(r, c) }
func (g offsetUnitGrid) X(c int) float64 {
	_, n := g.Data.Dims()
	if c < 0 || c >= n {
		panic("index out of range")
	}
	return float64(c) + g.XOffset
}
func (g offsetUnitGrid) Y(r int) float64 {
	m, _ := g.Data.Dims()
	if r < 0 || r >= m {
		panic("index out of range")
	}
	return float64(r) + g.YOffset
}

func ExampleHeatMap() {
	m := offsetUnitGrid{
		XOffset: -2,
		YOffset: -1,
		Data: mat64.NewDense(3, 4, []float64{
			1, 2, 3, 4,
			5, 6, 7, 8,
			9, 10, 11, 12,
		})}
	h := NewHeatMap(m, palette.Heat(12, 1))

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Heat map"

	p.Add(h)

	p.X.Padding = 0
	p.Y.Padding = 0
	p.X.Max = 1.5
	p.Y.Max = 1.5

	err = p.Save(100, 100, "testdata/heatMap.png")
	if err != nil {
		log.Panic(err)
	}
}

func TestHeatMap(t *testing.T) {
	cmpimg.CheckPlot(ExampleHeatMap, t, "heatMap.png")
}

func TestFlatHeat(t *testing.T) {
	m := offsetUnitGrid{
		XOffset: -2,
		YOffset: -1,
		Data:    mat64.NewDense(3, 4, nil),
	}
	h := NewHeatMap(m, palette.Heat(12, 1))

	p, err := plot.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p.Add(h)

	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic for flat data")
			}
			const want = "heatmap: non-positive Z range"
			if r != want {
				t.Errorf("unexpected panic message: got:%q want:%q", r, want)
			}
		}()
		c := draw.NewCanvas(new(recorder.Canvas), 72, 72)
		p.Draw(c)
	}()
}
