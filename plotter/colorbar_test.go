// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestColorBar_horizontal(t *testing.T) {
	cmpimg.CheckPlot(ExampleColorBar_horizontal, t, "colorBarHorizontal.png")
}

func TestColorBar_horizontal_log(t *testing.T) {
	cmpimg.CheckPlot(ExampleColorBar_horizontal_log, t, "colorBarHorizontalLog.png")
}

func TestColorBar_vertical(t *testing.T) {
	cmpimg.CheckPlot(ExampleColorBar_vertical, t, "colorBarVertical.png")
}
