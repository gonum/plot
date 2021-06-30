// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"image/color"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func TestLabels(t *testing.T) {
	cmpimg.CheckPlot(ExampleLabels, t, "labels.png")
	cmpimg.CheckPlot(ExampleLabels_inCanvasCoordinates, t, "labels_cnv_coords.png")
}

// TestLabelsWithGlyphBoxes tests the position of the glyphbox around
// a block of text, checking whether we correctly take into account
// the descent+ascent of a glyph.
func TestLabelsWithGlyphBoxes(t *testing.T) {
	cmpimg.CheckPlot(
		func() {
			const fontSize = 24

			p := plot.New()
			p.Title.Text = "Labels"
			p.X.Min = -1
			p.X.Max = +1
			p.Y.Min = -1
			p.Y.Max = +1

			const (
				left   = 0.00
				middle = 0.02
				right  = 0.04
			)

			labels, err := plotter.NewLabels(plotter.XYLabels{
				XYs: []plotter.XY{
					{X: -0.8 + left, Y: -0.5},   // Aq + y-align bottom
					{X: -0.6 + middle, Y: -0.5}, // Aq + y-align center
					{X: -0.4 + right, Y: -0.5},  // Aq + y-align top

					{X: -0.8 + left, Y: +0.5}, // ditto for Aq\nAq
					{X: -0.6 + middle, Y: +0.5},
					{X: -0.4 + right, Y: +0.5},

					{X: +0.0 + left, Y: +0}, // ditto for Bg\nBg\nBg
					{X: +0.2 + middle, Y: +0},
					{X: +0.4 + right, Y: +0},
				},
				Labels: []string{
					"Aq", "Aq", "Aq",
					"Aq\nAq", "Aq\nAq", "Aq\nAq",

					"Bg\nBg\nBg",
					"Bg\nBg\nBg",
					"Bg\nBg\nBg",
				},
			})
			if err != nil {
				t.Fatalf("could not creates labels plotter: %+v", err)
			}
			for i := range labels.TextStyle {
				sty := &labels.TextStyle[i]
				sty.Font.Size = vg.Length(fontSize)
			}
			labels.TextStyle[0].YAlign = draw.YBottom
			labels.TextStyle[1].YAlign = draw.YCenter
			labels.TextStyle[2].YAlign = draw.YTop

			labels.TextStyle[3].YAlign = draw.YBottom
			labels.TextStyle[4].YAlign = draw.YCenter
			labels.TextStyle[5].YAlign = draw.YTop

			labels.TextStyle[6].YAlign = draw.YBottom
			labels.TextStyle[7].YAlign = draw.YCenter
			labels.TextStyle[8].YAlign = draw.YTop

			lred, err := plotter.NewLabels(plotter.XYLabels{
				XYs: []plotter.XY{
					{X: -0.8 + left, Y: +0.5},
					{X: +0.0 + left, Y: +0},
				},
				Labels: []string{
					"Aq", "Bg",
				},
			})
			if err != nil {
				t.Fatalf("could not creates labels plotter: %+v", err)
			}
			for i := range lred.TextStyle {
				sty := &lred.TextStyle[i]
				sty.Font.Size = vg.Length(fontSize)
				sty.Color = color.RGBA{R: 255, A: 255}
				sty.YAlign = draw.YBottom
			}

			// label with positive x/y-offsets
			loffp, err := plotter.NewLabels(plotter.XYLabels{
				XYs:    []plotter.XY{{X: left}},
				Labels: []string{"Bg"},
			})
			if err != nil {
				t.Fatalf("could not creates labels plotter: %+v", err)
			}
			for i := range loffp.TextStyle {
				sty := &loffp.TextStyle[i]
				sty.Font.Size = vg.Length(fontSize)
				sty.Color = color.RGBA{G: 255, A: 255}
			}
			loffp.Offset = vg.Point{X: 75, Y: 75}

			// label with negative x/y-offsets
			loffm, err := plotter.NewLabels(plotter.XYLabels{
				XYs:    []plotter.XY{{X: left}},
				Labels: []string{"Bg"},
			})
			if err != nil {
				t.Fatalf("could not creates labels plotter: %+v", err)
			}
			for i := range loffm.TextStyle {
				sty := &loffm.TextStyle[i]
				sty.Font.Size = vg.Length(fontSize)
				sty.Color = color.RGBA{B: 255, A: 255}
			}
			loffm.Offset = vg.Point{X: -40, Y: -40}

			m5 := plotter.NewFunction(func(float64) float64 { return -0.5 })
			m5.LineStyle.Color = color.RGBA{R: 255, A: 255}

			l0 := plotter.NewFunction(func(float64) float64 { return 0 })
			l0.LineStyle.Color = color.RGBA{G: 255, A: 255}

			p5 := plotter.NewFunction(func(float64) float64 { return +0.5 })
			p5.LineStyle.Color = color.RGBA{B: 255, A: 255}

			p.Add(labels, lred, m5, l0, p5, loffp, loffm)
			p.Add(plotter.NewGrid())
			p.Add(plotter.NewGlyphBoxes())

			err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/labels_glyphboxes.png")
			if err != nil {
				t.Fatalf("could save plot: %+v", err)
			}
		},
		t, "labels_glyphboxes.png",
	)
}
