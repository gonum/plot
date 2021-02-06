// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"log"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func ExampleContour() {
	rnd := rand.New(rand.NewSource(1234))

	const stddev = 2
	data := make([]float64, 6400)
	for i := range data {
		r := float64(i/80) - 40
		c := float64(i%80) - 40

		data[i] = rnd.NormFloat64()*stddev + math.Hypot(r, c)
	}

	var (
		grid   = unitGrid{mat.NewDense(80, 80, data)}
		levels = []float64{-1, 3, 7, 9, 13, 15, 19, 23, 27, 31}

		c = plotter.NewContour(
			grid,
			levels,
			palette.Rainbow(10, palette.Blue, palette.Red, 1, 1, 1),
		)
	)

	p := plot.New()
	p.Title.Text = "Contour"
	p.X.Padding = 0
	p.Y.Padding = 0
	p.X.Max = 79.5
	p.Y.Max = 79.5

	p.Add(c)

	err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/contour.png")
	if err != nil {
		log.Fatalf("could not save plot: %+v", err)
	}
}

type unitGrid struct{ mat.Matrix }

func (g unitGrid) Dims() (c, r int)   { r, c = g.Matrix.Dims(); return c, r }
func (g unitGrid) Z(c, r int) float64 { return g.Matrix.At(r, c) }
func (g unitGrid) X(c int) float64 {
	_, n := g.Matrix.Dims()
	if c < 0 || c >= n {
		panic("index out of range")
	}
	return float64(c)
}
func (g unitGrid) Y(r int) float64 {
	m, _ := g.Matrix.Dims()
	if r < 0 || r >= m {
		panic("index out of range")
	}
	return float64(r)
}

func TestContour(t *testing.T) {
	cmpimg.CheckPlotApprox(ExampleContour, t, 0.01, "contour.png")
}
