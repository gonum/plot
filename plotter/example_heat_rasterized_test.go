// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"log"
	"os"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func ExampleHeatMap_rasterized() {
	m := offsetUnitGrid{
		XOffset: -2,
		YOffset: -1,
		Data: mat.NewDense(3, 4, []float64{
			1, 2, 3, 4,
			5, 6, 7, 8,
			9, 10, 11, 12,
		})}

	pal := palette.Heat(12, 1)
	plt, err := plot.New()
	if err != nil {
		log.Panic(err)
	}

	raster := plotter.NewHeatMap(&m, pal)
	raster.Rasterized = true
	plt.Add(raster)

	err := plt.Save(5*vg.Centimeter, 3.5*vg.Centimeter, "testdata/rasterHeatMap.png")
	if err != nil {
		log.Panic(err)
	}
}
