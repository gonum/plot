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
	if err := Example().Save(4, 4, "test.png"); err != nil {
		t.Error(err)
	}
}

func TestDrawEps(t *testing.T) {
	if err := Example().Save(4, 4, "test.eps"); err != nil {
		t.Error(err)
	}
}

func TestDrawSvg(t *testing.T) {
	if err := Example_functions().Save(4, 4, "test.svg"); err != nil {
		t.Error(err)
	}
}

// An example of making and saving a plot.
func Example() *plot.Plot {
	rand.Seed(int64(0))
	pts := MakeXYLabelErrors(10)
	for i := range pts.XYs {
		if i == 0 {
			pts.XYs[i].X = rand.Float64()
		} else {
			pts.XYs[i].X = pts.XYs[i-1].X + rand.Float64()
		}
		pts.XYs[i].Y = rand.Float64()
		pts.Labels[i] = fmt.Sprintf("%d", i)
		pts.XErrors[i].Low = -rand.Float64() / 2
		pts.XErrors[i].High = rand.Float64() / 2
		pts.YErrors[i].Low = -rand.Float64() / 2
		pts.YErrors[i].High = rand.Float64() / 2
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Plot Title"
	p.X.Label.Text = "X Values"
	p.Y.Label.Text = "Y Values"
	line := Line{pts, DefaultLineStyle}
	scatter := Scatter{pts, DefaultGlyphStyle}
	errbars, err := MakeErrorBars(pts)
	if err != nil {
		panic(err)
	}
	labels, err := MakeLabels(pts)
	if err != nil {
		panic(err)
	}
	labels.XOffs = scatter.Radius
	labels.YOffs = scatter.Radius
	p.Add(line, scatter, errbars, labels)
	p.Legend.Add("line", line, scatter)
	p.Legend.Left = true
	return p
}

// An example of plotting a function.
func Example_functions() *plot.Plot {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Functions"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	quad := MakeFunction(func(x float64) float64 { return x * x })
	quad.Color = color.RGBA{B: 255, A: 255}

	exp := MakeFunction(func(x float64) float64 { return math.Pow(2, x) })
	exp.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	exp.Width = vg.Points(2)
	exp.Color = color.RGBA{G: 255, A: 255}

	sin := MakeFunction(func(x float64) float64 { return 10*math.Sin(x) + 50 })
	sin.Dashes = []vg.Length{vg.Points(4), vg.Points(5)}
	sin.Width = vg.Points(4)
	sin.Color = color.RGBA{R: 255, A: 255}

	p.Add(quad, exp, sin)
	p.Legend.Add("x^2", quad)
	p.Legend.Add("2^x", exp)
	p.Legend.Add("10*sin(x)+50", sin)

	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 100
	return p
}

// Draw the plotinum logo.
func Example_logo() *plot.Plot {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	DefaultLineStyle.Width = vg.Points(1)
	DefaultGlyphStyle.Radius = vg.Points(3)

	p.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{0, "0"}, {0.25, ""}, {0.5, "0.5"}, {0.75, ""}, {1, "1"},
	})
	p.X.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{0, "0"}, {0.25, ""}, {0.5, "0.5"}, {0.75, ""}, {1, "1"},
	})

	pts := XYs{{0, 0}, {0, 1}, {0.5, 1}, {0.5, 0.6}, {0, 0.6}}
	line := Line{pts, DefaultLineStyle}
	scatter := Scatter{pts, DefaultGlyphStyle}
	p.Add(line, scatter)

	pts = XYs{{1, 0}, {0.75, 0}, {0.75, 0.75}}
	line = Line{pts, DefaultLineStyle}
	scatter = Scatter{pts, DefaultGlyphStyle}
	p.Add(line, scatter)

	pts = XYs{{0.5, 0.5}, {1, 0.5}}
	line = Line{pts, DefaultLineStyle}
	scatter = Scatter{pts, DefaultGlyphStyle}
	p.Add(line, scatter)

	return p
}

// An example of making a box plot.
func Example_boxPlot() *plot.Plot {
	rand.Seed(int64(0))
	n := 10
	uniform := make(Ys, n)
	normal := make(Ys, n)
	expon := make(Ys, n)
	for i := 0; i < n; i++ {
		uniform[i] = rand.Float64()
		normal[i] = rand.NormFloat64()
		expon[i] = rand.ExpFloat64()
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Plot Title"
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

// An example of making a horizontal box plot.
func Example_horizontalBoxes() *plot.Plot {
	rand.Seed(int64(0))
	n := 10
	uniform := make(Ys, n)
	normal := make(Ys, n)
	expon := make(Ys, n)
	for i := 0; i < n; i++ {
		uniform[i] = rand.Float64()
		normal[i] = rand.NormFloat64()
		expon[i] = rand.ExpFloat64()
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Plot Title"
	p.X.Label.Text = "Values"

	// Make horizontal boxes for our data and add
	// them to the plot.
	p.Add(MakeHorizBoxPlot(vg.Points(20), 0, uniform),
		MakeHorizBoxPlot(vg.Points(20), 1, normal),
		MakeHorizBoxPlot(vg.Points(20), 2, expon))

	// Set the Y axis of the plot to nominal with
	// the given names for y=0, y=1 and y=2.
	p.NominalY("Uniform\nDistribution", "Normal\nDistribution",
		"Exponential\nDistribution")
	return p
}
