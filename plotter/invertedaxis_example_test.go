// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"image/color"
	"log"
	"math"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/internal/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Example_invertedScale shows how to create a plot with an inverted Y-axis.
// This is nearly identical to the Log example, except it inverts the X and Y axes
func Example_invertedScale() {
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	p.Title.Text = "Example of inverted axes"
	p.Y.Scale = plot.InvertedScale{Normalizer: plot.LogScale{}}
	p.X.Scale = plot.InvertedScale{Normalizer: plot.LinearScale{}}
	p.Y.Tick.Marker = plot.LogTicks{}
	p.X.Label.Text = "x"
	p.Y.Label.Text = "f(x)"

	f := plotter.NewFunction(math.Exp)
	f.XMin = 0.2
	f.XMax = 10

	f.Color = color.RGBA{R: 255, A: 255}

	p.Add(f, plotter.NewGrid())
	p.Legend.Add("exp(x)", f)

	// Notice that both .Min and .Max for X and Y are both in 'normal' order
	p.X.Min = f.XMin
	p.X.Max = f.XMax
	p.Y.Min = math.Exp(f.XMin)
	p.Y.Max = math.Exp(f.XMax)

	err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/invertedlogscale.png")
	if err != nil {
		log.Panic(err)
	}
}

func TestInvertedScale(t *testing.T) {
	cmpimg.CheckPlot(Example_invertedScale, t, "invertedlogscale.png")
}
