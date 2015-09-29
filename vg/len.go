// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vg

// A Length is a unit-independent representation of length.
// Internally, the length is stored in postscript points.
type Length float64

// Points returns a length for the given number of points.
func Points(pt float64) Length {
	return Length(pt)
}

// Common lengths.
const (
	Inch       Length = 72
	Centimeter        = Inch / 2.54
	Millimeter        = Centimeter / 10
)

// Dots returns the length in dots for the given resolution.
func (l Length) Dots(dpi float64) float64 {
	return float64(l) / Inch.Points() * dpi
}

// Points returns the length in postscript points.
func (l Length) Points() float64 {
	return float64(l)
}
