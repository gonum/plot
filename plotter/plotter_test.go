// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"testing"
)

func TestDrawPng(t *testing.T) {
	if err := Example_barChart().Save(4, 4, "test.png"); err != nil {
		t.Error(err)
	}
}

func TestDrawEps(t *testing.T) {
	if err := Example_errbars().Save(4, 4, "test.eps"); err != nil {
		t.Error(err)
	}
}

func TestDrawSvg(t *testing.T) {
	if err := Example_points().Save(4, 4, "test.svg"); err != nil {
		t.Error(err)
	}
}

func TestDrawTiff(t *testing.T) {
	if err := Example_points().Save(4, 4, "test.tiff"); err != nil {
		t.Error(err)
	}
}

func TestDrawJpg(t *testing.T) {
	if err := Example_points().Save(4, 4, "test.jpg"); err != nil {
		t.Error(err)
	}
}

func TestDrawPdf(t *testing.T) {
	if err := Example_points().Save(4, 4, "test.pdf"); err != nil {
		t.Error(err)
	}
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

	quad := NewFunction(func(x float64) float64 { return x * x })
	quad.Color = color.RGBA{B: 255, A: 255}

	exp := NewFunction(func(x float64) float64 { return math.Pow(2, x) })
	exp.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	exp.Width = vg.Points(2)
	exp.Color = color.RGBA{G: 255, A: 255}

	sin := NewFunction(func(x float64) float64 { return 10*math.Sin(x) + 50 })
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
	uniform := make(Values, n)
	normal := make(Values, n)
	expon := make(Values, n)
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
	p.Y.Label.Text = "Values"

	// Make boxes for our data and add them to the plot.
	p.Add(NewBoxPlot(vg.Points(20), 0, uniform),
		NewBoxPlot(vg.Points(20), 1, normal),
		NewBoxPlot(vg.Points(20), 2, expon))

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
	p.Y.Label.Text = "Values"

	// Make boxes for our data and add them to the plot.
	uniBox := NewBoxPlot(vg.Points(20), 0, uniform)
	uniLabels, err := uniBox.OutsideLabels(uniform)
	if err != nil {
		panic(err)
	}

	normBox := NewBoxPlot(vg.Points(20), 1, normal)
	normLabels, err := normBox.OutsideLabels(normal)
	if err != nil {
		panic(err)
	}

	expBox := NewBoxPlot(vg.Points(20), 2, expon)
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
	p.X.Label.Text = "Values"

	// Make boxes for our data and add them to the plot.
	uniBox := HorizBoxPlot{NewBoxPlot(vg.Points(20), 0, uniform)}
	uniLabels, err := uniBox.OutsideLabels(uniform)
	if err != nil {
		panic(err)
	}

	normBox := HorizBoxPlot{NewBoxPlot(vg.Points(20), 1, normal)}
	normLabels, err := normBox.OutsideLabels(normal)
	if err != nil {
		panic(err)
	}

	expBox := HorizBoxPlot{NewBoxPlot(vg.Points(20), 2, expon)}
	expLabels, err := expBox.OutsideLabels(expon)
	if err != nil {
		panic(err)
	}
	p.Add(uniBox, uniLabels, normBox, normLabels, expBox, expLabels)

	// Add a GlyphBox plotter for debugging.
	p.Add(NewGlyphBoxes())

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
	p.Add(NewGrid())

	s := NewScatter(scatterData)
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	s.GlyphStyle.Radius = vg.Points(3)

	l := NewLine(lineData)
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	lpLine, lpPoints := NewLinePoints(linePointsData)
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
func randomPoints(n int) XYs {
	pts := make(XYs, n)
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

// Example_errbars draws points and error bars.
func Example_errbars() *plot.Plot {

	type errPoints struct {
		XYs
		YErrors
		XErrors
	}

	rand.Seed(int64(0))
	n := 15
	data := errPoints{
		XYs: randomPoints(n),
		YErrors: YErrors(randomError(n)),
		XErrors: XErrors(randomError(n)),
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	scatter := NewScatter(data)
	scatter.Shape = plot.CrossGlyph{}
	p.Add(scatter, NewXErrorBars(data), NewYErrorBars(data))
	p.Add(NewGlyphBoxes())

	return p
}

func randomError(n int) Errors {
	err := make(Errors, n)
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

	bs := NewBubbles(bubbleData, vg.Points(1), vg.Points(20))
	bs.Color = color.RGBA{R: 196, B: 128, A: 255}
	p.Add(bs)

	return p
}

// randomTriples returns some random x, y, z triples
// with some interesting kind of trend.
func randomTriples(n int) XYZs {
	data := make(XYZs, n)
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
	vals := make(Values, n)
	for i := 0; i < n; i++ {
		vals[i] = rand.NormFloat64()
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Histogram"
	h := NewHist(vals, 16)
	h.Normalize(1)
	p.Add(h)

	// The normal distribution function
	norm := NewFunction(stdNorm)
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
	groupA := Values{20, 35, 30, 35, 27}
	groupB := Values{25, 32, 34, 20, 25}
	groupC := Values{12, 28, 15, 21, 8}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Bar chart"
	p.Y.Label.Text = "Heights"

	w := vg.Points(15)

	barsA := NewBarChart(groupA, w)
	barsA.Color = color.RGBA{R: 255, A: 255}
	barsA.Offset = -w

	barsB := NewBarChart(groupB, w)
	barsB.Color = color.RGBA{R: 196, G: 196, A: 255}
	barsB.Offset = 0

	barsC := NewBarChart(groupC, w)
	barsC.Color = color.RGBA{B: 255, A: 255}
	barsC.Offset = w

	p.Add(barsA, barsB, barsC)
	p.Legend.Add("Group A", barsA)
	p.Legend.Add("Group B", barsB)
	p.Legend.Add("Group C", barsC)
	p.Legend.Top = true
	p.NominalX("One", "Two", "Three", "Four", "Five")

	return p
}

func TestEmpty(t *testing.T) {
	p, err := plot.New()
	if err != nil {
		t.Error(err)
	}

	if err := p.Save(4, 4, "empty.svg"); err != nil {
		t.Error(err)
	}
}
