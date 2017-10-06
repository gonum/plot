// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"gonum.org/v1/plot/palette"
	"log"
	"math/rand"
	"testing"

	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/internal/cmpimg"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"os"
)

// ExampleScatter draws some scatter coloured points.
func ExampleScatterColor() {
	rnd := rand.New(rand.NewSource(1))

	// randomTriples returns some random x, y, z triples
	// with some interesting kind of trend.
	randomTriples := func(n int) XYZs {
		data := make(XYZs, n)
		for i := range data {
			if i == 0 {
				data[i].X = rnd.Float64()
			} else {
				data[i].X = data[i-1].X + 2*rnd.Float64()
			}
			data[i].Y = data[i].X + 10*rnd.Float64()
			data[i].Z = data[i].X
		}
		return data
	}

	n := 15
	scatterColorData := randomTriples(n)

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}

	p.Title.Text = "Scatter Colour Plot"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(NewGrid())

	pal := palette.Heat(12, 1)

	sc, err := NewScatterColor(scatterColorData, pal)
	if err != nil {
		log.Panic(err)
	}
	p.Add(sc)

	// Create a legend.
	thumbs := PaletteThumbnailers(pal)
	for i := len(thumbs) - 1; i >= 0; i-- {
		t := thumbs[i]
		if i != 0 && i != len(thumbs)-1 {
			p.Legend.Add("", t)
			continue
		}
		var val float64
		switch i {
		case 0:
			val = sc.Min
		case len(thumbs) - 1:
			val = sc.Max
		}
		p.Legend.Add(fmt.Sprintf("%.2g", val), t)
	}
	// This is the width of the legend, experimentally determined.
	const legendWidth = 1.25 * vg.Centimeter
	// Slide the legend over so it doesn't overlap the ScatterColorPlot.
	p.Legend.XOffs = legendWidth

	img := vgimg.New(300, 200)
	dc := draw.New(img)
	dc = draw.Crop(dc, 0, -legendWidth, 0, 0) // Make space for the legend.
	p.Draw(dc)
	w, err := os.Create("testdata/scattercolor.png")
	if err != nil {
		log.Panic(err)
	}
	png := vgimg.PngCanvas{Canvas: img}
	if _, err = png.WriteTo(w); err != nil {
		log.Panic(err)
	}
}

func TestScatterColor(t *testing.T) {
	cmpimg.CheckPlot(ExampleScatterColor, t, "scattercolor.png")
}
