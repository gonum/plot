// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"image/color"
	"log"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/internal/cmpimg"
)

func filledLine() {
	rnd := rand.New(rand.NewSource(1))

	// randomPoints returns some random x, y points
	// with some interesting kind of trend.
	randomPoints := func(n int, x float64) XYs {
		pts := make(XYs, n)
		for i := range pts {
			if i == 0 {
				pts[i].X = x + rnd.Float64()
			} else {
				pts[i].X = pts[i-1].X + 0.5 + rnd.Float64()
			}
			pts[i].Y = -5. + 10*rnd.Float64()
		}
		return pts
	}

	n := 4

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Filled Line Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(NewGrid())

	filled, err := NewLine(randomPoints(n, 0))
	if err != nil {
		log.Panic(err)
	}
	filled.FillColor = color.RGBA{R: 196, G: 255, B: 196, A: 255}

	p.Add(filled)
	p.Y.Max = -1

	err = p.Save(200, 200, "testdata/filledLine.png")
	if err != nil {
		log.Panic(err)
	}
}

func TestFilledLine(t *testing.T) {
	cmpimg.CheckPlot(filledLine, t, "filledLine.png")
}
