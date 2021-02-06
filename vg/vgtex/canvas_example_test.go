// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgtex_test

import (
	"image/color"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgtex"
)

// An example of making a LaTeX plot.
func Example() {
	p := plot.New()
	// p.HideAxes()
	p.Title.Text = `A scatter plot: $\sqrt{\frac{e^{3i\pi}}{2\cos 3\pi}}$`
	p.Title.TextStyle.Font.Size = 16
	p.X.Label.Text = `$x = \eta$`
	p.Y.Label.Text = `$y$ is some $\Phi$`

	scatter1, err := plotter.NewScatter(plotter.XYs{{X: 1, Y: 1}, {X: 0, Y: 1}, {X: 0, Y: 0}})
	if err != nil {
		log.Fatal(err)
	}
	scatter1.Color = color.RGBA{R: 255, A: 200}

	scatter2, err := plotter.NewScatter(plotter.XYs{{X: 1, Y: 0}, {X: 1, Y: 0.5}})
	if err != nil {
		log.Fatal(err)
	}
	scatter2.GlyphStyle.Shape = draw.PyramidGlyph{}
	scatter2.GlyphStyle.Radius = 2
	scatter2.Color = color.RGBA{B: 255, A: 200}

	p.Add(scatter1, scatter2)

	txtFont := p.TextHandler.Cache().Lookup(
		p.X.Label.TextStyle.Font,
		p.X.Label.TextStyle.Font.Size,
	)

	c := vgtex.NewDocument(5*vg.Centimeter, 5*vg.Centimeter)
	p.Draw(draw.New(c))

	c.SetColor(color.Black)
	c.FillString(txtFont, vg.Point{X: 2.5 * vg.Centimeter, Y: 2.5 * vg.Centimeter}, "x")

	f, err := os.Create("testdata/scatter.tex")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err = c.WriteTo(f); err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}
