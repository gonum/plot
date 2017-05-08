// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package palette

import (
	"fmt"
	"image/color"
	"math"

	"github.com/gonum/floats"
)

// A ColorMap maps scalar values to colors.
type ColorMap interface {
	// At returns the color associated with the given value.
	// If the value is not between Max() and Min(), an error is returned.
	At(float64) (color.Color, error)

	// Max returns the current maximum value of the ColorMap.
	Max() float64

	// SetMax sets the maximum value of the ColorMap.
	SetMax(float64)

	// Min returns the current minimum value of the ColorMap.
	Min() float64

	// SetMin sets the minimum value of the ColorMap.
	SetMin(float64)

	// Alpha returns the opacity value of the ColorMap.
	Alpha() float64

	// SetAlpha sets the opacity value of the ColorMap. Zero is transparent
	// and one is completely opaque. The default value of alpha should be
	// expected to be one. The function should be expected to panic
	// if alpha is not between zero and one.
	SetAlpha(float64)

	// Palette creates a Palette with the specified number of colors
	// from the ColorMap.
	Palette(colors int) Palette
}

// DivergingColorMap maps scalar values to colors that diverge
// from a central value.
type DivergingColorMap interface {
	ColorMap

	// SetConvergePoint sets the value where the diverging colors
	// should meet. The default value should be expected to be
	// (Min() + Max()) / 2. It should be expected that calling either
	// SetMax() or SetMin() will set a new default value, so for a
	// custom convergence point this function should be called after
	// SetMax() and SetMin(). The function should be expected to panic
	// if the value is not between Min() and Max().
	SetConvergePoint(float64)

	// ConvergePoint returns the value where the diverging colors meet.
	ConvergePoint() float64
}

// FromPalette creates a ColorMap that maps values to colors based
// on RGB interpolation among the colors in p.
func FromPalette(p Palette) ColorMap {
	scalars := make([]float64, len(p.Colors()))
	return &colorMap{
		colors:  p.Colors(),
		scalars: floats.Span(scalars, 0, 1),
		alpha:   1,
	}
}

type colorMap struct {
	// colors are the colors in the underlying Palette.
	colors []color.Color

	// scalars are len(colors) equally spaced points between
	// 0 and 1, used for interpolation.
	scalars []float64

	// alpha represents the opacity of the returned
	// colors in the range (0,1). It is set to 1 by default.
	alpha float64

	// min and max are the minimum and maximum values of the range of scalars
	// that can be mapped to colors using this ColorMap.
	min, max float64
}

// At implements the ColorMap interface for a colorMap value.
func (cm *colorMap) At(v float64) (color.Color, error) {
	if err := checkRange(cm.min, cm.max, v); err != nil {
		return nil, err
	}
	scalar := (v - cm.min) / cm.max
	i := searchFloat64s(cm.scalars, scalar)
	if i == 0 {
		return cm.colors[i], nil
	}
	c1 := color.NRGBA64Model.Convert(cm.colors[i-1]).(color.NRGBA64)
	c2 := color.NRGBA64Model.Convert(cm.colors[i]).(color.NRGBA64)
	frac := (scalar - cm.scalars[i-1]) / (cm.scalars[i] - cm.scalars[i-1])
	return color.NRGBA64{
		R: uint16(frac*(float64(c2.R)-float64(c1.R)) + float64(c1.R)),
		G: uint16(frac*(float64(c2.G)-float64(c1.G)) + float64(c1.G)),
		B: uint16(frac*(float64(c2.B)-float64(c1.B)) + float64(c1.B)),
		A: uint16(float64(math.MaxUint16) * cm.alpha),
	}, nil
}

func checkRange(min, max, val float64) error {
	if max == min {
		return fmt.Errorf("palette: color map max == min == %g", max)
	}
	if min > max {
		return fmt.Errorf("palette: color map max (%g) < min (%g)", max, min)
	}
	if val < min {
		return ErrUnderflow
	}
	if val > max {
		return ErrOverflow
	}
	if math.IsNaN(val) {
		return ErrNaN
	}
	return nil
}

// searchFloat64s acts the same as sort.SearchFloat64s, except
// it uses a simple search algorithm instead of binary search.
func searchFloat64s(vals []float64, val float64) int {
	for j, v := range vals {
		if val <= v {
			return j
		}
	}
	return len(vals)
}

// SetMax implements the palette.ColorMap interface for a colorMap value.
func (cm *colorMap) SetMax(v float64) {
	cm.max = v
}

// SetMin implements the palette.ColorMap interface for a colorMap value.
func (cm *colorMap) SetMin(v float64) {
	cm.min = v
}

// Max implements the palette.ColorMap interface for a colorMap value.
func (cm *colorMap) Max() float64 {
	return cm.max
}

// Min implements the palette.ColorMap interface for a colorMap value.
func (cm *colorMap) Min() float64 {
	return cm.min
}

// SetAlpha sets the opacity value of this color map. Zero is transparent
// and one is completely opaque.
// The function will panic is alpha is not between zero and one.
func (cm *colorMap) SetAlpha(alpha float64) {
	if !inUnitRange(alpha) {
		panic(fmt.Errorf("palette: invalid alpha: %g", alpha))
	}
	cm.alpha = alpha
}

func inUnitRange(v float64) bool { return 0 <= v && v <= 1 }

// Alpha returns the opacity value of this color map.
func (cm *colorMap) Alpha() float64 {
	return cm.alpha
}

// Palette returns a value that fulfills the Palette interface,
// where n is the number of desired colors.
func (cm colorMap) Palette(n int) Palette {
	if cm.Max() == 0 && cm.Min() == 0 {
		cm.SetMin(0)
		cm.SetMax(1)
	}
	delta := (cm.max - cm.min) / float64(n-1)
	var v float64
	c := make([]color.Color, n)
	for i := 0; i < n; i++ {
		v = cm.min + delta*float64(i)
		var err error
		c[i], err = cm.At(v)
		if err != nil {
			panic(err)
		}
	}
	return palette(c)
}
