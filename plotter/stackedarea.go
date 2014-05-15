package plotter

import (
	"errors"
	"image/color"

	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
)

type YData struct {
	// Data is the dataset, with a length matching the x-slice length.
	Data []float64

	// Color is the fill color for the area.
	color.Color

	// Label will appear on the plot as the dataset's label.
	Label string
}

// StackedArea implements the Plotter interface, drawing a stacked area plot,
// adding as many datasets as desired.
type StackedArea struct {
	// x is the shared data set for the X axis
	x []float64

	// ys contains all of the Y datasets
	ys []YData
}

// NewStackedArea creates as new stacked area plot plotter for the given data,
// with an x-axis.
func NewStackedArea(x []float64) (*StackedArea, error) {
	cpy := make([]float64, len(x))
	for i := range cpy {
		cpy[i] = x[i]
		if err := CheckFloats(cpy[i]); err != nil {
			return nil, err
		}
	}

	return &StackedArea{
		x:  cpy,
		ys: make([]YData, 0),
	}, nil
}

// Add adds a dataset to the plot
func (sa *StackedArea) Add(yd YData) error {
	if len(sa.x) != len(yd.Data) {
		return errors.New("XData length doesn't match X-axis length")
	}

	// Stack the data by adding the incoming data to the previous values
	stackedData := make([]float64, len(yd.Data))
	for i := range stackedData {
		stackedData[i] = yd.Data[i]
		if 0 != len(sa.ys) {
			stackedData[i] += sa.ys[len(sa.ys)-1].Data[i]
		}
	}

	sa.ys = append(sa.ys, YData{
		Data:  stackedData,
		Color: yd.Color,
		Label: yd.Label,
	})

	return nil
}

// AddToPlot adds all of the data to the given Plot.
func (sa *StackedArea) AddToPlot(p *plot.Plot) error {
	numPlots := len(sa.ys)
	if numPlots == 0 {
		return errors.New("No data has been added")
	}

	for i := numPlots - 1; i >= 0; i-- {
		// Make a line plotter and set its style.
		l, err := NewLine(sa.pointsForPlot(i))
		if err != nil {
			return err
		}

		l.LineStyle.Width = vg.Points(0)
		l.EnableShading(sa.ys[i].Color)

		p.Add(l)
		p.Legend.Add(sa.ys[i].Label, l)
	}

	p.Legend.Top = true
	p.Legend.Left = true

	return nil
}

// pointsForPlot returns the XY points for one of the StackedArea's datasets.
func (sa *StackedArea) pointsForPlot(n int) XYs {
	pts := make(XYs, len(sa.x))
	for i, x := range sa.x {
		pts[i].X = x
		pts[i].Y = sa.ys[n].Data[i]
	}
	return pts
}

// Normalize will normalize the current data. It should (obviously) be run after
// all data has been added.
func (sa *StackedArea) Normalize() {
	for i := range sa.x {
		total := sa.ys[len(sa.ys)-1].Data[i]
		for _, y := range sa.ys {
			y.Data[i] /= total
			y.Data[i] *= 100.0
		}
	}
}
