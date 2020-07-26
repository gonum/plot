package plotter_test

import (
	"log"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
)

func ExampleRasterHeatMap() {
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
	raster := plotter.NewRasterHeatMap(&m, pal)
	plt.Add(raster)
	err = plt.Save(200, 200, "demoRaster.png")
	if err != nil {
		log.Panic(err)
	}
}
