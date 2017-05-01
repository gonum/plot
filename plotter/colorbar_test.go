// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"log"
	"testing"

	"github.com/gonum/plot"
	"github.com/gonum/plot/internal/cmpimg"
	"github.com/gonum/plot/palette/moreland"
)

func ExampleColorBar_horizontal() {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	l := &ColorBar{ColorMap: moreland.ExtendedBlackBody()}
	l.ColorMap.SetMax(1)
	p.Add(l)
	p.HideY()
	p.X.Padding = 0
	p.Title.Text = "Title"

	if err = p.Save(300, 48, "testdata/colorBarHorizontal.png"); err != nil {
		log.Panic(err)
	}
}

func TestColorBar_horizontal(t *testing.T) {
	cmpimg.CheckPlot(ExampleColorBar_horizontal, t, "colorBarHorizontal.png")
}

func ExampleColorBar_vertical() {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	l := &ColorBar{ColorMap: moreland.ExtendedBlackBody()}
	l.ColorMap.SetMax(1)
	l.Vertical = true
	p.Add(l)
	p.HideX()
	p.Y.Padding = 0
	p.Title.Text = "Title"

	if err = p.Save(40, 300, "testdata/colorBarVertical.png"); err != nil {
		log.Panic(err)
	}
}

func TestColorBar_vertical(t *testing.T) {
	cmpimg.CheckPlot(ExampleColorBar_vertical, t, "colorBarVertical.png")
}
