// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"image/color"
	"log"

	"golang.org/x/exp/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func ExampleBarChart() {
	// Create the plot values and labels.
	values := plotter.Values{0.5, 10, 20, 30}
	verticalLabels := []string{"A", "B", "C", "D"}
	horizontalLabels := []string{"Label A", "Label B", "Label C", "Label D"}

	// Create a vertical BarChart
	p1, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	verticalBarChart, err := plotter.NewBarChart(values, 0.5*vg.Centimeter)
	if err != nil {
		log.Panic(err)
	}
	p1.Add(verticalBarChart)
	p1.NominalX(verticalLabels...)
	err = p1.Save(100, 100, "testdata/verticalBarChart.png")
	if err != nil {
		log.Panic(err)
	}

	// Create a horizontal BarChart
	p2, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	horizontalBarChart, err := plotter.NewBarChart(values, 0.5*vg.Centimeter)
	horizontalBarChart.Horizontal = true // Specify a horizontal BarChart.
	if err != nil {
		log.Panic(err)
	}
	p2.Add(horizontalBarChart)
	p2.NominalY(horizontalLabels...)
	err = p2.Save(100, 100, "testdata/horizontalBarChart.png")
	if err != nil {
		log.Panic(err)
	}

	// Now, make a different type of BarChart.
	groupA := plotter.Values{20, 35, 30, 35, 27}
	groupB := plotter.Values{25, 32, 34, 20, 25}
	groupC := plotter.Values{12, 28, 15, 21, 8}
	groupD := plotter.Values{30, 42, 6, 9, 12}

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Bar chart"
	p.Y.Label.Text = "Heights"

	w := vg.Points(8)

	barsA, err := plotter.NewBarChart(groupA, w)
	if err != nil {
		log.Panic(err)
	}
	barsA.Color = color.RGBA{R: 255, A: 255}
	barsA.Offset = -w / 2

	barsB, err := plotter.NewBarChart(groupB, w)
	if err != nil {
		log.Panic(err)
	}
	barsB.Color = color.RGBA{R: 196, G: 196, A: 255}
	barsB.Offset = w / 2

	barsC, err := plotter.NewBarChart(groupC, w)
	if err != nil {
		log.Panic(err)
	}
	barsC.XMin = 6
	barsC.Color = color.RGBA{B: 255, A: 255}
	barsC.Offset = -w / 2

	barsD, err := plotter.NewBarChart(groupD, w)
	if err != nil {
		log.Panic(err)
	}
	barsD.Color = color.RGBA{B: 255, R: 255, A: 255}
	barsD.XMin = 6
	barsD.Offset = w / 2

	p.Add(barsA, barsB, barsC, barsD)
	p.Legend.Add("A", barsA)
	p.Legend.Add("B", barsB)
	p.Legend.Add("C", barsC)
	p.Legend.Add("D", barsD)
	p.Legend.Top = true
	p.NominalX("Zero", "One", "Two", "Three", "Four", "",
		"Six", "Seven", "Eight", "Nine", "Ten")

	p.Add(plotter.NewGlyphBoxes())
	err = p.Save(300, 250, "testdata/barChart2.png")
	if err != nil {
		log.Panic(err)
	}

	// Now, make a stacked BarChart.
	p, err = plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Bar chart"
	p.Y.Label.Text = "Heights"

	w = vg.Points(15)

	barsA, err = plotter.NewBarChart(groupA, w)
	if err != nil {
		log.Panic(err)
	}
	barsA.Color = color.RGBA{R: 255, A: 255}
	barsA.Offset = -w / 2

	barsB, err = plotter.NewBarChart(groupB, w)
	if err != nil {
		log.Panic(err)
	}
	barsB.Color = color.RGBA{R: 196, G: 196, A: 255}
	barsB.StackOn(barsA)

	barsC, err = plotter.NewBarChart(groupC, w)
	if err != nil {
		log.Panic(err)
	}
	barsC.Offset = w / 2
	barsC.Color = color.RGBA{B: 255, A: 255}

	barsD, err = plotter.NewBarChart(groupD, w)
	if err != nil {
		log.Panic(err)
	}
	barsD.StackOn(barsC)
	barsD.Color = color.RGBA{B: 255, R: 255, A: 255}

	p.Add(barsA, barsB, barsC, barsD)
	p.Legend.Add("A", barsA)
	p.Legend.Add("B", barsB)
	p.Legend.Add("C", barsC)
	p.Legend.Add("D", barsD)
	p.Legend.Top = true
	p.NominalX("Zero", "One", "Two", "Three", "Four")

	p.Add(plotter.NewGlyphBoxes())
	err = p.Save(250, 250, "testdata/stackedBarChart.png")
	if err != nil {
		log.Panic(err)
	}
}

// This example shows a bar chart with both positive and negative values.
func ExampleBarChart_positiveNegative() {
	rnd := rand.New(rand.NewSource(1))

	// Create random data points between -1 and 1.
	const n = 6
	data1 := make(plotter.Values, n)
	data2 := make(plotter.Values, n)
	net := make(plotter.XYs, n) // net = data1 + data2
	for i := 0; i < n; i++ {
		data1[i] = rnd.Float64()*2 - 1
		data2[i] = rnd.Float64()*2 - 1
		net[i].X = data1[i] + data2[i]
		net[i].Y = float64(i)
	}

	// splitBySign splits an array into two arrays containing the positive and
	// negative values, respectively, from the original array.
	splitBySign := func(d plotter.Values) (pos, neg plotter.Values) {
		pos = make(plotter.Values, len(d))
		neg = make(plotter.Values, len(d))
		for i, v := range d {
			if v > 0 {
				pos[i] = v
			} else {
				neg[i] = v
			}
		}
		return
	}

	data1Pos, data1Neg := splitBySign(data1)
	data2Pos, data2Neg := splitBySign(data2)

	const barWidth = 0.3 * vg.Centimeter
	pos1, err := plotter.NewBarChart(data1Pos, barWidth)
	if err != nil {
		log.Panic(err)
	}
	pos2, err := plotter.NewBarChart(data2Pos, barWidth)
	if err != nil {
		log.Panic(err)
	}
	neg1, err := plotter.NewBarChart(data1Neg, barWidth)
	if err != nil {
		log.Panic(err)
	}
	neg2, err := plotter.NewBarChart(data2Neg, barWidth)
	if err != nil {
		log.Panic(err)
	}

	netDots, err := plotter.NewScatter(net)
	if err != nil {
		log.Panic(err)
	}
	netDots.Radius = vg.Points(1.25)

	pos2.StackOn(pos1) // Specify that pos2 goes on top of pos1.
	neg2.StackOn(neg1) // Specify that neg2 goes on top of neg1.

	color1 := color.NRGBA{R: 112, G: 22, B: 0, A: 255}
	color2 := color.NRGBA{R: 91, G: 194, B: 54, A: 100}

	pos1.Color, neg1.Color = color1, color1
	pos2.Color, neg2.Color = color2, color2

	// Specify that we want a horizontal bar chart.
	pos1.Horizontal, pos2.Horizontal, neg1.Horizontal, neg2.Horizontal = true, true, true, true

	// Create a line at zero.
	zero, err := plotter.NewLine(plotter.XYs{{0, 0}, {0, 5}})
	if err != nil {
		log.Panic(err)
	}

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Add(zero, pos1, pos2, neg1, neg2, netDots)
	p.NominalY("Alpha", "Bravo", "Charlie", "Echo", "Foxtrot", "Golf")

	p.Legend.Add("1", pos1)
	p.Legend.Add("2", pos2)
	p.Legend.Add("Sum", netDots)
	p.Legend.Top = true
	p.Legend.ThumbnailWidth = 2 * vg.Millimeter
	p.Legend.Font.Size = 2 * vg.Millimeter

	err = p.Save(100, 100, "testdata/barChart_positiveNegative.png")
	if err != nil {
		log.Panic(err)
	}
}
