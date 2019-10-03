// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

const runImageLaTeX = false

func TestImagePlot(t *testing.T) {
	cmpimg.CheckPlot(ExampleImage, t, "image_plot.png")
}

func TestImagePlot_log(t *testing.T) {
	cmpimg.CheckPlot(ExampleImage_log, t, "image_plot_log.png")
}
