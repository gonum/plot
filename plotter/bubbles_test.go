// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/internal/cmpimg"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// ExampleScatter_bubbles draws some scatter points.
// Each point is plotted with a different radius size depending on
// external criteria.
func ExampleScatter_bubbles() {
	rnd := rand.New(rand.NewSource(1))

	// randomTriples returns some random but correlated x, y, z triples
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

	n := 10
	scatterData := randomTriples(n)

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

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Bubbles"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	sc, err := NewScatter(scatterData)
	if err != nil {
		log.Panic(err)
	}

	// Specify style for individual points.
	sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		c := color.RGBA{R: 196, B: 128, A: 255}
		var minRadius, maxRadius = vg.Points(1), vg.Points(20)
		rng := maxRadius - minRadius
		_, _, z := scatterData.XYZ(i)
		d := (z - minZ) / (maxZ - minZ)
		r := vg.Length(d)*rng + minRadius
		return draw.GlyphStyle{Color: c, Radius: r, Shape: draw.CircleGlyph{}}
	}
	p.Add(sc)

	err = p.Save(200, 200, "testdata/bubbles.png")
	if err != nil {
		log.Panic(err)
	}
}

func TestNewBubbles(t *testing.T) {
	cmpimg.CheckPlot(ExampleScatter_bubbles, t, "bubbles.png")
}
