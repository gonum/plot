// +build ignore

package main

import (
	"flag"
	"image/color"
	"math"
	"math/rand"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

var (
	style = flag.String("style", "gnuplot", "plot style (default|gnuplot)")
)

func main() {
	flag.Parse()

	// Draw some random values from the standard
	// normal distribution.
	rand.Seed(int64(0))
	v := make(plotter.Values, 10000)
	for i := range v {
		v[i] = rand.NormFloat64()
	}

	// Make a plot and set its title.
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Histogram"
	p.X.Label.Text = "X-axis"
	p.Y.Label.Text = "Y-axis"

	if *style == "gnuplot" {
		p.Style = GnuplotStyle{}
		p.X.Padding = 0
		p.Y.Padding = 0
	}

	// Draw a grid behind the data
	p.Add(plotter.NewGrid())
	p.Add(plotter.NewGlyphBoxes())

	// Create a histogram of our values drawn
	// from the standard normal.
	h, err := plotter.NewHist(v, 16)
	if err != nil {
		panic(err)
	}
	// Normalize the area under the histogram to
	// sum to one.
	h.Normalize(1)
	p.Add(h)

	// The normal distribution function
	norm := plotter.NewFunction(stdNorm)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)

	// Save the plot to a PNG file.
	if err := p.Save(4, 4, "hist.png"); err != nil {
		panic(err)
	}
}

// stdNorm returns the probability of drawing a
// value from a standard normal distribution.
func stdNorm(x float64) float64 {
	const sigma = 1.0
	const mu = 0.0
	const root2π = 2.50662827459517818309
	return 1.0 / (sigma * root2π) * math.Exp(-((x-mu)*(x-mu))/(2*sigma*sigma))
}

type GnuplotStyle struct{}

func (s GnuplotStyle) DrawPlot(p *plot.Plot, c draw.Canvas) {
	if p.BackgroundColor != nil {
		c.SetColor(p.BackgroundColor)
		c.Fill(c.Rectangle.Path())
	}
	if p.Title.Text != "" {
		cx := p.DataCanvas(c)
		c.FillText(p.Title.TextStyle, cx.Center().X, c.Max.Y, -0.5, -1, p.Title.Text)
		c.Max.Y -= p.Title.Height(p.Title.Text) - p.Title.Font.Extents().Descent
		c.Max.Y -= p.Title.Padding
	}

	p.X.SanitizeRange()
	x := plot.HorizontalAxis{p.X}
	p.Y.SanitizeRange()
	y := plot.VerticalAxis{p.Y}

	ywidth := y.Size()
	xheight := x.Size()

	xda := plot.PadX(p, draw.Crop(c, ywidth-y.Width-y.Padding, 0, 0, 0))
	yda := plot.PadY(p, draw.Crop(c, 0, xheight-x.Width-x.Padding, 0, 0))

	x.Draw(xda)
	y.Draw(yda)
	xmin := xda.Min.X
	xmax := xda.Max.X
	ymin := yda.Min.Y
	ymax := xda.Max.Y
	xda.StrokeLine2(x.LineStyle, xmin, ymax, xmax, ymax)
	xda.StrokeLine2(x.LineStyle, xmin, ymin, xmax, ymin)
	yda.StrokeLine2(y.LineStyle, xmin, ymin, xmin, ymax)
	yda.StrokeLine2(y.LineStyle, xmax, ymin, xmax, ymax)

	datac := plot.PadY(p, plot.PadX(p, draw.Crop(c, ywidth, xheight, 0, 0)))
	for _, data := range p.Plotters() {
		data.Plot(datac, p)
	}

	p.Legend.Draw(draw.Crop(draw.Crop(c, ywidth, 0, 0, 0), 0, 0, xheight, 0))
}
