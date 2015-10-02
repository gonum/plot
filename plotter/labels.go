// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"errors"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

var (
	// DefaultFont is the default font for label text.
	DefaultFont = plot.DefaultFont

	// DefaultFontSize is the default font.
	DefaultFontSize = vg.Points(10)
)

// Labels implements the Plotter interface,
// drawing a set of labels at specified points.
type Labels struct {
	XYs

	// Labels is the set of labels corresponding
	// to each point.
	Labels []string

	// TextStyle is the style of the label text.
	draw.TextStyle

	// XAlign and YAlign are multiplied by the width
	// and height of each label respectively and the
	// added to the final location.  E.g., XAlign=-0.5
	// and YAlign=-0.5 centers the label at the given
	// X, Y location, and XAlign=0, YAlign=0 aligns
	// the text to the left of the point, and XAlign=-1,
	// YAlign=0 aligns the text to the right of the point.
	XAlign, YAlign float64

	// XOffset and YOffset are added directly to the final
	// label X and Y location respectively.
	XOffset, YOffset vg.Length
}

// NewLabels returns a new Labels using the DefaultFont and
// the DefaultFontSize.
func NewLabels(d XYLabeller) (*Labels, error) {
	xys, err := CopyXYs(d)
	if err != nil {
		return nil, err
	}

	if d.Len() != len(xys) {
		return nil, errors.New("Number of points does not match the number of labels")
	}

	strs := make([]string, d.Len())
	for i := range strs {
		strs[i] = d.Label(i)
	}

	fnt, err := vg.MakeFont(DefaultFont, DefaultFontSize)
	if err != nil {
		return nil, err
	}

	return &Labels{
		XYs:       xys,
		Labels:    strs,
		TextStyle: draw.TextStyle{Font: fnt},
	}, nil
}

// Plot implements the Plotter interface, drawing labels.
func (l *Labels) Plot(c draw.Canvas, p *plot.Plot) {
	trX, trY := p.Transforms(&c)
	for i, label := range l.Labels {
		x := trX(l.XYs[i].X)
		y := trY(l.XYs[i].Y)
		if !c.Contains(draw.Point{X: x, Y: y}) {
			continue
		}
		x += l.XOffset
		y += l.YOffset
		c.FillText(l.TextStyle, x, y, l.XAlign, l.YAlign, label)
	}
}

// DataRange returns the minimum and maximum X and Y values
func (l *Labels) DataRange() (xmin, xmax, ymin, ymax float64) {
	return XYRange(l)
}

// GlyphBoxes returns a slice of GlyphBoxes,
// one for each of the labels, implementing the
// plot.GlyphBoxer interface.
func (l *Labels) GlyphBoxes(p *plot.Plot) []plot.GlyphBox {
	bs := make([]plot.GlyphBox, len(l.Labels))
	for i, label := range l.Labels {
		bs[i].X = p.X.Norm(l.XYs[i].X)
		bs[i].Y = p.Y.Norm(l.XYs[i].Y)
		w := l.Width(label)
		h := l.Height(label)
		bs[i].Rectangle.Min.X = w*vg.Length(l.XAlign) + l.XOffset
		bs[i].Rectangle.Min.Y = h*vg.Length(l.YAlign) + l.YOffset
		bs[i].Rectangle.Max.X = w + w*vg.Length(l.XAlign) + l.XOffset
		bs[i].Rectangle.Max.Y = h + h*vg.Length(l.YAlign) + l.YOffset
	}
	return bs
}

// XYLabeller combines the XYer and Labeller types.
type XYLabeller interface {
	XYer
	Labeller
}

// XYLabels holds XY data with labels.
// The ith label corresponds to the ith XY.
type XYLabels struct {
	XYs
	Labels []string
}

// Label returns the label for point index i.
func (l XYLabels) Label(i int) string {
	return l.Labels[i]
}
