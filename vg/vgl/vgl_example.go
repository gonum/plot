// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build ignore

package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgl"
)

func init() {
	runtime.LockOSThread()
}

func main() {

	err := glfw.Init()
	if err != nil {
		panic(fmt.Errorf("glfw init failed: %v\n", err))
	}
	defer glfw.Terminate()

	err = vgl.Run(run)
	if err != nil {
		panic(err)
	}
}

func run() error {
	// Draw some random values from the standard
	// normal distribution.
	rand.Seed(int64(0))
	v := make(plotter.Values, 10000)
	for i := range v {
		v[i] = rand.NormFloat64()
	}

	// Make a plot and set its title.
	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = "Histogram"
	p.X.Label.Text = "X-axis"
	p.Y.Label.Text = "Y-axis"

	// Draw a grid behind the data
	p.Add(plotter.NewGrid())
	p.Add(plotter.NewGlyphBoxes())

	// Create a histogram of our values drawn
	// from the standard normal.
	h, err := plotter.NewHist(v, 16)
	if err != nil {
		return err
	}
	// Normalize the area under the histogram to
	// sum to one.
	h.Normalize(1)
	p.Add(h)

	// The normal distribution function
	norm := plotter.NewFunction(stdNorm)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)

	scatter, err := plotter.NewScatter(plotter.XYs{{1, 1}, {0, 1}, {0, 0}})
	if err != nil {
		return err
	}
	p.Add(scatter)

	//cnvs, err := vgl.New(4*96, 4*96, "Example")
	cnvs, err := vgl.New(800, 600, "Example")
	if err != nil {
		return err
	}
	p.Draw(draw.New(cnvs))
	cnvs.Paint()
	return err
}

// stdNorm returns the probability of drawing a
// value from a standard normal distribution.
func stdNorm(x float64) float64 {
	const sigma = 1.0
	const mu = 0.0
	const root2π = 2.50662827459517818309
	return 1.0 / (sigma * root2π) * math.Exp(-((x-mu)*(x-mu))/(2*sigma*sigma))
}
