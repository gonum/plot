// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
	"fmt"
)

// Labels implements the Plotter interface, drawing
// a set of labels on the plot.
type Labels struct {
	// XYer is a set of coordinates that also implements
	// the Labeller interface.
	XYer

	// TextStyle gives the style of the labels.
	plot.TextStyle

	// XAlign and YAlign are multiplied by the width
	// and height of each label respectively and the
	// added to the final location.  E.g., XAlign=-0.5
	// and YAlign=-0.5 centers the label at the given
	// X, Y location, and XAlign=0, YAlign=0 aligns
	// the text to the left of the point, and XAlign=-1,
	// YAlign=0 aligns the text to the right of the point.
	XAlign, YAlign float64

	// XOffs and YOffs are added directly to the final
	// label X and Y location respectively.
	XOffs, YOffs vg.Length
}

// Labels returns a Labels using the default TextStyle,
// with the labels left-aligned above the corresponding
// X, Y point.
func MakeLabels(xy XYer) (Labels, error) {
	fnt, err := vg.MakeFont(DefaultFont, vg.Points(10))
	if err != nil {
		return Labels{}, err
	}
	if _, ok := xy.(Labeller); !ok {
		err := fmt.Errorf("Type %T doesn't implement Labeller", xy)
		return Labels{}, err
	}
	return Labels{ XYer: xy, TextStyle: plot.TextStyle{Font: fnt} }, nil
}

// Plot implements the Plotter interface for Labels.
func (l Labels) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)
	labels := l.XYer.(Labeller)
	for i := 0; i < l.Len(); i++ {
		x, y := trX(l.X(i))+l.XOffs, trY(l.Y(i))+l.YOffs
		if da.Contains(plot.Point{x, y}) {
			label := labels.Label(i)
			da.FillText(l.TextStyle, x, y, l.XAlign, l.YAlign, label)
		}
	}
}

// DataRange returns the minimum and maximum X and Y values
func (l Labels) DataRange() (xmin, xmax, ymin, ymax float64) {
	return XYRange(l)
}

// GlyphBoxes returns a slice of GlyphBoxes, one for
// each of the labels.
func (l Labels) GlyphBoxes(p *plot.Plot) (boxes []plot.GlyphBox) {
	labels := l.XYer.(Labeller)
	for i := 0; i < l.Len(); i++ {
		x, y := p.X.Norm(l.X(i)), p.Y.Norm(l.Y(i))
		txt := labels.Label(i)
		rect := l.Rect(txt)
		rect.Min.X += l.Width(txt)*vg.Length(l.XAlign) + l.XOffs
		rect.Min.Y += l.Height(txt)*vg.Length(l.YAlign) + l.YOffs
		box := plot.GlyphBox{X: x, Y: y, Rect: rect}
		boxes = append(boxes, box)
	}
	return
}
