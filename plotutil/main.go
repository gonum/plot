// +build ignore

package main

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/plotutil"
	"math/rand"
)

var examples = []struct {
	name   string
	mkplot func() *plot.Plot
}{
	{"example_errpoints", Example_errpoints},
}

func main() {
	for _, ex := range examples {
		drawEps(ex.name, ex.mkplot)
		drawSvg(ex.name, ex.mkplot)
		drawPng(ex.name, ex.mkplot)
		drawTiff(ex.name, ex.mkplot)
		drawJpg(ex.name, ex.mkplot)
		drawPdf(ex.name, ex.mkplot)
	}
}

func drawEps(name string, mkplot func() *plot.Plot) {
	if err := mkplot().Save(4, 4, name+".eps"); err != nil {
		panic(err)
	}
}

func drawPdf(name string, mkplot func() *plot.Plot) {
	if err := mkplot().Save(4, 4, name+".pdf"); err != nil {
		panic(err)
	}
}

func drawSvg(name string, mkplot func() *plot.Plot) {
	if err := mkplot().Save(4, 4, name+".svg"); err != nil {
		panic(err)
	}
}

func drawPng(name string, mkplot func() *plot.Plot) {
	if err := mkplot().Save(4, 4, name+".png"); err != nil {
		panic(err)
	}
}

func drawTiff(name string, mkplot func() *plot.Plot) {
	if err := mkplot().Save(4, 4, name+".tiff"); err != nil {
		panic(err)
	}
}

func drawJpg(name string, mkplot func() *plot.Plot) {
	if err := mkplot().Save(4, 4, name+".jpg"); err != nil {
		panic(err)
	}
}

// Example_errpoints draws some error points.
func Example_errpoints() *plot.Plot {
	// Get some random data.
	n, m := 5, 10
	pts := make([]plotter.XYer, n)
	for i := range pts {
		xys := make(plotter.XYs, m)
		pts[i] = xys
		center := float64(i)
		for j := range xys {
			xys[j].X = center + (rand.Float64() - 0.5)
			xys[j].Y = center + (rand.Float64() - 0.5)
		}
	}

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}

	mean95, err := plotutil.NewErrorPoints(plotutil.MeanAndConf95, pts...)
	if err != nil {
		panic(err)
	}
	medMinMax, err := plotutil.NewErrorPoints(plotutil.MedianAndMinMax, pts...)
	if err != nil {
		panic(err)
	}
	plotutil.AddLinePoints(plt,
		"mean and 95% confidence", mean95,
		"median and minimum and maximum", medMinMax)
	if err := plotutil.AddErrorBars(plt, mean95, medMinMax); err != nil {
		panic(err)
	}
	if err := plotutil.AddScatters(plt, pts[0], pts[1], pts[2], pts[3], pts[4]); err != nil {
		panic(err)
	}

	return plt
}
