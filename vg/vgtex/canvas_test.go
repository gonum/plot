// Copyright 2016 The gonum Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgtex

import (
	"log"
	"os"
	"testing"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

func Example() {
	scatter, err := plotter.NewScatter(plotter.XYs{{1, 1}, {0, 1}, {0, 0}})
	if err != nil {
		log.Fatal(err)
	}
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.Add(scatter)
	// p.HideAxes()
	p.Title.Text = `A scatter plot: $\sqrt{\frac{e^{3i\pi}}{2\cos 3\pi}}$`
	p.X.Label.Text = `$x = \eta$`
	p.Y.Label.Text = `$y$ is some $\Phi$`

	c := NewDocument(5*vg.Centimeter, 5*vg.Centimeter)
	p.Draw(draw.New(c))
	c.FillString(p.Title.Font, 2.5*vg.Centimeter, 2.5*vg.Centimeter, "x")

	f, err := os.Create("testdata/scatter.tex")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err = c.WriteTo(f); err != nil {
		log.Fatal(err)
	}
}

func TestTexCanvas(t *testing.T) {
	Example()
}
