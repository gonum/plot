// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"image/color"
	"log"
	"math"

	"golang.org/x/exp/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// An example of making a histogram.
func ExampleHistogram() {
	rnd := rand.New(rand.NewSource(1))

	// stdNorm returns the probability of drawing a
	// value from a standard normal distribution.
	stdNorm := func(x float64) float64 {
		const sigma = 1.0
		const mu = 0.0
		const root2π = 2.50662827459517818309
		return 1.0 / (sigma * root2π) * math.Exp(-((x-mu)*(x-mu))/(2*sigma*sigma))
	}

	n := 10000
	vals := make(plotter.Values, n)
	for i := 0; i < n; i++ {
		vals[i] = rnd.NormFloat64()
	}

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Histogram"
	h, err := plotter.NewHist(vals, 16)
	if err != nil {
		log.Panic(err)
	}
	h.Normalize(1)
	p.Add(h)

	// The normal distribution function
	norm := plotter.NewFunction(stdNorm)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)

	err = p.Save(200, 200, "testdata/histogram.png")
	if err != nil {
		log.Panic(err)
	}
}

func ExampleHistogram_logScaleY() {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Histogram in log-y"
	p.Y.Scale = plot.LogScale{}
	p.Y.Tick.Marker = plot.LogTicks{}
	p.Y.Label.Text = "Y"
	p.X.Label.Text = "X"

	h1, err := plotter.NewHist(plotter.Values{
		-2, -2,
		-1,
		+3, +3, +3, +3,
		+1, +1, +1, +1, +1, +1, +1, +1, +1, +1,
		+1, +1, +1, +1, +1, +1, +1, +1, +1, +1,
	}, 16)
	if err != nil {
		log.Fatal(err)
	}
	h1.LogY = true
	h1.FillColor = color.RGBA{255, 0, 0, 255}

	h2, err := plotter.NewHist(plotter.Values{
		-3, -3, -3,
		+2, +2, +2, +2, +2,
	}, 16)

	if err != nil {
		log.Fatal(err)
	}

	h2.LogY = true
	h2.FillColor = color.RGBA{0, 0, 255, 255}

	p.Add(h1, h2, plotter.NewGrid())

	err = p.Save(200, 200, "testdata/histogram_logy.png")
	if err != nil {
		log.Fatal(err)
	}
}
