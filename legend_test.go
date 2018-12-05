// Copyright ©2018 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"image/color"
	"os"
	"testing"

	"gonum.org/v1/plot/internal/cmpimg"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type exampleThumbnailer struct {
	color.Color
}

// Thumbnail fulfills the plot.Thumbnailer interface.
func (et exampleThumbnailer) Thumbnail(c *draw.Canvas) {
	pts := []vg.Point{
		{c.Min.X, c.Min.Y},
		{c.Min.X, c.Max.Y},
		{c.Max.X, c.Max.Y},
		{c.Max.X, c.Min.Y},
	}
	poly := c.ClipPolygonY(pts)
	c.FillPolygon(et.Color, poly)

	pts = append(pts, vg.Point{X: c.Min.X, Y: c.Min.Y})
	outline := c.ClipLinesY(pts)
	c.StrokeLines(draw.LineStyle{
		Color: color.Black,
		Width: vg.Points(1),
	}, outline...)
}

// This example creates a some standalone legends with borders around them.
func ExampleLegend_standalone() {
	c := vgimg.New(vg.Points(120), vg.Points(100))
	dc := draw.New(c)

	// These example thumbnailers could be replaced with any of Plotters
	// in the plotter subpackage.
	red := exampleThumbnailer{Color: color.NRGBA{R: 255, A: 255}}
	green := exampleThumbnailer{Color: color.NRGBA{G: 255, A: 255}}
	blue := exampleThumbnailer{Color: color.NRGBA{B: 255, A: 255}}

	l, err := NewLegend()
	if err != nil {
		panic(err)
	}
	l.Add("red", red)
	l.Add("green", green)
	l.Add("blue", blue)
	l.Padding = vg.Millimeter

	// purpleRectangle draws a purple rectangle around the given Legend.
	purpleRectangle := func(l Legend) {
		r := l.Rectangle(dc)
		dc.StrokeLines(draw.LineStyle{
			Color: color.NRGBA{R: 255, B: 255, A: 255},
			Width: vg.Points(1),
		}, []vg.Point{
			{r.Min.X, r.Min.Y}, {r.Min.X, r.Max.Y}, {r.Max.X, r.Max.Y},
			{r.Max.X, r.Min.Y}, {r.Min.X, r.Min.Y},
		})
	}

	l.Draw(dc)
	purpleRectangle(l)

	l.Left = true
	l.Draw(dc)
	purpleRectangle(l)

	l.Top = true
	l.Draw(dc)
	purpleRectangle(l)

	l.Left = false
	l.Draw(dc)
	purpleRectangle(l)

	w, err := os.Create("testdata/legend_standalone.png")
	if err != nil {
		panic(err)
	}

	png := vgimg.PngCanvas{Canvas: c}
	if _, err := png.WriteTo(w); err != nil {
		panic(err)
	}
}

func TestLegend_standalone(t *testing.T) {
	cmpimg.CheckPlot(ExampleLegend_standalone, t, "legend_standalone.png")
}
