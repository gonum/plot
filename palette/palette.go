// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package palette provides basic color palette handling.
package palette

import (
	"image/color"
	"math"
)

// Palette is a collection of colors ordered into a palette.
type Palette interface {
	Colors() []color.Color
}

// DivergingPalette is a collection of colors ordered into a palette with
// a critical class or break in the middle of the color range.
type DivergingPalette interface {
	Palette

	// CriticalIndex returns the indices of the lightest
	// (median) color or colors in the DivergingPalette.
	// The low and high index values will be equal when
	// there is a single median color.
	CriticalIndex() (low, high int)
}

// Hue represents a hue in HSV color space. Valid Hues are within [0, 1].
type Hue float64

const (
	Red Hue = Hue(iota) / 6
	Yellow
	Green
	Cyan
	Blue
	Magenta
)

// Complement returns the complementary hue of a Hue.
func (h Hue) Complement() Hue { return Hue(math.Mod(float64(h+0.5), 1)) }

type palette []color.Color

func (p palette) Colors() []color.Color { return p }

type divergingPalette []color.Color

func (p divergingPalette) Colors() []color.Color { return p }

func (d divergingPalette) CriticalIndex() (low, high int) {
	l := len(d)
	return (l - 1) / 2, l / 2
}

// Rainbow returns a rainbow palette with the specified number of colors, saturation
// value and alpha, and hues in the specified range.
func Rainbow(colors int, start, end Hue, sat, val, alpha float64) Palette {
	p := make(palette, colors)
	hd := float64(end-start) / float64(colors-1)
	c := HSVA{H: float64(start), V: val, S: sat, A: alpha}
	for i := range p {
		p[i] = color.NRGBAModel.Convert(c)
		c.H += hd
	}

	return p
}

// Heat returns a red to yellow palette with the specified number of colors and alpha.
func Heat(colors int, alpha float64) Palette {
	p := make(palette, colors)
	j := colors / 4
	i := colors - j

	hd := float64(Yellow-Red) / float64(i-1)
	c := HSVA{H: float64(Red), V: 1, S: 1, A: alpha}
	for k := range p[:i] {
		p[k] = color.NRGBAModel.Convert(c)
		c.H += hd
	}
	if j == 0 {
		return p
	}

	c.H = float64(Yellow)
	start, end := 1-1/(2*float64(j)), 1/(2*float64(j))
	c.S = start
	sd := (end - start) / float64(j-1)
	for k := range p[i:] {
		p[k+i] = color.NRGBAModel.Convert(c)
		c.S += sd
	}

	return p
}

// Radial return a diverging palette across the specified range, through white and with
// the specified alpha.
func Radial(colors int, start, end Hue, alpha float64) DivergingPalette {
	p := make(divergingPalette, colors)
	h := colors / 2
	c := HSVA{S: 0.5, V: 1, A: alpha}
	ds := 0.5 / float64(h)
	for i := range p[:h] {
		c.H = float64(start)
		p[i] = color.NRGBAModel.Convert(c)
		c.H = float64(end)
		p[len(p)-1-i] = color.NRGBAModel.Convert(c)
		c.S -= ds
	}
	if colors%2 != 0 {
		p[colors/2] = color.NRGBA{0xff, 0xff, 0xff, byte(math.MaxUint8 * alpha)}
	}

	return p
}
