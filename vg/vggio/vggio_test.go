// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vggio

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"runtime"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

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

	p, err := plot.New()
	if err != nil {
		log.Fatalf("could not create plot: %+v", err)
	}
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
