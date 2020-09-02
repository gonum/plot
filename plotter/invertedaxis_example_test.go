// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"image/color"
	"log"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func Example_invertedScale() {
	// This example is nearly identical to the LogScale, other than
	// both the X and Y axes are inverted. InvertedScale expects to act
	// on another Normalizer - which should allow for more flexibility
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

	// Neither .Min nor .Max for the X and Y axes are 'swapped'.
	// The minimal value is retained in .Min, and the maximal
	// value stays in .Max.
	p.X.Min = f.XMin
	p.X.Max = f.XMax
	p.Y.Min = math.Exp(f.XMin)
	p.Y.Max = math.Exp(f.XMax)

	err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/invertedlogscale.png")
	if err != nil {
		log.Panic(err)
	}
}
