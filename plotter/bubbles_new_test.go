// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"image/color"
	"log"
	"math/rand"
	"os"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/internal/cmpimg"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

// ExampleScatter_bubbles draws some scatter points.
// Each point is plotted with a different radius size depending on
// some external criteria.
func ExampleScatter_bubbles() {
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
	scatterDataNew := randomPoints(n)

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Bubble Plot"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	sc, err := NewScatter(scatterDataNew)
	if err != nil {
		log.Panic(err)
	}

	//Specify style for individual points.
	sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		c := color.RGBA{R: 196, B: 128, A: 255}
		r := float64(i)
		return draw.GlyphStyle{Color: c, Radius: vg.Points(r), Shape: draw.CircleGlyph{}}
	}
	p.Add(sc)

	img := vgimg.New(200, 200)
	dc := draw.New(img)
	p.Draw(dc)
	w, err := os.Create("testdata/bubblesNew.png")
	defer w.Close()
	if err != nil {
		log.Panic(err)
	}
	png := vgimg.PngCanvas{Canvas: img}
	if _, err = png.WriteTo(w); err != nil {
		log.Panic(err)
	}
	if err = w.Close(); err != nil {
		log.Panic(err)
	}
}

func TestNewBubbles(t *testing.T) {
	cmpimg.CheckPlot(ExampleScatter_bubbles, t, "bubblesNew.png")
}
