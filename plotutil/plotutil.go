// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

// plotutil contains a small number of utilites for creating plots.
package plotutil

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/vg"
	"fmt"
	"image/color"
)

// DefaultColors is a set of colors used by the Color funciton.
var DefaultColors = []color.Color{
	color.RGBA{A: 255, R: 190},
	color.RGBA{A: 255, G: 190},
	color.RGBA{A: 255, B: 190},
	color.RGBA{A: 255, R: 190, B: 190},
	color.RGBA{A: 255, G: 190, B: 190},
	color.RGBA{A: 255, R: 250, G: 190, B: 10},
}

// Color returns the ith default color, wrapping
// if i is less than zero or greater than the max
// number of colors in the DefaultColors slice.
func Color(i int) color.Color {
	n := len(DefaultColors)
	if i < 0 {
		return DefaultColors[i%n+n]
	}
	return DefaultColors[i%n]
}

// DefaultGlyphShapes is a set of GlyphDrawers used by
// the Shape function.
var DefaultGlyphShapes = []plot.GlyphDrawer{
	plot.RingGlyph{},
	plot.SquareGlyph{},
	plot.TriangleGlyph{},
	plot.CrossGlyph{},
	plot.PlusGlyph{},
	plot.CircleGlyph{},
	plot.BoxGlyph{},
	plot.PyramidGlyph{},
}

// Shape returns the ith default glyph shape,
// wrapping if i is less than zero or greater
// than the max number of GlyphDrawers
// in the DefaultGlyphShapes slice.
func Shape(i int) plot.GlyphDrawer {
	n := len(DefaultGlyphShapes)
	if i < 0 {
		return DefaultGlyphShapes[i%n+n]
	}
	return DefaultGlyphShapes[i%n]
}

// DefaultDashes is a set of dash patterns used by
// the Dashes function.
var DefaultDashes = [][]vg.Length{
	{},

	{vg.Points(6), vg.Points(2)},

	{vg.Points(2), vg.Points(2)},

	{vg.Points(1), vg.Points(1)},

	{vg.Points(5), vg.Points(2), vg.Points(1), vg.Points(2)},

	{vg.Points(10), vg.Points(2), vg.Points(2), vg.Points(2),
		vg.Points(2), vg.Points(2), vg.Points(2), vg.Points(2)},

	{vg.Points(10), vg.Points(2), vg.Points(2), vg.Points(2)},

	{vg.Points(5), vg.Points(2), vg.Points(5), vg.Points(2),
		vg.Points(2), vg.Points(2), vg.Points(2), vg.Points(2)},

	{vg.Points(4), vg.Points(2), vg.Points(4), vg.Points(1),
		vg.Points(1), vg.Points(1), vg.Points(1), vg.Points(1),
		vg.Points(1), vg.Points(1)},
}

// Dashes returns the ith default dash pattern,
// wrapping if i is less than zero or greater
// than the max number of dash patters
// in the DefaultDashes slice.
func Dashes(i int) []vg.Length {
	n := len(DefaultDashes)
	if i < 0 {
		return DefaultDashes[i%n+n]
	}
	return DefaultDashes[i%n]
}

// AddBoxPlots adds box plot plotters to a plot and
// sets the X axis of the plot to be nominal.
// The variadic arguments must be either strings
// or plotter.Valuers.  Each valuer adds a box plot
// to the plot at the X location corresponding to
// the number of box plots added before it.  If a
// plotter.Valuer is immediately preceeded by a
// string then the string value is used to label the
// tick mark for the box plot's X location.
func AddBoxPlots(plt *plot.Plot, width vg.Length, vs ...interface{}) {
	var names []string
	name := ""
	for _, v := range vs {
		switch t := v.(type) {
		case string:
			name = t

		case plotter.Valuer:
			plt.Add(plotter.NewBoxPlot(width, float64(len(names)), t))
			names = append(names, name)
			name = ""

		default:
			panic(fmt.Sprintf("AddScatters handles strings and plotter.XYers, got %T", t))
		}
	}
	plt.NominalX(names...)
}

// AddScatters adds Scatter plotters to a plot.
// The variadic arguments must be either strings
// or plotter.XYers.  Each plotter.XYer is added to
// the plot using the next color, and glyph shape
// via the Color and Shape functions. If a
// plotter.XYer is immediately preceeded by
// a string then a legend entry is added to the plot
// using the string as the name.
func AddScatters(plt *plot.Plot, vs ...interface{}) {
	name := ""
	var i int
	for _, v := range vs {
		switch t := v.(type) {
		case string:
			name = t

		case plotter.XYer:
			s := plotter.NewScatter(t)
			s.Color = Color(i)
			s.Shape = Shape(i)
			i++
			plt.Add(s)
			if name != "" {
				plt.Legend.Add(name, s)
				name = ""
			}

		default:
			panic(fmt.Sprintf("AddScatters handles strings and plotter.XYers, got %T", t))
		}
	}
}

// AddLines adds Line plotters to a plot.
// The variadic arguments must be either strings
// or plotter.XYers.  Each plotter.XYer is added to
// the plot using the next color and dashes
// shape via the Color and Dashes functions.
// If a plotter.XYer is immediately preceeded by
// a string then a legend entry is added to the plot
// using the string as the name.
func AddLines(plt *plot.Plot, vs ...interface{}) {
	name := ""
	var i int
	for _, v := range vs {
		switch t := v.(type) {
		case string:
			name = t

		case plotter.XYer:
			l := plotter.NewLine(t)
			l.Color = Color(i)
			l.Dashes = Dashes(i)
			i++
			plt.Add(l)
			if name != "" {
				plt.Legend.Add(name, l)
				name = ""
			}

		default:
			panic(fmt.Sprintf("AddLines handles strings and plotter.XYers, got %T", t))
		}
	}
}

// AddLinePoints adds Line and Scatter plotters to a
// plot.  The variadic arguments must be either strings
// or plotter.XYers.  Each plotter.XYer is added to
// the plot using the next color, dashes, and glyph
// shape via the Color, Dashes, and Shape functions.
// If a plotter.XYer is immediately preceeded by
// a string then a legend entry is added to the plot
// using the string as the name.
func AddLinePoints(plt *plot.Plot, vs ...interface{}) {
	name := ""
	var i int
	for _, v := range vs {
		switch t := v.(type) {
		case string:
			name = t

		case plotter.XYer:
			l, s := plotter.NewLinePoints(t)
			l.Color = Color(i)
			l.Dashes = Dashes(i)
			s.Color = Color(i)
			s.Shape = Shape(i)
			i++
			plt.Add(l, s)
			if name != "" {
				plt.Legend.Add(name, l, s)
				name = ""
			}

		default:
			panic(fmt.Sprintf("AddLinePoints handles strings and plotter.XYers, got %T", t))
		}
	}
}
