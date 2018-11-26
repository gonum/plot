// Copyright ©2018 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"image/color"
	"log"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/internal/cmpimg"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// ExampleStep draws some filled step lines.
func ExampleStep() {
	rnd := rand.New(rand.NewSource(1))

	// randomPoints returns some random x, y points
	// with some interesting kind of trend.
	randomPoints := func(n int, x float64) XYs {
		pts := make(XYs, n)
		for i := range pts {
			if i == 0 {
				pts[i].X = x + rnd.Float64()
			} else {
				pts[i].X = pts[i-1].X + 0.5 + rnd.Float64()
			}
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
	p.Add(NewGrid())

	stepPre, err := NewStep(randomPoints(n, 0.))
	if err != nil {
		log.Panic(err)
	}
	stepPre.FillColor = color.RGBA{A: 40}

	stepMid, err := NewStep(randomPoints(n, 3.5))
	if err != nil {
		log.Panic(err)
	}
	stepMid.StepStyle = MidStep
	stepMid.LineStyle = &draw.LineStyle{Color: color.RGBA{R: 196, B: 128, A: 255}, Width: vg.Points(1)}

	stepPost, err := NewStep(randomPoints(n, 7.))
	if err != nil {
		log.Panic(err)
	}
	stepPost.StepStyle = PostStep
	stepPost.LineStyle = nil
	stepPost.FillColor = color.RGBA{B: 255, A: 255}

	p.Add(stepPre, stepMid, stepPost)
	p.Legend.Add("pre", stepPre)
	p.Legend.Add("mid", stepMid)
	p.Legend.Add("post", stepPost)
	p.Legend.Top = true

	err = p.Save(200, 200, "testdata/step.png")
	if err != nil {
		log.Panic(err)
	}
}

func TestStep(t *testing.T) {
	cmpimg.CheckPlot(ExampleStep, t, "step.png")
}
