// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package palette_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestReverse(t *testing.T) {
	cmpimg.CheckPlot(ExampleReverse, t, "reverse.png")
}

func TestReverse_Palette(t *testing.T) {
	cmpimg.CheckPlot(ExampleReverse_palette, t, "reverse_palette.png")
}
