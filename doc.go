// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package plot provides an API for setting up plots, and primitives for
// drawing on plots.
//
// Plot is the basic type for creating a plot, setting the title, axis
// labels, legend, tick marks, etc.  Types implementing the Plotter
// interface can draw to the data area of a plot using the primitives
// made available by this package.  Some standard implementations
// of the Plotter interface can be found in the
// github.com/gonum/plot/plotter package
// which is documented here:
// https://godoc.org/github.com/gonum/plot/plotter
package plot // import "github.com/gonum/plot"
