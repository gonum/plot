// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package palette_test

import (
	"fmt"
	"log"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/internal/cmpimg"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/plotter"
)

// This example creates a color bar and a second color bar where the
// direction of the colors are reversed.
func ExampleReverse() {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	l := &plotter.ColorBar{ColorMap: moreland.Kindlmann()}
	l2 := &plotter.ColorBar{ColorMap: palette.Reverse(moreland.Kindlmann())}
	l.ColorMap.SetMin(0.5)
	l.ColorMap.SetMax(2.5)
	l2.ColorMap.SetMin(2.5)
	l2.ColorMap.SetMax(4.5)

	p.Add(l, l2)
	p.HideY()
	p.X.Padding = 0
	p.Title.Text = "A ColorMap and its Reverse"

	if err = p.Save(300, 48, "testdata/reverse.png"); err != nil {
		log.Panic(err)
	}
}

func TestReverse(t *testing.T) {
	cmpimg.CheckPlot(ExampleReverse, t, "reverse.png")
}

// This example creates a color palette from a reversed ColorMap.
func ExampleReverse_Palette() {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	thumbs := plotter.PaletteThumbnailers(palette.Reverse(moreland.Kindlmann()).Palette(10))
	for i, t := range thumbs {
		p.Legend.Add(fmt.Sprint(i), t)
	}
	p.HideAxes()
	p.X.Padding = 0
	p.Y.Padding = 0

	if err = p.Save(35, 120, "testdata/reverse_palette.png"); err != nil {
		log.Panic(err)
	}
}

func TestReverse_Palette(t *testing.T) {
	cmpimg.CheckPlot(ExampleReverse_Palette, t, "reverse_palette.png")
}
