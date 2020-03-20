// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgimg_test

import (
	"image"
	"image/color"
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func ExampleUseDPI() {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	p.Title.Text = "Title"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	img := vgimg.NewWith(
		vgimg.UseWH(10*vg.Centimeter, 10*vg.Centimeter),
		vgimg.UseDPI(72),
	)

	dc := draw.New(img)
	p.Draw(dc)
}

func ExampleUseImage() {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	p.Title.Text = "Title"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	img := image.NewRGBA(image.Rect(0, 100, 0, 100))
	c := vgimg.NewWith(
		vgimg.UseImage(img),
	)

	dc := draw.New(c)
	p.Draw(dc)
}

func ExampleUseBackgroundColor() {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	p.Title.Text = "Title"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	c := vgimg.NewWith(
		vgimg.UseWH(10*vg.Centimeter, 10*vg.Centimeter),
		vgimg.UseBackgroundColor(color.Transparent),
	)

	dc := draw.New(c)
	p.Draw(dc)
}
