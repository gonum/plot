// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestInvertedScale(t *testing.T) {
	cmpimg.CheckPlot(Example_invertedScale, t, "invertedlogscale.png")
}
