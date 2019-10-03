// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestTimeSeries(t *testing.T) {
	cmpimg.CheckPlot(Example_timeSeries, t, "timeseries.png")
}
