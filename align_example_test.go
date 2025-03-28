// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot_test

import (
	"log"
	"math"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func ExampleAlign() {
	const rows, cols = 4, 3
	plots := make([][]*plot.Plot, rows)
	for j := range rows {
		plots[j] = make([]*plot.Plot, cols)
		for i := range cols {
			if i == 0 && j == 2 {
				// This shows what happens when there are nil plots.
				continue
			}

			p := plot.New()

			if j == 0 && i == 2 {
				// This shows what happens when the axis padding
				// is different among plots.
				p.X.Padding, p.Y.Padding = 0, 0
			}

			if j == 1 && i == 1 {
				// To test the Align function, we make the axis labels
				// on one of the plots stick out.
				p.Y.Max = 1e9
				p.X.Max = 1e9
				p.X.Tick.Label.Rotation = math.Pi / 2
				p.X.Tick.Label.XAlign = draw.XRight
				p.X.Tick.Label.YAlign = draw.YCenter
				p.X.Tick.Label.Font.Size = 8
				p.Y.Tick.Label.Font.Size = 8
			} else {
				p.Y.Max = 1e9
				p.X.Max = 1e9
				p.X.Tick.Label.Font.Size = 1
				p.Y.Tick.Label.Font.Size = 1
			}

			plots[j][i] = p
		}
	}

	img := vgimg.New(vg.Points(150), vg.Points(175))
	dc := draw.New(img)

	t := draw.Tiles{
		Rows:      rows,
		Cols:      cols,
		PadX:      vg.Millimeter,
		PadY:      vg.Millimeter,
		PadTop:    vg.Points(2),
		PadBottom: vg.Points(2),
		PadLeft:   vg.Points(2),
		PadRight:  vg.Points(2),
	}

	canvases := plot.Align(plots, t, dc)
	for j := range rows {
		for i := range cols {
			if plots[j][i] != nil {
				plots[j][i].Draw(canvases[j][i])
			}
		}
	}

	w, err := os.Create("testdata/align.png")
	if err != nil {
		panic(err)
	}
	defer w.Close()
	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		panic(err)
	}
}

func ExampleAxis_labelsPosition() {
	p := plot.New()
	p.Title.Text = "Title"
	p.X.Label.Text = "X [mm]"
	p.Y.Label.Text = "Y [A.U.]"
	p.X.Label.Position = draw.PosRight
	p.Y.Label.Position = draw.PosTop
	p.X.Min = -10
	p.X.Max = +10
	p.Y.Min = -10
	p.Y.Max = +10

	err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/axis_labels.png")
	if err != nil {
		log.Fatalf("could not save plot: %+v", err)
	}
}
