// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func ExampleLabels() {
	p := plot.New()
	p.Title.Text = "Labels"
	p.X.Min = -10
	p.X.Max = +10
	p.Y.Min = 0
	p.Y.Max = +20

	labels, err := plotter.NewLabels(plotter.XYLabels{
		XYs: []plotter.XY{
			{X: -5, Y: 5},
			{X: +5, Y: 5},
			{X: +5, Y: 15},
			{X: -5, Y: 15},
		},
		Labels: []string{"A", "B", "C", "D"},
	})
	if err != nil {
		log.Fatalf("could not creates labels plotter: %+v", err)
	}

	p.Add(labels)

	err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/labels.png")
	if err != nil {
		log.Fatalf("could save plot: %+v", err)
	}
}

// ExampleLabels_inCanvasCoordinates shows how to write a label
// in a plot using the canvas coordinates (instead of the data coordinates.)
// It can be useful in the situation where one wants to draw some label
// always in the same location of a plot, irrespective of the minute
// details of a particular plot data range.
func ExampleLabels_inCanvasCoordinates() {
	p := plot.New()
	p.Title.Text = "Labels - X"
	p.X.Min = -10
	p.X.Max = +10
	p.Y.Min = 0
	p.Y.Max = +20

	labels, err := plotter.NewLabels(plotter.XYLabels{
		XYs: []plotter.XY{
			{X: -5, Y: 5},
			{X: +5, Y: 5},
			{X: +5, Y: 15},
			{X: -5, Y: 15},
		},
		Labels: []string{"A", "B", "C", "D"},
	},
	)
	if err != nil {
		log.Fatalf("could not creates labels plotter: %+v", err)
	}

	p.Add(labels)

	f, err := os.Create("testdata/labels_cnv_coords.png")
	if err != nil {
		log.Fatalf("could not create output plot file: %+v", err)
	}
	defer f.Close()

	cnv := vgimg.PngCanvas{
		Canvas: vgimg.New(10*vg.Centimeter, 10*vg.Centimeter),
	}

	dc := draw.New(cnv)
	p.Draw(dc)

	// Put an 'X' in the middle of the data-canvas.
	{
		fnt := p.TextHandler.Cache().Lookup(plotter.DefaultFont, vg.Points(12))
		da := p.DataCanvas(dc)
		da.FillString(fnt, vg.Point{X: da.X(0.5), Y: da.Y(0.5)}, "X")
	}

	_, err = cnv.WriteTo(f)
	if err != nil {
		log.Fatalf("could not write to output plot file: %+v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could save plot: %+v", err)
	}
}
