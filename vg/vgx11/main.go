// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build ignore

package main

import (
	"math/rand"
	"time"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgx11"
)

func main() {
	rand.Seed(0) // The default random seed is 1.

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	err = plotutil.AddLinePoints(
		p,
		"First", randomPoints(15),
	)
	if err != nil {
		panic(err)
	}

	cnvs, err := vgx11.New(4*96, 4*96, "Example")
	if err != nil {
		panic(err)
	}

	p.Draw(draw.New(cnvs))
	cnvs.Paint()
	time.Sleep(5 * time.Second)

	err = plotutil.AddLinePoints(
		p,
		"Second", randomPoints(15),
		"Third", randomPoints(15),
	)
	if err != nil {
		panic(err)
	}

	p.Draw(draw.New(cnvs))
	cnvs.Paint()
	time.Sleep(10 * time.Second)

	// Save the plot to a PNG file.
	//        if err := p.Save(4, 4, "points.png"); err != nil {
	//                panic(err)
	//        }
}

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}
