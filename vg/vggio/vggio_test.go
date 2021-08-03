// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vggio

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"runtime"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

const deltaGio = 0.05 // empirical value from experimentation.

// init makes sure the headless display is ready for tests with Gio.
// On GitHub Actions and on linux, that headless display may take some time to
// be properly available and appears to be setup "on demand".
// So we request it by trying to take a screenshot twice:
//  - the first time around might fail
//  - the second time shouldn't.
func init() {
	if runtime.GOOS != "linux" {
		return
	}

	const (
		w   = 20 * vg.Centimeter
		h   = 15 * vg.Centimeter
		dpi = 96
	)
	gtx := layout.Context{
		Ops: new(op.Ops),
		Constraints: layout.Exact(image.Pt(
			int(w.Dots(dpi)),
			int(h.Dots(dpi)),
		)),
	}

	var err error
	for try := 0; try < 2; try++ {
		_, err = New(gtx, w, h, UseDPI(dpi)).Screenshot()
		if err == nil {
			return
		}
	}

	panic(fmt.Errorf("vg/vggio_test: could not setup headless display: %+v", err))
}

func TestCanvas(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("TODO: github actions for darwin with headless setup.")
	}

	const fname = "testdata/func.png"

	const (
		w   = 20 * vg.Centimeter
		h   = 15 * vg.Centimeter
		dpi = 96
	)

	cmpimg.CheckPlotApprox(func() {
		p := plot.New()
		p.Title.Text = "My title"
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"

		quad := plotter.NewFunction(func(x float64) float64 { return x * x })
		quad.Color = color.RGBA{B: 255, A: 255}

		exp := plotter.NewFunction(func(x float64) float64 { return math.Pow(2, x) })
		exp.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
		exp.Width = vg.Points(2)
		exp.Color = color.RGBA{G: 255, A: 255}

		sin := plotter.NewFunction(func(x float64) float64 { return 10*math.Sin(x) + 50 })
		sin.Dashes = []vg.Length{vg.Points(4), vg.Points(5)}
		sin.Width = vg.Points(4)
		sin.Color = color.RGBA{R: 255, A: 255}

		p.Add(quad, exp, sin)
		p.Legend.Add("x^2", quad)
		p.Legend.Add("2^x", exp)
		p.Legend.Add("10*sin(x)+50", sin)
		p.Legend.ThumbnailWidth = 0.5 * vg.Inch

		p.X.Min = 0
		p.X.Max = 10
		p.Y.Min = 0
		p.Y.Max = 100

		p.Add(plotter.NewGrid())

		gtx := layout.Context{
			Ops: new(op.Ops),
			Constraints: layout.Exact(image.Pt(
				int(w.Dots(dpi)),
				int(h.Dots(dpi)),
			)),
		}
		cnv := New(gtx, w, h, UseDPI(dpi))
		p.Draw(draw.New(cnv))

		img, err := cnv.Screenshot()
		if err != nil {
			t.Fatalf("could not create screenshot: %+v", err)
		}
		f, err := os.Create(fname)
		if err != nil {
			t.Fatalf("could not create output file: %+v", err)
		}
		defer f.Close()

		err = png.Encode(f, img)
		if err != nil {
			t.Fatalf("could not encode screenshot: %+v", err)
		}

		err = f.Close()
		if err != nil {
			t.Fatalf("could not save screenshot: %+v", err)
		}
	}, t, deltaGio, "func.png",
	)
}

func TestCollectionName(t *testing.T) {
	for _, tc := range []struct {
		name string
		want string
	}{
		{"Liberation", "Liberation"},
		{"LiberationSerif-Bold", "LiberationSerif"},
		{"LiberationSerif-BoldItalic", "LiberationSerif"},
		{"LiberationSerif-BoldItalic-Extra", "LiberationSerif"},

		{"LiberationMono", "LiberationMono"},
		{"LiberationMono-Regular", "LiberationMono"},

		{"Times-Roman", "Times"},
		{"Times-Bold", "Times"},
	} {
		got := collectionName(tc.name)
		if got != tc.want {
			t.Errorf(
				"%s: invalid collection name: got=%q, want=%q",
				tc.name, got, tc.want,
			)
		}
	}
}

func TestLabels(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("TODO: github actions for darwin with headless setup.")
	}

	const fname = "testdata/labels.png"

	const (
		w   = 20 * vg.Centimeter
		h   = 15 * vg.Centimeter
		dpi = 96
	)

	cmpimg.CheckPlotApprox(func() {
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
			sty.Font.Size = vg.Length(34)
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
			sty.Font.Size = vg.Length(34)
			sty.Color = color.RGBA{R: 255, A: 255}
			sty.YAlign = draw.YBottom
		}

		m5 := plotter.NewFunction(func(float64) float64 { return -0.5 })
		m5.LineStyle.Color = color.RGBA{R: 255, A: 255}

		l0 := plotter.NewFunction(func(float64) float64 { return 0 })
		l0.LineStyle.Color = color.RGBA{G: 255, A: 255}

		p5 := plotter.NewFunction(func(float64) float64 { return +0.5 })
		p5.LineStyle.Color = color.RGBA{B: 255, A: 255}

		p.Add(labels, lred, m5, l0, p5)
		p.Add(plotter.NewGrid())
		p.Add(plotter.NewGlyphBoxes())

		gtx := layout.Context{
			Ops: new(op.Ops),
			Constraints: layout.Exact(image.Pt(
				int(w.Dots(dpi)),
				int(h.Dots(dpi)),
			)),
		}
		cnv := New(gtx, w, h, UseDPI(dpi))
		p.Draw(draw.New(cnv))

		img, err := cnv.Screenshot()
		if err != nil {
			t.Fatalf("could not create screenshot: %+v", err)
		}
		f, err := os.Create(fname)
		if err != nil {
			t.Fatalf("could not create output file: %+v", err)
		}
		defer f.Close()

		err = png.Encode(f, img)
		if err != nil {
			t.Fatalf("could not encode screenshot: %+v", err)
		}

		err = f.Close()
		if err != nil {
			t.Fatalf("could not save screenshot: %+v", err)
		}
	}, t, deltaGio, "labels.png",
	)
}
