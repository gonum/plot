// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgimg_test

import (
	"image"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func ExampleUseDPI() {
	p := plot.New()
	p.Title.Text = "Title"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	const (
		width  = 10 * vg.Centimeter
		height = 10 * vg.Centimeter
	)

	// Create a new canvas with the given dimensions,
	// and specify an explicit DPI value (dot per inch)
	// for the plot.
	c := vgimg.NewWith(
		vgimg.UseWH(width, height),
		vgimg.UseDPI(72),
	)

	dc := draw.New(c)
	p.Draw(dc)
}

func ExampleUseImage() {
	p := plot.New()
	p.Title.Text = "Title"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	img := image.NewRGBA(image.Rect(0, 100, 0, 100))
	// Create a new canvas using the given image to specify the dimensions
	// of the plot.
	//
	// Note that modifications applied to the canvas will not be reflected on
	// the input image.
	c := vgimg.NewWith(
		vgimg.UseImage(img),
	)

	dc := draw.New(c)
	p.Draw(dc)
}

func ExampleUseBackgroundColor() {
	p := plot.New()
	p.Title.Text = "Title"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	const (
		width  = 10 * vg.Centimeter
		height = 10 * vg.Centimeter
	)

	// Create a new canvas with the given dimensions,
	// and specify an explicit background color for the plot.
	c := vgimg.NewWith(
		vgimg.UseWH(width, height),
		vgimg.UseBackgroundColor(color.Transparent),
	)

	dc := draw.New(c)
	p.Draw(dc)
}
