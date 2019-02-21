// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"image/png"
	"log"
	"os"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/vg"
)

const runImageLaTeX = false

// An example of embedding an image in a plot.
func ExampleImage() {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	p.Title.Text = "A Logo"

	// load an image
	f, err := os.Open("testdata/image_plot_input.png")
	if err != nil {
		log.Fatalf("error opening image file: %v\n", err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Fatalf("error decoding image file: %v\n", err)
	}

	p.Add(NewImage(img, 100, 100, 200, 200))

	const (
		w = 5 * vg.Centimeter
		h = 5 * vg.Centimeter
	)

	err = p.Save(5*vg.Centimeter, 5*vg.Centimeter, "testdata/image_plot.png")
	if err != nil {
		log.Fatalf("error saving image plot: %v\n", err)
	}
}

func TestImagePlot(t *testing.T) {
	cmpimg.CheckPlot(ExampleImage, t, "image_plot.png")
}

// An example of embedding an image in a plot with non-linear axes.
func ExampleImage_log() {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	p.Title.Text = "A Logo"

	// load an image
	f, err := os.Open("../../gonum/gopher.png")
	if err != nil {
		log.Fatalf("error opening image file: %v\n", err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Fatalf("error decoding image file: %v\n", err)
	}

	p.Add(NewImage(img, 100, 100, 10000, 10000))

	// Transform axes.
	p.X.Scale = plot.LogScale{}
	p.Y.Scale = plot.LogScale{}
	p.X.Tick.Marker = plot.LogTicks{}
	p.Y.Tick.Marker = plot.LogTicks{}

	const (
		w = 5 * vg.Centimeter
		h = 5 * vg.Centimeter
	)

	err = p.Save(5*vg.Centimeter, 5*vg.Centimeter, "testdata/image_plot_log.png")
	if err != nil {
		log.Fatalf("error saving image plot: %v\n", err)
	}
}

func TestImagePlot_log(t *testing.T) {
	cmpimg.CheckPlot(ExampleImage_log, t, "image_plot_log.png")
}
