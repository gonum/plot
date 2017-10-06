// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
)

// ScatterColor implements the Plotter interface, drawing
// a scatter plot of x, y, z triples where the z value
// determines the colour of a scatter plot point.
type ScatterColor struct {
	XYZs

	// Palette is the color palette used to render
	// the scatter plot. Palette must not be nil or
	// return a zero length []color.Color.
	Palette palette.Palette

	// Min and Max define the dynamic range of the
	// scatter plot.
	Min, Max float64

	// GlyphStyle is the style of the glyphs drawn
	// at each point.
	draw.GlyphStyle
}

// NewScatterColor creates a new scatter plot plotter for
// the given data, with a minimum and maximum
// colour range.
func NewScatterColor(xyz XYZer, p palette.Palette) (*ScatterColor, error) {
	cpy, err := CopyXYZs(xyz)

	if err != nil {
		return nil, err
	}

	minz := cpy[0].Z
	maxz := cpy[0].Z

	for _, d := range cpy {
		minz = math.Min(minz, d.Z)
		maxz = math.Max(maxz, d.Z)
	}
	return &ScatterColor{
		XYZs:    cpy,
		Palette: p,
		Min:     minz,
		Max:     maxz,
	}, nil
}

// Plot implements the Plot method of the plot.Plotter interface.
func (sc *ScatterColor) Plot(c draw.Canvas, plt *plot.Plot) {

	pal := sc.Palette.Colors()
	if len(pal) == 0 {
		panic("scattercolor: empty palette")
	}
	// ps scales the palette uniformly across the data range.
	ps := float64(len(pal)-1) / (sc.Max - sc.Min)

	trX, trY := plt.Transforms(&c)

	var p vg.Path

	var col color.Color

	for _, d := range sc.XYZs {
		//Transform the data x,y coordinate of this scatter plot element
		//to the corresponding drawing coordinate.
		col = pal[int((d.Z-sc.Min)*ps+0.5)] // Apply palette scaling.
		GlyphNewStyle := draw.GlyphStyle{
			Color:  col,
			Radius: vg.Points(3),
			Shape:  draw.CircleGlyph{},
		}
		c.DrawGlyph(GlyphNewStyle, vg.Point{X: trX(d.X), Y: trY(d.Y)})

		if col != nil {
			c.SetColor(col)
			c.Fill(p)
		}
	}
}

// DataRange implements the DataRange method
// of the plot.DataRanger interface.
func (sc *ScatterColor) DataRange() (xmin, xmax, ymin, ymax float64) {
	return XYRange(XYValues{sc.XYZs})
}

// GlyphBoxes implements the GlyphBoxes method
// of the plot.GlyphBoxer interface.
func (bs *ScatterColor) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, len(bs.XYZs))
	for i, d := range bs.XYZs {
		boxes[i].X = plt.X.Norm(d.X)
		boxes[i].Y = plt.Y.Norm(d.Y)
		boxes[i].Rectangle = vg.Rectangle{
			Min: vg.Point{X: -3, Y: -3}, //3 is a radius of glyphs
			Max: vg.Point{X: +3, Y: +3},
		}
	}
	return boxes
}
