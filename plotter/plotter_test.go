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

func TestDrawImage(t *testing.T) {
	if err := Example_horizontalBoxPlots().Save(4, 4, "test.png"); err != nil {
		t.Error(err)
	}
}

func TestDrawEps(t *testing.T) {
	if err := Example_boxPlots().Save(4, 4, "test.eps"); err != nil {
		t.Error(err)
	}
}

func TestDrawSvg(t *testing.T) {
	if err := Example_horizontalBoxPlots().Save(4, 4, "test.svg"); err != nil {
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
	uniform := make(ValueLabels, n)
	normal := make(ValueLabels, n)
	expon := make(ValueLabels, n)
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
	uniLabels, err := uniBox.PointLabels(uniform)
	if err != nil {
		panic(err)
	}

	normBox := NewBoxPlot(vg.Points(20), 1, normal)
	normLabels, err := normBox.PointLabels(normal)
	if err != nil {
		panic(err)
	}

	expBox := NewBoxPlot(vg.Points(20), 2, expon)
	expLabels, err := expBox.PointLabels(expon)
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

// Example_horizontalBoxPlots draws horizontal boxplots
// with some labels on their points.
func Example_horizontalBoxPlots() *plot.Plot {
	rand.Seed(int64(0))
	n := 100
	uniform := make(ValueLabels, n)
	normal := make(ValueLabels, n)
	expon := make(ValueLabels, n)
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
	uniLabels, err := uniBox.PointLabels(uniform)
	if err != nil {
		panic(err)
	}

	normBox := HorizBoxPlot{NewBoxPlot(vg.Points(20), 1, normal)}
	normLabels, err := normBox.PointLabels(normal)
	if err != nil {
		panic(err)
	}

	expBox := HorizBoxPlot{NewBoxPlot(vg.Points(20), 2, expon)}
	expLabels, err := expBox.PointLabels(expon)
	if err != nil {
		panic(err)
	}
	p.Add(uniBox, uniLabels, normBox, normLabels, expBox, expLabels)

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

	s := NewScatter(scatterData)
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	s.GlyphStyle.Radius = vg.Points(3)

	l := NewLine(lineData)
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	lp := NewLinePoints(linePointsData)
	lp.LineStyle.Color = color.RGBA{G: 255, A: 255}
	lp.GlyphStyle.Shape = plot.CircleGlyph
	lp.GlyphStyle.Color = color.RGBA{R: 255, A: 255}

	p.Add(s, l, lp)
	p.Legend.Add("scatter", s)
	p.Legend.Add("line", l)
	p.Legend.Add("line points", lp)

	return p
}

// randomPoints returns some random x, y points.
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

// An example of making a histogram.
func Example_histogram() *plot.Plot {
	rand.Seed(int64(0))
	n := 10000
	vals := make(XYs, n)
	for i := 0; i < n; i++ {
		vals[i].X = rand.NormFloat64()
		vals[i].Y = 1
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Histogram"
	h := NewHistogram(vals, 16)
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

func TestEmpty(t *testing.T) {
	p, err := plot.New()
	if err != nil {
		t.Error(err)
	}

	if err := p.Save(4, 4, "empty.svg"); err != nil {
		t.Error(err)
	}
}
