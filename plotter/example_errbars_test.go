// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"log"

	"golang.org/x/exp/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

// ExampleErrors draws points and error bars.
func ExampleErrors() {
	rnd := rand.New(rand.NewSource(1))

	randomError := func(n int) plotter.Errors {
		err := make(plotter.Errors, n)
		for i := range err {
			err[i].Low = rnd.Float64()
			err[i].High = rnd.Float64()
		}
		return err
	}
	// randomPoints returns some random x, y points
	// with some interesting kind of trend.
	randomPoints := func(n int) plotter.XYs {
		pts := make(plotter.XYs, n)
		for i := range pts {
			if i == 0 {
				pts[i].X = rnd.Float64()
			} else {
				pts[i].X = pts[i-1].X + rnd.Float64()
			}
			pts[i].Y = pts[i].X + 10*rnd.Float64()
		}
		return pts
	}

	type errPoints struct {
		plotter.XYs
		plotter.YErrors
		plotter.XErrors
	}

	n := 15
	data := errPoints{
		XYs:     randomPoints(n),
		YErrors: plotter.YErrors(randomError(n)),
		XErrors: plotter.XErrors(randomError(n)),
	}

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	scatter, err := plotter.NewScatter(data)
	if err != nil {
		log.Panic(err)
	}
	scatter.Shape = draw.CrossGlyph{}
	xerrs, err := plotter.NewXErrorBars(data)
	if err != nil {
		log.Panic(err)
	}
	yerrs, err := plotter.NewYErrorBars(data)
	if err != nil {
		log.Panic(err)
	}
	p.Add(scatter, xerrs, yerrs)

	err = p.Save(200, 200, "testdata/errorBars.png")
	if err != nil {
		log.Panic(err)
	}
}
