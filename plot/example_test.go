package plot

import (
	"code.google.com/p/plotinum/vg"
	"code.google.com/p/plotinum/vg/veceps"
	"math/rand"
)

// An example of making and saving a plot.
func Example() {
	// Get some data to plot.
	pts := make(XYs, 10)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = rand.Float64()
	}

	// Make an encapsulated postscript DrawArea
	w, h := vg.Inches(3), vg.Inches(2)
	da := NewDrawArea(veceps.New(w, h, "test"), w, h)

	// Make our plot and set some labels.
	p, err := New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Plot Title"
	p.X.Label.Text = "X Values"
	p.Y.Label.Text = "Y Values"
	line := Line{pts, DefaultLineStyle}
	scatter := Scatter{pts, DefaultGlyphStyle}
	p.Add(line, scatter)
	p.Legend.Add("line", line, scatter)
	p.Legend.Top = true
	p.Draw(da)

	// Save to test.eps
	err = da.Canvas.(*veceps.Canvas).Save("test.eps")
	if err != nil {
		panic(err)
	}
}

// An example of making a box plot.
func Example_boxPlot(da *DrawArea) {
	// Get some data to plot.
	n := 10
	uniform := make(Ys, n)
	normal := make(Ys, n)
	expon := make(Ys, n)
	for i := 0; i < n; i++ {
		uniform[i] = rand.Float64()
		normal[i] = rand.NormFloat64()
		expon[i] = rand.ExpFloat64()
	}

	// Make our plot and set some labels.
	p, err := New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Plot Title"
	p.Y.Label.Text = "Values"

	// Make boxes for our data and add them to the plot.
	p.Add(NewBox(vg.Points(20), 0, uniform),
		NewBox(vg.Points(20), 1, normal),
		NewBox(vg.Points(20), 2, expon))

	// Set the X axis of the plot to nominal with
	// the given names for x=0, x=1 and x=2.
	p.NominalX("Uniform\nDistribution", "Normal\nDistribution",
		"Exponential\nDistribution")

	p.Draw(da)
}

// An example of making a horizontal box plot.
func Example_horizontalBoxes(da *DrawArea) {
	// Get some data to plot.
	n := 10
	uniform := make(Ys, n)
	normal := make(Ys, n)
	expon := make(Ys, n)
	for i := 0; i < n; i++ {
		uniform[i] = rand.Float64()
		normal[i] = rand.NormFloat64()
		expon[i] = rand.ExpFloat64()
	}

	// Make our plot and set some labels.
	p, err := New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Plot Title"
	p.X.Label.Text = "Values"

	// Make horizontal boxes for our data and add
	// them to the plot.
	p.Add(MakeHorizBox(vg.Points(20), 0, uniform),
		MakeHorizBox(vg.Points(20), 1, normal),
		MakeHorizBox(vg.Points(20), 2, expon))

	// Set the Y axis of the plot to nominal with
	// the given names for y=0, y=1 and y=2.
	p.NominalY("Uniform\nDistribution", "Normal\nDistribution",
		"Exponential\nDistribution")

	p.Draw(da)
}
