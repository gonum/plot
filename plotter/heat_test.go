// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"log"
	"testing"

	"github.com/gonum/matrix/mat64"
	"github.com/gonum/plot"
	"github.com/gonum/plot/palette"
)

func ExampleHeatMap() {
	m := unitGrid{mat64.NewDense(3, 4, []float64{
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
	p.X.Max = 3.5
	p.Y.Max = 2.5

	err = p.Save(100, 100, "testdata/heatMap.png")
	if err != nil {
		log.Panic(err)
	}
}

func TestHeatMap(t *testing.T) {
	checkPlot(ExampleHeatMap, t, "heatMap.png")
}
