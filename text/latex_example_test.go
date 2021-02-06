// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text_test

import (
	"image/color"
	"log"
	"math"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/font/liberation"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func ExampleLatex() {
	fonts := font.NewCache(liberation.Collection())
	plot.DefaultTextHandler = text.Latex{
		Fonts: fonts,
	}

	p := plot.New()
	p.Title.Text = `$f(x) = \sqrt{\alpha x \Gamma}$`
	p.X.Label.Text = `$\frac{\sqrt{x}}{2\pi\Gamma\gamma}$`
	p.Y.Label.Text = `$\beta$`

	p.X.Min = -1
	p.X.Max = +1
	p.Y.Min = -1
	p.Y.Max = +1

	labels, err := plotter.NewLabels(plotter.XYLabels{
		XYs: []plotter.XY{
			{X: +0.0, Y: +0.0},
			{X: -0.9, Y: -0.9},
			{X: +0.6, Y: +0.0},
			{X: +0.5, Y: -0.9},
		},
		Labels: []string{
			`$\frac{\sqrt{x}}{2\pi\Gamma\gamma}$`,
			`$LaTeX$`,
			"plain",
			`$\frac{\sqrt{x}}{2\beta}$`,
		},
	})
	if err != nil {
		log.Fatalf("could not create labels: %+v", err)
	}
	labels.TextStyle[0].Font.Size = 24
	labels.TextStyle[0].Color = color.RGBA{B: 255, A: 255}
	labels.TextStyle[0].XAlign = draw.XCenter
	labels.TextStyle[0].YAlign = draw.YCenter

	labels.TextStyle[1].Font.Size = 24
	labels.TextStyle[1].Color = color.RGBA{R: 255, A: 255}
	labels.TextStyle[1].Rotation = math.Pi / 4

	labels.TextStyle[2].Font.Size = 24
	labels.TextStyle[2].Rotation = math.Pi / 4
	labels.TextStyle[2].YAlign = draw.YCenter
	labels.TextStyle[2].Handler = &text.Plain{Fonts: fonts}

	labels.TextStyle[3].Font.Size = 24
	labels.TextStyle[3].Rotation = math.Pi / 2

	p.Add(labels)
	p.Add(plotter.NewGlyphBoxes())
	p.Add(plotter.NewGrid())

	err = p.Save(10*vg.Centimeter, 5*vg.Centimeter, "testdata/latex.png")
	if err != nil {
		log.Fatalf("could not save plot: %+v", err)
	}
}

func TestLatex(t *testing.T) {
	cmpimg.CheckPlot(ExampleLatex, t, "latex.png")
}
