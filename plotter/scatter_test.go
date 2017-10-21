// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"image/color"
	"log"
	"math/rand"
	"testing"

	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/internal/cmpimg"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"os"
)

// ExampleScatter draws some scatter points, a line,
// and a line with points.
func ExampleScatter() {
	rnd := rand.New(rand.NewSource(1))

	// randomPoints returns some random x, y points
	// with some interesting kind of trend.
	randomPoints := func(n int) XYs {
		pts := make(XYs, n)
		for i := range pts {
			if i == 0 {
				pts[i].X = rnd.Float64()
			} else {
				pts[i].X = pts[i-1].X + rnd.Float64()
			}
			pts[i].Y = pts[i].X + 10*rnd.Float64()
		}
		return pts
	}

	n := 15
	scatterData := randomPoints(n)
	lineData := randomPoints(n)
	linePointsData := randomPoints(n)
	scatterDataNew := randomPoints(n)

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Points Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(NewGrid())

	s, err := NewScatter(scatterData)
	if err != nil {
		log.Panic(err)
	}

	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	s.GlyphStyle.Radius = vg.Points(3)

	l, err := NewLine(lineData)
	if err != nil {
		log.Panic(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	lpLine, lpPoints, err := NewLinePoints(linePointsData)
	if err != nil {
		log.Panic(err)
	}
	lpLine.Color = color.RGBA{G: 255, A: 255}
	lpPoints.Shape = draw.CircleGlyph{}
	lpPoints.Color = color.RGBA{R: 255, A: 255}

	sc, err := NewScatter(scatterDataNew)
	if err != nil {
		log.Panic(err)
	}

	colors := moreland.Kindlmann() // Initialize a color map.
	colors.SetMax(255)
	colors.SetMin(0)

	z := []float64{31, 41, 51, 61, 71, 81, 91, 101, 111, 121, 131, 141, 151, 161, 171, 181}

	sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		c, err := colors.At(z[i])
		if err != nil {
			log.Panic(err)
		}
		return draw.GlyphStyle{Color: c, Radius: vg.Points(3), Shape: draw.CircleGlyph{}}
	}

	p.Add(s, l, lpLine, lpPoints, sc)
	p.Legend.Add("scatter", s)
	p.Legend.Add("line", l)
	p.Legend.Add("line points", lpLine, lpPoints)
	p.Legend.Add("scatterColor", sc)

	// Create a legend.
	thumbs := PaletteThumbnailers(colors.Palette(n))
	for i := len(thumbs) - 1; i >= 0; i-- {
		t := thumbs[i]
		if i != 0 && i != len(thumbs)-1 {
			p.Legend.Add("", t)
			continue
		}
		var val float64
		switch i {
		case 0:
			val = z[0]
		case len(thumbs) - 1:
			val = z[n]
		}
		p.Legend.Add(fmt.Sprintf("%g", val), t)
	}

	// This is the width of the legend, experimentally determined.
	const legendWidth = 3 * vg.Centimeter
	// Slide the legend over so it doesn't overlap the ScatterPlot.
	p.Legend.XOffs = legendWidth

	img := vgimg.New(350, 290)
	dc := draw.New(img)
	dc = draw.Crop(dc, 0, -legendWidth, 0, 0) // Make space for the legend.
	p.Draw(dc)
	w, err := os.Create("testdata/scatter.png")
	if err != nil {
		log.Panic(err)
	}
	png := vgimg.PngCanvas{Canvas: img}
	if _, err = png.WriteTo(w); err != nil {
		log.Panic(err)
	}
}

func TestScatter(t *testing.T) {
	cmpimg.CheckPlot(ExampleScatter, t, "scatter.png")
}
