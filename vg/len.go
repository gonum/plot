// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package vg

// TODO(eaburns): These should be more like time.Duration.  I.e., 5*vg.Inch, 6*vg.Centimeter, etc., instead of requiring a function for each.

const (
	// PtInch is the number of points in an inch.
	ptInch = 72

	// PtCentimeter is the number of points in a centimeter.
	ptCentimeter = 28.3464567

	// PtMillimeter is the number of ponints in a millimeter.
	ptMillimeter = ptCentimeter / 10
)

// A Length is a unit-independent representation of length.
// Internally, the length is stored in postscript points.
type Length float64

// Points returns a length for the given number of points.
func Points(pt float64) Length {
	return Length(pt)
}

// Inches returns a length for the given numer of inches.
func Inches(in float64) Length {
	return Length(in * ptInch)
}

// Centimeters returns a length for the given number of centimeters.
func Centimeters(cm float64) Length {
	return Length(cm * ptCentimeter)
}

// Millimeters returns a length for the given number of millimeters.
func Millimeters(mm float64) Length {
	return Length(mm * ptMillimeter)
}

// Dots returns the length in dots for the given Canvas.
func (l Length) Dots(c Canvas) float64 {
	return float64(l) / ptInch * c.DPI()
}

// Points returns the length in postscript points.
func (l Length) Points() float64 {
	return float64(l)
}

// Inches returns the length in inches.
func (l Length) Inches() float64 {
	return float64(l) / ptInch
}

// Centimeters returns the length in centimeters.
func (l Length) Centimeters() float64 {
	return float64(l) / ptCentimeter
}

// Millimeters returns the length in millimeters.
func (l Length) Millimeters() float64 {
	return float64(l) / ptMillimeter
}
