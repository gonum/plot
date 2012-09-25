// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

// +build ignore

// A simple test program to test plotters.
package main

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/vg"
	"fmt"
	"image/color"
	"math"
	"math/rand"
)

var examples = []struct{
	name string
	mkplot func()*plot.Plot
} {
	{ "example_logo", Example_logo },
	{ "example_functions", Example_functions },
	{ "example_boxPlots", Example_boxPlots },
	{ "example_verticalBoxPlots", Example_verticalBoxPlots },
	{ "example_horizontalBoxPlots", Example_horizontalBoxPlots },
	{ "example_points", Example_points },
	{ "example_errBars", Example_errBars },
	{ "example_bubbles", Example_bubbles },
	{ "example_histogram", Example_histogram },
	{ "example_barChart", Example_barChart },
	{ "example_stackedBarChart", Example_stackedBarChart },
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

func drawEps(name string, mkplot func()*plot.Plot) {
	if err := mkplot().Save(4, 4, name+".eps"); err != nil {
		panic(err)
	}
}

func drawPdf(name string, mkplot func()*plot.Plot) {
	if err := mkplot().Save(4, 4, name+".pdf"); err != nil {
		panic(err)
	}
}

func drawSvg(name string, mkplot func()*plot.Plot) {
	if err := mkplot().Save(4, 4, name+".svg"); err != nil {
		panic(err)
	}
}

func drawPng(name string, mkplot func()*plot.Plot) {
	if err := mkplot().Save(4, 4, name+".png"); err != nil {
		panic(err)
	}
}

func drawTiff(name string, mkplot func()*plot.Plot) {
	if err := mkplot().Save(4, 4, name+".tiff"); err != nil {
		panic(err)
	}
}

func drawJpg(name string, mkplot func()*plot.Plot) {
	if err := mkplot().Save(4, 4, name+".jpg"); err != nil {
		panic(err)
	}
}

// Draw the plotinum logo.
func Example_logo() *plot.Plot {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	plotter.DefaultLineStyle.Width = vg.Points(1)
	plotter.DefaultGlyphStyle.Radius = vg.Points(3)

	p.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{0, "0"}, {0.25, ""}, {0.5, "0.5"}, {0.75, ""}, {1, "1"},
	})
	p.X.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{0, "0"}, {0.25, ""}, {0.5, "0.5"}, {0.75, ""}, {1, "1"},
	})

	pts := plotter.XYs{{0, 0}, {0, 1}, {0.5, 1}, {0.5, 0.6}, {0, 0.6}}
	line := plotter.NewLine(pts)
	scatter := plotter.NewScatter(pts)
	p.Add(line, scatter)

	pts = plotter.XYs{{1, 0}, {0.75, 0}, {0.75, 0.75}}
	line = plotter.NewLine(pts)
	scatter = plotter.NewScatter(pts)
	p.Add(line, scatter)

	pts = plotter.XYs{{0.5, 0.5}, {1, 0.5}}
	line = plotter.NewLine(pts)
	scatter = plotter.NewScatter(pts)
	p.Add(line, scatter)

	return p
}

// Example_functions draws some functions.
func Example_functions() *plot.Plot {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Functions"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	quad := plotter.NewFunction(func(x float64) float64 { return x * x })
	quad.Color = color.RGBA{B: 255, A: 255}

	exp := plotter.NewFunction(func(x float64) float64 { return math.Pow(2, x) })
	exp.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	exp.Width = vg.Points(2)
	exp.Color = color.RGBA{G: 255, A: 255}

	sin := plotter.NewFunction(func(x float64) float64 { return 10*math.Sin(x) + 50 })
	sin.Dashes = []vg.Length{vg.Points(4), vg.Points(5)}
	sin.Width = vg.Points(4)
	sin.Color = color.RGBA{R: 255, A: 255}

	p.Add(quad, exp, sin)
	p.Legend.Add("x^2", quad)
	p.Legend.Add("2^x", exp)
	p.Legend.Add("10*sin(x)+50", sin)
	p.Legend.ThumbnailWidth = vg.Inches(0.5)

	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 100
	return p
}

// Example_boxPlots draws vertical boxplots.
func Example_boxPlots() *plot.Plot {
	rand.Seed(int64(0))
	n := 100
	uniform := make(plotter.Values, n)
	normal := make(plotter.Values, n)
	expon := make(plotter.Values, n)
	for i := 0; i < n; i++ {
		uniform[i] = rand.Float64()
		normal[i] = rand.NormFloat64()
		expon[i] = rand.ExpFloat64()
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Box Plot"
	p.Y.Label.Text = "plotter.Values"

	// Make boxes for our data and add them to the plot.
	p.Add(plotter.NewBoxPlot(vg.Points(20), 0, uniform),
		plotter.NewBoxPlot(vg.Points(20), 1, normal),
		plotter.NewBoxPlot(vg.Points(20), 2, expon))

	// Set the X axis of the plot to nominal with
	// the given names for x=0, x=1 and x=2.
	p.NominalX("Uniform\nDistribution", "Normal\nDistribution",
		"Exponential\nDistribution")
	return p
}

// Example_verticalBoxPlots draws vertical boxplots
// with some labels on their points.
func Example_verticalBoxPlots() *plot.Plot {
	rand.Seed(int64(0))
	n := 100
	uniform := make(valueLabels, n)
	normal := make(valueLabels, n)
	expon := make(valueLabels, n)
	for i := 0; i < n; i++ {
		uniform[i].Value = rand.Float64()
		uniform[i].Label = fmt.Sprintf("%4.4f", uniform[i].Value)
		normal[i].Value = rand.NormFloat64()
		normal[i].Label = fmt.Sprintf("%4.4f", normal[i].Value)
		expon[i].Value = rand.ExpFloat64()
		expon[i].Label = fmt.Sprintf("%4.4f", expon[i].Value)
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Box Plot"
	p.Y.Label.Text = "plotter.Values"

	// Make boxes for our data and add them to the plot.
	uniBox := plotter.NewBoxPlot(vg.Points(20), 0, uniform)
	uniLabels, err := uniBox.OutsideLabels(uniform)
	if err != nil {
		panic(err)
	}

	normBox := plotter.NewBoxPlot(vg.Points(20), 1, normal)
	normLabels, err := normBox.OutsideLabels(normal)
	if err != nil {
		panic(err)
	}

	expBox := plotter.NewBoxPlot(vg.Points(20), 2, expon)
	expLabels, err := expBox.OutsideLabels(expon)
	if err != nil {
		panic(err)
	}

	p.Add(uniBox, uniLabels, normBox, normLabels, expBox, expLabels)

	// Set the X axis of the plot to nominal with
	// the given names for x=0, x=1 and x=2.
	p.NominalX("Uniform\nDistribution", "Normal\nDistribution",
		"Exponential\nDistribution")
	return p
}

// valueLabels implements both the Valuer
// and Labellel interfaces.
type valueLabels []struct {
	Value float64
	Label string
}

func (vs valueLabels) Len() int {
	return len(vs)
}

func (vs valueLabels) Value(i int) float64 {
	return vs[i].Value
}

func (vs valueLabels) Label(i int) string {
	return vs[i].Label
}

// Example_horizontalBoxPlots draws horizontal boxplots
// with some labels on their points.
func Example_horizontalBoxPlots() *plot.Plot {
	rand.Seed(int64(0))
	n := 100
	uniform := make(valueLabels, n)
	normal := make(valueLabels, n)
	expon := make(valueLabels, n)
	for i := 0; i < n; i++ {
		uniform[i].Value = rand.Float64()
		uniform[i].Label = fmt.Sprintf("%4.4f", uniform[i].Value)
		normal[i].Value = rand.NormFloat64()
		normal[i].Label = fmt.Sprintf("%4.4f", normal[i].Value)
		expon[i].Value = rand.ExpFloat64()
		expon[i].Label = fmt.Sprintf("%4.4f", expon[i].Value)
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Horizontal Box Plot"
	p.X.Label.Text = "plotter.Values"

	// Make boxes for our data and add them to the plot.
	uniBox := plotter.HorizBoxPlot{plotter.NewBoxPlot(vg.Points(20), 0, uniform)}
	uniLabels, err := uniBox.OutsideLabels(uniform)
	if err != nil {
		panic(err)
	}

	normBox := plotter.HorizBoxPlot{plotter.NewBoxPlot(vg.Points(20), 1, normal)}
	normLabels, err := normBox.OutsideLabels(normal)
	if err != nil {
		panic(err)
	}

	expBox := plotter.HorizBoxPlot{plotter.NewBoxPlot(vg.Points(20), 2, expon)}
	expLabels, err := expBox.OutsideLabels(expon)
	if err != nil {
		panic(err)
	}
	p.Add(uniBox, uniLabels, normBox, normLabels, expBox, expLabels)

	// Add a GlyphBox plotter for debugging.
	p.Add(plotter.NewGlyphBoxes())

	// Set the Y axis of the plot to nominal with
	// the given names for y=0, y=1 and y=2.
	p.NominalY("Uniform\nDistribution", "Normal\nDistribution",
		"Exponential\nDistribution")
	return p
}

// Example_points draws some scatter points, a line,
// and a line with points.
func Example_points() *plot.Plot {
	rand.Seed(int64(0))

	n := 15
	scatterData := randomPoints(n)
	lineData := randomPoints(n)
	linePointsData := randomPoints(n)

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Points Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	s := plotter.NewScatter(scatterData)
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	s.GlyphStyle.Radius = vg.Points(3)

	l := plotter.NewLine(lineData)
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	lpLine, lpPoints := plotter.NewLinePoints(linePointsData)
	lpLine.Color = color.RGBA{G: 255, A: 255}
	lpPoints.Shape = plot.CircleGlyph{}
	lpPoints.Color = color.RGBA{R: 255, A: 255}

	p.Add(s, l, lpLine, lpPoints)
	p.Legend.Add("scatter", s)
	p.Legend.Add("line", l)
	p.Legend.Add("line points", lpLine, lpPoints)

	return p
}

// randomPoints returns some random x, y points
// with some interesting kind of trend.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}

// Example_errBars draws points and error bars.
func Example_errBars() *plot.Plot {

	type errPoints struct {
		plotter.XYs
		plotter.YErrors
		plotter.XErrors
	}

	rand.Seed(int64(0))
	n := 15
	data := errPoints{
		plotter.XYs:     randomPoints(n),
		plotter.YErrors: plotter.YErrors(randomError(n)),
		plotter.XErrors: plotter.XErrors(randomError(n)),
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	scatter := plotter.NewScatter(data)
	scatter.Shape = plot.CrossGlyph{}
	p.Add(scatter, plotter.NewXErrorBars(data), plotter.NewYErrorBars(data))
	p.Add(plotter.NewGlyphBoxes())

	return p
}

func randomError(n int) plotter.Errors {
	err := make(plotter.Errors, n)
	for i := range err {
		err[i].Low = rand.Float64()
		err[i].High = rand.Float64()
	}
	return err
}

func Example_bubbles() *plot.Plot {
	rand.Seed(int64(0))
	n := 10
	bubbleData := randomTriples(n)

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Bubbles"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	bs := plotter.NewBubbles(bubbleData, vg.Points(1), vg.Points(20))
	bs.Color = color.RGBA{R: 196, B: 128, A: 255}
	p.Add(bs)

	return p
}

// randomTriples returns some random x, y, z triples
// with some interesting kind of trend.
func randomTriples(n int) plotter.XYZs {
	data := make(plotter.XYZs, n)
	for i := range data {
		if i == 0 {
			data[i].X = rand.Float64()
		} else {
			data[i].X = data[i-1].X + 2*rand.Float64()
		}
		data[i].Y = data[i].X + 10*rand.Float64()
		data[i].Z = data[i].X
	}
	return data
}

// An example of making a histogram.
func Example_histogram() *plot.Plot {
	rand.Seed(int64(0))
	n := 10000
	vals := make(plotter.Values, n)
	for i := 0; i < n; i++ {
		vals[i] = rand.NormFloat64()
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Histogram"
	h := plotter.NewHist(vals, 16)
	h.Normalize(1)
	p.Add(h)

	// The normal distribution function
	norm := plotter.NewFunction(stdNorm)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)

	return p
}

// stdNorm returns the probability of drawing a
// value from a standard normal distribution.
func stdNorm(x float64) float64 {
	const sigma = 1.0
	const mu = 0.0
	const root2π = 2.50662827459517818309
	return 1.0 / (sigma * root2π) * math.Exp(-((x-mu)*(x-mu))/(2*sigma*sigma))
}

// An example of making a bar chart.
func Example_barChart() *plot.Plot {
	groupA := plotter.Values{20, 35, 30, 35, 27}
	groupB := plotter.Values{25, 32, 34, 20, 25}
	groupC := plotter.Values{12, 28, 15, 21, 8}
	groupD := plotter.Values{30, 42, 6, 9, 12}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Bar chart"
	p.Y.Label.Text = "Heights"

	w := vg.Points(8)

	barsA := plotter.NewBarChart(groupA, w)
	barsA.Color = color.RGBA{R: 255, A: 255}
	barsA.Offset = -w / 2

	barsB := plotter.NewBarChart(groupB, w)
	barsB.Color = color.RGBA{R: 196, G: 196, A: 255}
	barsB.Offset = w / 2

	barsC := plotter.NewBarChart(groupC, w)
	barsC.Color = color.RGBA{B: 255, A: 255}
	barsC.XMin = 6
	barsC.Offset = -w / 2

	barsD := plotter.NewBarChart(groupD, w)
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

	return p
}

// An example of making a stacked bar chart.
func Example_stackedBarChart() *plot.Plot {
	groupA := plotter.Values{20, 35, 30, 35, 27}
	groupB := plotter.Values{25, 32, 34, 20, 25}
	groupC := plotter.Values{12, 28, 15, 21, 8}
	groupD := plotter.Values{30, 42, 6, 9, 12}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Bar chart"
	p.Y.Label.Text = "Heights"

	w := vg.Points(15)

	barsA := plotter.NewBarChart(groupA, w)
	barsA.Color = color.RGBA{R: 255, A: 255}
	barsA.Offset = -w / 2

	barsB := plotter.NewBarChart(groupB, w)
	barsB.Color = color.RGBA{R: 196, G: 196, A: 255}
	barsB.StackOn(barsA)

	barsC := plotter.NewBarChart(groupC, w)
	barsC.Color = color.RGBA{B: 255, A: 255}
	barsC.Offset = w / 2

	barsD := plotter.NewBarChart(groupD, w)
	barsD.Color = color.RGBA{B: 255, R: 255, A: 255}
	barsD.StackOn(barsC)

	p.Add(barsA, barsB, barsC, barsD)
	p.Legend.Add("A", barsA)
	p.Legend.Add("B", barsB)
	p.Legend.Add("C", barsC)
	p.Legend.Add("D", barsD)
	p.Legend.Top = true
	p.NominalX("Zero", "One", "Two", "Three", "Four", "",
		"Six", "Seven", "Eight", "Nine", "Ten")

	return p
}
