// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Some descriptive text is taken from http://www.kennethmoreland.com/,
// ©2009--2016 Kenneth Moreland.

// Color conversion functions are largely based on a spreadsheet
// (http://www.kennethmoreland.com/color-maps/DivergingColorMapWorkshop.xls)
// and a Python script (http://www.kennethmoreland.com/color-maps/diverging_map.py),
// which are copyright their respective authors.

// Package moreland provides color maps for pseudocoloring scalar fields.
// The color maps are described at http://www.kennethmoreland.com/color-advice/
// and in the following publications:
//
// "Why We Use Bad Color Maps and What You Can Do About It." Kenneth Moreland.
// In Proceedings of Human Vision and Electronic Imaging (HVEI), 2016. (To appear)
//
// "Diverging Color Maps for Scientific Visualization." Kenneth Moreland.
// In Proceedings of the 5th International Symposium on Visual Computing,
// December 2009. DOI 10.1007/978-3-642-10520-3_9.
package moreland

import (
	"fmt"
	"image/color"
	"math"

	"github.com/gonum/plot/palette"
)

// smoothDiverging is a smooth diverging color palette as described in
// "Diverging Color Maps for Scientific Visualization." by Kenneth Moreland,
// in Proceedings of the 5th International Symposium on Visual Computing,
// December 2009. DOI 10.1007/978-3-642-10520-3_9.
type smoothDiverging struct {
	// start and end are the beginning and ending colors
	start, end msh

	// convergeM is the MSH magnitude of the convergence point.
	// It is 88 by default.
	convergeM float64

	// alpha represents the opacity of the returned
	// colors in the range (0,1). It is 1 by default.
	alpha float64

	// min and max are the minimum and maximum values of the range of
	// scalars that can be mapped to colors using this palette.
	min, max float64

	// convergePoint is a number between min and max where the colors
	// should converge.
	convergePoint float64
}

// NewSmoothDiverging creates a new smooth diverging ColorMap as described in
// "Diverging Color Maps for Scientific Visualization." by Kenneth Moreland,
// in Proceedings of the 5th International Symposium on Visual Computing,
// December 2009. DOI 10.1007/978-3-642-10520-3_9.
//
// start and end are the start- and end-point colors and
// convergeM is the magnitude of the convergence point in
// magnitude-saturation-hue (MSH) color space. Note that
// convergeM specifies the color of the convergence point; it does not
// specify the location of the convergence point.
func NewSmoothDiverging(start, end color.Color, convergeM float64) palette.DivergingColorMap {
	return newSmoothDiverging(colorToMSH(start), colorToMSH(end), convergeM)
}

// newSmoothDiverging creates a new smooth diverging ColorMap
// where start and end are the start and end point colors in MSH space and
// convergeM is the MSH magnitude of the convergence point. Note that
// convergeM specifies the color of the convergence point; it does not
// specify the location of the convergence point.
func newSmoothDiverging(start, end msh, convergeM float64) palette.DivergingColorMap {
	return &smoothDiverging{
		start:         start,
		end:           end,
		convergeM:     convergeM,
		convergePoint: math.NaN(),
		alpha:         1,
	}
}

// At implements the palette.ColorMap interface.
func (p *smoothDiverging) At(v float64) (color.Color, error) {
	if err := checkRange(p.min, p.max, v); err != nil {
		return nil, err
	}
	convergePoint := (p.convergePoint - p.min) / (p.max - p.min)
	scalar := (v - p.min) / (p.max - p.min)
	o := p.interpolateMSHDiverging(scalar, convergePoint).cieLAB().cieXYZ().rgb().sRGBA(p.alpha)
	if !inUnitRange(o.R) || !inUnitRange(o.G) || !inUnitRange(o.B) || !inUnitRange(o.A) {
		return nil, fmt.Errorf("moreland: invalid color r:%g, g:%g, b:%g, a:%g", o.R, o.G, o.B, o.A)
	}
	return o, nil
}

func inUnitRange(v float64) bool { return 0 <= v && v <= 1 }

// SetMax implements the palette.ColorMap interface.
func (p *smoothDiverging) SetMax(v float64) {
	p.max = v
	p.convergePoint = (p.min + p.max) / 2
}

// SetMin implements the palette.ColorMap interface.
func (p *smoothDiverging) SetMin(v float64) {
	p.min = v
	p.convergePoint = (p.min + p.max) / 2
}

// Max implements the palette.ColorMap interface.
func (p *smoothDiverging) Max() float64 {
	return p.max
}

// Min implements the palette.ColorMap interface.
func (p *smoothDiverging) Min() float64 {
	return p.min
}

// SetAlpha sets the opacity value of this color map. Zero is transparent
// and one is completely opaque.
// The function will panic is alpha is not between zero and one.
func (p *smoothDiverging) SetAlpha(alpha float64) {
	if !inUnitRange(alpha) {
		panic(fmt.Errorf("invalid alpha: %g", alpha))
	}
	p.alpha = alpha
}

// Alpha returns the opacity value of this color map.
func (p *smoothDiverging) Alpha() float64 {
	return p.alpha
}

// SetConvergePoint sets the value where the diverging colors
// should meet.
func (p *smoothDiverging) SetConvergePoint(val float64) {
	if val > p.Max() || val < p.Min() {
		panic(fmt.Errorf("moreland: convergence point (%g) must be between min (%g) and max (%g)",
			val, p.Min(), p.Max()))
	}
	p.convergePoint = val
}

// ConvergePoint returns the value where the diverging colors meet.
func (p *smoothDiverging) ConvergePoint() float64 {
	return p.convergePoint
}

// interpolateMSHDiverging performs a color interpolation through MSH space,
// where scalar is a number between 0 and 1 that the
// color should be evaluated at, and convergePoint is a number between 0 and
// 1 where the colors should converge.
func (p *smoothDiverging) interpolateMSHDiverging(scalar, convergePoint float64) msh {
	startHTwist := hueTwist(p.start, p.convergeM)
	endHTwist := hueTwist(p.end, p.convergeM)
	if scalar < convergePoint {
		// interpolation factor
		interp := scalar / convergePoint
		return msh{
			M: (p.convergeM-p.start.M)*interp + p.start.M,
			S: p.start.S * (1 - interp),
			H: p.start.H + startHTwist*interp,
		}
	}
	// interpolation factors
	interp1 := (scalar - 1) / (convergePoint - 1)
	interp2 := (scalar/convergePoint - 1)
	var H float64
	if scalar > convergePoint {
		H = p.end.H + endHTwist*interp1
	}
	return msh{
		M: (p.convergeM-p.end.M)*interp1 + p.end.M,
		S: p.end.S * interp2,
		H: H,
	}
}

// Palette returns a palette.Palette with the specified number of colors.
func (p smoothDiverging) Palette(n int) palette.Palette {
	if p.Max() == 0 && p.Min() == 0 {
		p.SetMin(0)
		p.SetMax(1)
	}
	delta := (p.max - p.min) / float64(n-1)
	var v float64
	c := make([]color.Color, n)
	for i := range c {
		v = p.min + delta*float64(i)
		var err error
		c[i], err = p.At(v)
		if err != nil {
			panic(err)
		}
		v += delta
	}
	return plte(c)
}

// SmoothBlueRed is a SmoothDiverging-class ColorMap ranging from blue to red.
func SmoothBlueRed() palette.DivergingColorMap {
	start := msh{
		M: 80,
		S: 1.08,
		H: -1.1,
	}
	end := msh{
		M: 80,
		S: 1.08,
		H: 0.5,
	}
	return newSmoothDiverging(start, end, 88)
}

// SmoothPurpleOrange is a SmoothDiverging-class ColorMap ranging from purple to orange.
func SmoothPurpleOrange() palette.DivergingColorMap {
	start := msh{
		M: 64.97539711,
		S: 0.899434815,
		H: -0.899431964,
	}
	end := msh{
		M: 85.00850996,
		S: 0.949730284,
		H: 0.950636521,
	}
	return newSmoothDiverging(start, end, 88)
}

// SmoothGreenPurple is a SmoothDiverging-class ColorMap ranging from green to purple.
func SmoothGreenPurple() palette.DivergingColorMap {
	start := msh{
		M: 78.04105346,
		S: 0.885011982,
		H: 2.499491379,
	}
	end := msh{
		M: 64.97539711,
		S: 0.899434815,
		H: -0.899431964,
	}
	return newSmoothDiverging(start, end, 88)
}

// SmoothBlueTan is a SmoothDiverging-class ColorMap ranging from blue to tan.
func SmoothBlueTan() palette.DivergingColorMap {
	start := msh{
		M: 79.94788321,
		S: 0.798754784,
		H: -1.401313221,
	}
	end := msh{
		M: 80.07193125,
		S: 0.799798811,
		H: 1.401089787,
	}
	return newSmoothDiverging(start, end, 88)
}

// SmoothGreenRed is a SmoothDiverging-class ColorMap ranging from green to red.
func SmoothGreenRed() palette.DivergingColorMap {
	start := msh{
		M: 78.04105346,
		S: 0.885011982,
		H: 2.499491379,
	}
	end := msh{
		M: 76.96722122,
		S: 0.949483656,
		H: 0.499492043,
	}
	return newSmoothDiverging(start, end, 88)
}
