// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestBarChart(t *testing.T) {
	cmpimg.CheckPlot(ExampleBarChart, t, "verticalBarChart.png",
		"horizontalBarChart.png", "barChart2.png",
		"stackedBarChart.png")
}

func TestBarChart_positiveNegative(t *testing.T) {
	cmpimg.CheckPlot(ExampleBarChart_positiveNegative, t, "barChart_positiveNegative.png")
}
