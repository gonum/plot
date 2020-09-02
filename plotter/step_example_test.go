// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"image/color"
	"log"

	"golang.org/x/exp/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func ExampleLine_stepLine() {
	rnd := rand.New(rand.NewSource(1))

	// randomPoints returns some random x, y points
	// with some interesting kind of trend.
	randomPoints := func(n int, x float64) plotter.XYs {
		pts := make(plotter.XYs, n)
		for i := range pts {
			pts[i].X = float64(i) + x
			pts[i].Y = 5. + 10*rnd.Float64()
		}
		return pts
	}

	n := 4

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Step Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	stepPre, err := plotter.NewLine(randomPoints(n, 0))
	if err != nil {
		log.Panic(err)
	}
	stepPre.StepStyle = plotter.PreStep
	stepPre.FillColor = color.RGBA{R: 196, G: 255, B: 196, A: 255}

	stepMid, err := plotter.NewLine(randomPoints(n, 3.5))
	if err != nil {
		log.Panic(err)
	}
	stepMid.StepStyle = plotter.MidStep
	stepMid.LineStyle = draw.LineStyle{Color: color.RGBA{R: 196, B: 128, A: 255}, Width: vg.Points(1)}

	stepMidFilled, err := plotter.NewLine(randomPoints(n, 7))
	if err != nil {
		log.Panic(err)
	}
	stepMidFilled.StepStyle = plotter.MidStep
	stepMidFilled.LineStyle = draw.LineStyle{Color: color.RGBA{R: 196, B: 128, A: 255}, Width: vg.Points(1)}
	stepMidFilled.FillColor = color.RGBA{R: 255, G: 196, B: 196, A: 255}

	stepPost, err := plotter.NewLine(randomPoints(n, 10.5))
	if err != nil {
		log.Panic(err)
	}
	stepPost.StepStyle = plotter.PostStep
	stepPost.LineStyle.Width = 0
	stepPost.FillColor = color.RGBA{R: 196, G: 196, B: 255, A: 255}

	p.Add(stepPre, stepMid, stepMidFilled, stepPost)
	p.Legend.Add("pre", stepPre)
	p.Legend.Add("mid", stepMid)
	p.Legend.Add("midFilled", stepMidFilled)
	p.Legend.Add("post", stepPost)
	p.Legend.Top = true
	p.Y.Max = 20
	p.Y.Min = 0

	err = p.Save(200, 200, "testdata/step.png")
	if err != nil {
		log.Panic(err)
	}
}
