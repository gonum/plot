// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"image/color"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

var (
	// DefaultGridLineStyle is the default style for grid lines.
	DefaultGridLineStyle = draw.LineStyle{
		Color: color.Gray{128},
		Width: vg.Points(0.25),
	}
)

// Grid implements the plot.Plotter interface, drawing
// a set of grid lines at the major tick marks.
type Grid struct {
	// Vertical is the style of the vertical lines.
	Vertical draw.LineStyle

	// Horizontal is the style of the horizontal lines.
	Horizontal draw.LineStyle
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
func (g *Grid) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)

	if g.Vertical.Color == nil {
		goto horiz
	}
	for _, tk := range plt.X.Tick.Marker.Ticks(plt.X.Min, plt.X.Max) {
		if tk.IsMinor() {
			continue
		}
		x := trX(tk.Value)
		c.StrokeLine2(g.Vertical, x, c.Min.Y, x, c.Min.Y+c.Size().Y)
	}

horiz:
	if g.Horizontal.Color == nil {
		return
	}
	for _, tk := range plt.Y.Tick.Marker.Ticks(plt.Y.Min, plt.Y.Max) {
		if tk.IsMinor() {
			continue
		}
		y := trY(tk.Value)
		c.StrokeLine2(g.Horizontal, c.Min.X, y, c.Min.X+c.Size().X, y)
	}
}
