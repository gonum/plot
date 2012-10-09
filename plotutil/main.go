// +build ignore

package main

import (
	"code.google.com/p/plotinum/plotutil"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/plot"
	"math/rand"
)

func main() {
	ExampleCentroids()
}

func ExampleCentroids() {
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

	mean95 := plotutil.NewCentroids(plotutil.MeanAndConf95, pts...)
	meanLine, meanPts := plotter.NewLinePoints(mean95)
	meanLine.Color = plotutil.Color(n+1)
	meanPts.Color = plotutil.Color(n+1)
	plt.Add(meanLine)
	meanXerr := plotter.NewXErrorBars(mean95)
	meanXerr.Color = plotutil.Color(n+1)
	meanYerr := plotter.NewYErrorBars(mean95)
	meanYerr.Color = plotutil.Color(n+1)
	plt.Add(meanLine, meanPts, meanXerr, meanYerr)

	medMinMax := plotutil.NewCentroids(plotutil.MedianAndMinMax, pts...)
	medLine, medPts := plotter.NewLinePoints(medMinMax)
	medLine.Color = plotutil.Color(n+2)
	medPts.Color = plotutil.Color(n+2)
	medXerr := plotter.NewXErrorBars(medMinMax)
	medXerr.Color = plotutil.Color(n+2)
	medYerr := plotter.NewYErrorBars(medMinMax)
	medYerr.Color = plotutil.Color(n+2)
	plt.Add(medLine, medPts, medXerr, medYerr)

	for i, p := range pts {
		scatter := plotter.NewScatter(p)
		scatter.Color = plotutil.Color(i)
		scatter.Shape = plotutil.Shape(i)
		plt.Add(scatter)
	}

	plt.Save(4, 4, "centroids.png")
}