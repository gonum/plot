// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package palette provides basic color palette handling.
package palette

import (
	"errors"
	"fmt"
	"image/color"
	"math"
)

// Palette is a collection of colors ordered into a palette.
type Palette interface {
	Colors() []color.Color
}

// New returns a new palette from the specified colors
func New(colors ...color.Color) Palette {
	return palette(colors)
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

// A ColorMapInt maps an integer category value to a color.
// If there is no mapped color for the given category, an error is returned.
type ColorMapInt interface {
	// At returns the color associated with category cat.
	At(cat int) (color.Color, error)
}

// IntMap fulfils the ColorMapInt interface, mapping integer
// categories to colors.
type IntMap struct {
	Categories []int
	Colors     []color.Color
}

// At fulfils the ColorMapInt interface.
func (im *IntMap) At(cat int) (color.Color, error) {
	if len(im.Categories) != len(im.Colors) {
		panic(fmt.Errorf("palette: number of categories (%d) != number of colors (%d)", len(im.Categories), len(im.Colors)))
	}
	if i := searchInts(im.Categories, cat); i != -1 {
		return im.Colors[i], nil
	}
	return nil, fmt.Errorf("palette: category '%d' not found", cat)
}

// searchInts returns the index of ints that matches i, or -1 if there
// is no match.
func searchInts(ints []int, i int) int {
	for ii, iii := range ints {
		if iii == i {
			return ii
		}
	}
	return -1
}

// A ColorMapString maps a string category value to a color.
// If there is no mapped color for the given category, an error is returned.
type ColorMapString interface {
	// At returns the color associated with category cat.
	At(cat string) (color.Color, error)
}

// StringMap fulfils the ColorMapString interface, mapping integer
// categories to colors.
type StringMap struct {
	Categories []string
	Colors     []color.Color
}

// At fulfils the ColorMapInt interface.
func (sm *StringMap) At(cat string) (color.Color, error) {
	if len(sm.Categories) != len(sm.Colors) {
		panic(fmt.Errorf("palette: number of categories (%d) != number of colors (%d)", len(sm.Categories), len(sm.Colors)))
	}
	if i := searchStrings(sm.Categories, cat); i != -1 {
		return sm.Colors[i], nil
	}
	return nil, fmt.Errorf("palette: category '%s' not found", cat)
}

// searchStrings returns the index of strs that matches str, or -1 if there
// is no match.
func searchStrings(strs []string, str string) int {
	for i, s := range strs {
		if s == str {
			return i
		}
	}
	return -1
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

var (
	// ErrOverflow is the error returned by ColorMaps when the specified
	// value is greater than the maximum value.
	ErrOverflow = errors.New("palette: specified value > maximum")

	// ErrUnderflow is the error returned by ColorMaps when the specified
	// value is less than the minimum value.
	ErrUnderflow = errors.New("palette: specified value < minimum")

	// ErrNaN is the error returned by ColorMaps when the specified
	// value is NaN.
	ErrNaN = errors.New("palette: specified value == NaN")
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
	c := HSVA{V: val, S: sat, A: alpha}
	for i := range p {
		c.H = float64(start) + float64(i)*hd
		p[i] = color.NRGBAModel.Convert(c)
	}

	return p
}

// Heat returns a red to yellow palette with the specified number of colors and alpha.
func Heat(colors int, alpha float64) Palette {
	p := make(palette, colors)
	j := colors / 4
	i := colors - j

	hd := float64(Yellow-Red) / float64(i-1)
	c := HSVA{V: 1, S: 1, A: alpha}
	for k := range p[:i] {
		c.H = float64(Red) + float64(k)*hd
		p[k] = color.NRGBAModel.Convert(c)
	}
	if j == 0 {
		return p
	}

	c.H = float64(Yellow)
	start, end := 1-1/(2*float64(j)), 1/(2*float64(j))
	c.S = start
	sd := (end - start) / float64(j-1)
	for k := range p[i:] {
		c.S = start + float64(k)*sd
		p[k+i] = color.NRGBAModel.Convert(c)
	}

	return p
}

// Radial return a diverging palette across the specified range, through white and with
// the specified alpha.
func Radial(colors int, start, end Hue, alpha float64) DivergingPalette {
	p := make(divergingPalette, colors)
	h := colors / 2
	c := HSVA{V: 1, A: alpha}
	ds := 0.5 / float64(h)
	for i := range p[:h] {
		c.H = float64(start)
		c.S = 0.5 - float64(i)*ds
		p[i] = color.NRGBAModel.Convert(c)
		c.H = float64(end)
		p[len(p)-1-i] = color.NRGBAModel.Convert(c)
	}
	if colors%2 != 0 {
		p[colors/2] = color.NRGBA{0xff, 0xff, 0xff, byte(math.MaxUint8 * alpha)}
	}

	return p
}
