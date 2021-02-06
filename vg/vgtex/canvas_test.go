// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgtex_test

import (
	"image/color"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func TestTexCanvas(t *testing.T) {
	cmpimg.CheckPlot(Example, t, "scatter.tex")
}

func TestLineLatex(t *testing.T) {
	test := func(fname string) func() {
		return func() {
			p := plot.New()
			p.X.Min = -10
			p.X.Max = +10
			p.Y.Min = -10
			p.Y.Max = +10

			f1 := plotter.NewFunction(func(float64) float64 {
				return -7
			})
			f1.LineStyle.Color = color.Black
			f1.LineStyle.Width = 2
			f1.LineStyle.Dashes = []vg.Length{2, 1}

			f2 := plotter.NewFunction(func(float64) float64 {
				return -1
			})
			f2.LineStyle.Color = color.RGBA{R: 255, A: 255}
			f2.LineStyle.Width = 2
			f2.LineStyle.Dashes = []vg.Length{4, 2}

			f3 := plotter.NewFunction(func(float64) float64 {
				return +7
			})
			f3.LineStyle.Color = color.Black
			f3.LineStyle.Width = 2
			f3.LineStyle.Dashes = []vg.Length{2, 1}

			p.Add(f1, f2, f3)
			p.Add(plotter.NewGrid())

			const size = 5 * vg.Centimeter
			err := p.Save(size, size, fname)
			if err != nil {
				t.Fatalf("error: %+v", err)
			}
		}
	}
	cmpimg.CheckPlot(test("testdata/linestyle.tex"), t, "linestyle.tex")
	cmpimg.CheckPlot(test("testdata/linestyle.png"), t, "linestyle.png")
}

func TestFillStyle(t *testing.T) {
	cmpimg.CheckPlot(func() {
		p := plot.New()
		p.Title.Text = "Fill style"
		p.Legend.Top = true
		p.Legend.Left = true

		const n = 10
		xys := make(plotter.XYs, n)
		for i := range xys {
			xys[i].X = float64(i)
			xys[i].Y = float64(i)
		}
		h, err := plotter.NewHistogram(xys, n)
		if err != nil {
			t.Fatalf("could not create histogram: %+v", err)
		}
		h.FillColor = color.NRGBA{R: 255, A: 100}
		h.LineStyle.Color = color.Transparent
		p.Add(h)
		p.Legend.Add("h", h)

		const size = 5 * vg.Centimeter
		err = p.Save(size, size, "testdata/fillstyle.tex")
		if err != nil {
			t.Fatalf("error: %+v", err)
		}
	}, t, "fillstyle.tex")
}
