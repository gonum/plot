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
	scatter, err := plotter.NewScatter(plotter.XYs{{X: 1, Y: 1}, {X: 0, Y: 1}, {X: 0, Y: 0}})
	if err != nil {
		log.Fatal(err)
	}
	scatter.Color = color.RGBA{R: 255, A: 200}
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.Add(scatter)
	// p.HideAxes()
	p.Title.Text = `A scatter plot: $\sqrt{\frac{e^{3i\pi}}{2\cos 3\pi}}$`
	p.Title.TextStyle.Font.Size = 16
	p.X.Label.Text = `$x = \eta$`
	p.Y.Label.Text = `$y$ is some $\Phi$`

	txtFont := p.X.Label.TextStyle.Font

	c := vgtex.NewDocument(5*vg.Centimeter, 5*vg.Centimeter)
	p.Draw(draw.New(c))
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
