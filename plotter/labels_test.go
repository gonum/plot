// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestLabels(t *testing.T) {
	cmpimg.CheckPlot(ExampleLabels, t, "labels.png")
	cmpimg.CheckPlot(ExampleLabels_inCanvasCoordinates, t, "labels_cnv_coords.png")
}
