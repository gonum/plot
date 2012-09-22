// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
	"image/color"
)

var (
	// DefaultGridLineStyle is the default style for grid lines.
	DefaultGridLineStyle = plot.LineStyle{
		Color: color.Gray{128},
		Width: vg.Points(0.25),
	}
)

// Grid implements the plot.Plotter interface, drawing
// a set of grid lines at the major tick marks.
type Grid struct {
	// Vertical is the style of the vertical lines.
	Vertical plot.LineStyle

	// Horizontal is the style of the horizontal lines.
	Horizontal plot.LineStyle
}

// NewGrid returns a new grid with both vertical and
// horizontal lines using the default grid line style.
func NewGrid() *Grid {
	return &Grid{
		Vertical:   DefaultGridLineStyle,
		Horizontal: DefaultGridLineStyle,
	}
}

// Plot implements the plot.Plotter interface.
func (g *Grid) Plot(da plot.DrawArea, plt *plot.Plot) {
	trX, trY := plt.Transforms(&da)

	if g.Vertical.Color == nil {
		goto horiz
	}
	for _, tk := range plt.X.Tick.Marker(plt.X.Min, plt.X.Max) {
		if tk.IsMinor() {
			continue
		}
		x := trX(tk.Value)
		da.StrokeLine2(g.Vertical, x, da.Min.Y, x, da.Min.Y+da.Size.Y)
	}

horiz:
	if g.Horizontal.Color == nil {
		return
	}
	for _, tk := range plt.Y.Tick.Marker(plt.Y.Min, plt.Y.Max) {
		if tk.IsMinor() {
			continue
		}
		y := trY(tk.Value)
		da.StrokeLine2(g.Horizontal, da.Min.X, y, da.Min.X+da.Size.X, y)
	}
}
