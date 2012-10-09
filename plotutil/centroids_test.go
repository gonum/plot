package plotutil

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"math/rand"
)

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

	cs := NewCentroids(MeanAndConf95, pts...)

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}

	plt.Add(plotter.NewLinePoints(cs))
	plt.Add(plotter.NewXErrorBars(cs), plotter.NewYErrorBars(cs))

	for i, p := range pts {
		scatter := plotter.NewScatter(p)
		scatter.Color = Color(i)
		scatter.Shape = Shape(i)
		plt.Add(scatter)
	}

	plt.Save(4, 4, "centroids.png")
}
