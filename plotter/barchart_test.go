// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
)

func ExampleBarChart() {
	// Create the plot values and labels.
	values := Values{0.5, 10, 20, 30}
	verticalLabels := []string{"A", "B", "C", "D"}
	horizontalLabels := []string{"Label A", "Label B", "Label C", "Label D"}

	// Create a vertical BarChart
	p1, err := plot.New()
	handle(err)
	verticalBarChart, err := NewBarChart(values, 0.5*vg.Centimeter)
	handle(err)
	p1.Add(verticalBarChart)
	p1.NominalX(verticalLabels...)
	checkPlot("examplePlots", "verticalBarChart", "png", p1)

	// Create a horizontal BarChart
	p2, err := plot.New()
	handle(err)
	horizontalBarChart, err := NewBarChart(values, 0.5*vg.Centimeter)
	horizontalBarChart.Horizontal = true // Specify a horizontal BarChart.
	handle(err)
	p2.Add(horizontalBarChart)
	p2.NominalY(horizontalLabels...)
	checkPlot("examplePlots", "horizontalBarChart", "png", p2)

	// Output:
	// Image saved in dir examplePlots as verticalBarChart.png. Normally, you would use plot.Save().
	// Image saved in dir examplePlots as horizontalBarChart.png. Normally, you would use plot.Save().

}
