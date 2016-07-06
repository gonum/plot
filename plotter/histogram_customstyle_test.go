// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// An example of making a histogram with a custom style.
func ExampleHistogramCustomStyle() {
	// stdNorm returns the probability of drawing a
	// value from a standard normal distribution.
	stdNorm := func(x float64) float64 {
		const sigma = 1.0
		const mu = 0.0
		const root2π = 2.50662827459517818309
		return 1.0 / (sigma * root2π) * math.Exp(-((x-mu)*(x-mu))/(2*sigma*sigma))
	}

	n := 10000
	vals := make(Values, n)
	for i := 0; i < n; i++ {
		vals[i] = rand.NormFloat64()
	}

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Histogram"
	p.X.Label.Text = "X-axis"
	p.Y.Label.Text = "Y-axis"
	p.Style = CustomStyle{}

	h, err := NewHist(vals, 16)
	if err != nil {
		log.Panic(err)
	}
	h.Normalize(1)
	p.Add(h)

	// The normal distribution function
	norm := NewFunction(stdNorm)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)

	err = p.Save(200, 200, "testdata/histogram-custom-style.png")
	if err != nil {
		log.Panic(err)
	}
}

func TestHistogramCustomStyle(t *testing.T) {
	checkPlot(ExampleHistogramCustomStyle, t, "histogram-custom-style.png")
}

type CustomStyle struct{}

func (CustomStyle) DrawPlot(p *plot.Plot, c draw.Canvas) {
	if p.BackgroundColor != nil {
		c.SetColor(p.BackgroundColor)
		c.Fill(c.Rectangle.Path())
	}

	xpad := p.X.Padding
	ypad := p.Y.Padding
	defer func() {
		p.X.Padding = xpad
		p.Y.Padding = ypad
	}()

	p.X.Padding = 0
	p.Y.Padding = 0

	p.X.SanitizeRange()
	p.Y.SanitizeRange()
	xaxis := plot.HorizontalAxis{p.X}
	yaxis := plot.VerticalAxis{p.Y}

	ywidth := yaxis.Size()
	xheight := xaxis.Size()

	if p.Title.Text != "" {
		cx := draw.Crop(c, ywidth, 0, 0, 0)
		c.FillText(p.Title.TextStyle, vg.Point{cx.Center().X, c.Max.Y}, -0.5, -1, p.Title.Text)
		c.Max.Y -= p.Title.Height(p.Title.Text) - p.Title.Font.Extents().Descent
		c.Max.Y -= p.Title.Padding
	}

	xc := plot.PadX(p, draw.Crop(c, ywidth-yaxis.Width-yaxis.Padding, 0, 0, 0))
	yc := plot.PadY(p, draw.Crop(c, 0, xheight-xaxis.Width-xaxis.Padding, xheight, 0))

	xaxis.Draw(xc)
	yaxis.Draw(yc)
	xmin := xc.Min.X
	xmax := xc.Max.X
	ymin := yc.Min.Y
	ymax := xc.Max.Y
	xc.StrokeLine2(xaxis.LineStyle, xmin, ymax, xmax, ymax)
	xc.StrokeLine2(xaxis.LineStyle, xmin, ymin, xmax, ymin)
	yc.StrokeLine2(yaxis.LineStyle, xmin, ymin, xmin, ymax)
	yc.StrokeLine2(yaxis.LineStyle, xmax, ymin, xmax, ymax)

	datac := plot.PadY(p, plot.PadX(p, draw.Crop(c, ywidth, 0, xheight, 0)))
	for _, data := range p.Plotters() {
		data.Plot(datac, p)
	}

	p.Legend.Draw(draw.Crop(draw.Crop(c, ywidth, 0, 0, 0), 0, 0, xheight, 0))
}
