// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/internal/cmpimg"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"math"
)

// ExampleScatter_color draws some scatter points.
// Each point is plotted with a different color depending on
// some external criteria.
func ExampleScatter_color() {
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
	scatterData := randomTriples(n)

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Points Example Color"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(NewGrid())

	sc, err := NewScatter(scatterData)
	if err != nil {
		log.Panic(err)
	}

	// Calculate the range of Z values.
	minZ, maxZ := math.Inf(1), math.Inf(-1)
	for _, xyz := range scatterData {
		if xyz.Z > maxZ {
			maxZ = xyz.Z
		}
		if xyz.Z < minZ {
			minZ = xyz.Z
		}
	}

	colors := moreland.Kindlmann() // Initialize a color map.
	colors.SetMax(255)
	colors.SetMin(0)
	max := colors.Max()
	min := colors.Min()

	// Specify style and color for individual points.
	sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		_, _, z := scatterData.XYZ(i)
		d := (z - minZ) / (maxZ - minZ)
		rng := max - min
		k := d*rng + min
		c, err := colors.At(k)
		if err != nil {
			log.Panic(err)
		}
		return draw.GlyphStyle{Color: c, Radius: vg.Points(3), Shape: draw.CircleGlyph{}}
	}
	p.Add(sc)

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
			val = min
		case len(thumbs) - 1:
			val = max
		}
		p.Legend.Add(fmt.Sprintf("%g", val), t)
	}

	// This is the width of the legend, experimentally determined.
	const legendWidth = 1.5 * vg.Centimeter

	// Slide the legend over so it doesn't overlap the ScatterPlot.
	p.Legend.XOffs = legendWidth

	img := vgimg.New(300, 230)
	dc := draw.New(img)
	dc = draw.Crop(dc, 0, -legendWidth, 0, 0) // Make space for the legend.
	p.Draw(dc)
	w, err := os.Create("testdata/scatterColor.png")
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

func TestScatterColor(t *testing.T) {
	cmpimg.CheckPlot(ExampleScatter_color, t, "scatterColor.png")
}
