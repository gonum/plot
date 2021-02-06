// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vg_test

import (
	"image/color"
	"math"
	"reflect"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/recorder"
)

func TestMultiCanvas(t *testing.T) {
	p := plot.New()
	p.Title.Text = "Title"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"
	p.X.Min = -2 * math.Pi
	p.X.Max = +2 * math.Pi

	sin := plotter.NewFunction(math.Sin)
	sin.LineStyle.Color = color.RGBA{R: 255, A: 255}
	sin.LineStyle.Dashes = plotutil.Dashes(1)

	cos := plotter.NewFunction(math.Cos)
	cos.LineStyle.Color = color.RGBA{B: 255, A: 255}
	cos.LineStyle.Dashes = plotutil.Dashes(2)

	p.Add(sin, cos, plotter.NewGrid())

	c1 := new(recorder.Canvas)
	c2 := new(recorder.Canvas)

	const (
		width  = 10 * vg.Centimeter
		height = 10 * vg.Centimeter
	)

	p.Draw(draw.NewCanvas(
		vg.MultiCanvas(c1, c2, vg.MultiCanvas()),
		width, height,
	))

	if !reflect.DeepEqual(c1, c2) {
		t.Fatalf("tee canvas failed to replicate drawing calls")
	}
}
