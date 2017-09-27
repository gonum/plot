// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/internal/cmpimg"
	"gonum.org/v1/plot/plotter"
	"log"
	"testing"
)

func TestFloatPrecision(t *testing.T) {
	const fname = "precision.png"

	cmpimg.CheckPlot(func() {
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}

		p.X.Label.Text = "x"
		p.Y.Label.Text = "y"

		var data = make(plotter.XYs, 10)
		for i := range data {
			data[i].X = float64(i)
			data[i].Y = 1300
		}

		lines, points, err := plotter.NewLinePoints(data)
		if err != nil {
			log.Fatal(err)
		}
		p.Add(points, lines)
		p.Add(plotter.NewGrid())

		err = p.Save(300, 300, "testdata/"+fname)
		if err != nil {
			log.Fatal(err)
		}
	}, t, fname)
}
