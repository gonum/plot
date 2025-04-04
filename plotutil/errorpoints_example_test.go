// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotutil_test

import (
	"math/rand/v2"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

func ExampleErrorPoints() {
	rnd := rand.New(rand.NewPCG(1, 1))

	// Get some random data.
	n, m := 5, 10
	pts := make([]plotter.XYer, n)
	for i := range pts {
		xys := make(plotter.XYs, m)
		pts[i] = xys
		center := float64(i)
		for j := range xys {
			xys[j].X = center + (rnd.Float64() - 0.5)
			xys[j].Y = center + (rnd.Float64() - 0.5)
		}
	}

	plt := plot.New()

	mean95, err := plotutil.NewErrorPoints(plotutil.MeanAndConf95, pts...)
	if err != nil {
		panic(err)
	}
	medMinMax, err := plotutil.NewErrorPoints(plotutil.MedianAndMinMax, pts...)
	if err != nil {
		panic(err)
	}
	err = plotutil.AddLinePoints(plt,
		"mean and 95% confidence", mean95,
		"median and minimum and maximum", medMinMax)
	if err != nil {
		panic(err)
	}
	if err := plotutil.AddErrorBars(plt, mean95, medMinMax); err != nil {
		panic(err)
	}
	if err := plotutil.AddScatters(plt, pts[0], pts[1], pts[2], pts[3], pts[4]); err != nil {
		panic(err)
	}

	err = plt.Save(4, 4, "centroids.png")
	if err != nil {
		panic(err)
	}
}
