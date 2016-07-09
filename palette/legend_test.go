// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package palette

import (
	"log"
	"testing"

	"github.com/gonum/plot"
	"github.com/gonum/plot/internal/cmpimg"
	"github.com/gonum/plot/palette/moreland"
)

func ExampleColorMapLegend_horizontal() {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	cm := moreland.ExtendedBlackBody()
	cm.SetMax(1)
	l := NewColorMapLegend(cm)
	p.Add(l)
	l.SetupPlot(p)

	if err = p.Save(300, 40, "testdata/colorMapLegendHorizontal.png"); err != nil {
		log.Panic(err)
	}
}

func TestColorMapLegend_horizontal(t *testing.T) {
	cmpimg.CheckPlot(ExampleColorMapLegend_horizontal, t, "colorMapLegendHorizontal.png")
}

func ExampleColorMapLegend_vertical() {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	cm := moreland.ExtendedBlackBody()
	cm.SetMax(1)
	l := NewColorMapLegend(cm)
	l.Vertical = true
	p.Add(l)
	l.SetupPlot(p)

	if err = p.Save(40, 300, "testdata/colorMapLegendVertical.png"); err != nil {
		log.Panic(err)
	}
}

func TestColorMapLegend_vertical(t *testing.T) {
	cmpimg.CheckPlot(ExampleColorMapLegend_vertical, t, "colorMapLegendVertical.png")
}
