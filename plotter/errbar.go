// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
	"fmt"
	"math"
)

// ErrorBars implement the Plotter interface, drawing
// either vertical, horizontal, or both vertical and
// horizontal error bars.
type ErrorBars struct {
	// XYer gives the center location of the error
	// bar.  If the XYer also implements the XErrorer
	// interface then a horizontal error bar is plotted.
	// If the XYer also implements YErrorer then
	// a vertical error bar is drawn.
	XYer

	// LineStyle is the style of the errorbar and cap lines.
	plot.LineStyle

	// CapWidth is the width of the error bar caps.
	CapWidth vg.Length
}

// MakeErrorBars returns ErrorBars using the default
// line style and cap width.  If the XYer doesn't implement
// either XErrorer or YErrorer then a non-nil error
// is returned, and while the resulting ErrorBars can be
// successfully added to a plot, they will not draw
// anything.
func MakeErrorBars(xy XYer) (ErrorBars, error) {
	_, xerr := xy.(XErrorer)
	_, yerr := xy.(YErrorer)
	if !xerr && !yerr {
		err := fmt.Errorf("Type %T doesn't implement XErrorer or YErrorer", xy)
		return ErrorBars{}, err
	}
	return ErrorBars{
		XYer:      xy,
		CapWidth:  vg.Points(6),
		LineStyle: DefaultLineStyle,
	}, nil
}

// Plot implements the Plotter interface, drawing
// error bars in either the X or Y directions or both.
func (e ErrorBars) Plot(da plot.DrawArea, p *plot.Plot) {
	e.plotVerticalBars(&da, p)
	e.plotHorizontalBars(&da, p)
}

// plotVerticalBars plots the vertical error bars
// if this ErrorBars implements the YErrorer interface.
func (e ErrorBars) plotVerticalBars(da *plot.DrawArea, p *plot.Plot) {
	yerr, ok := e.XYer.(YErrorer)
	if !ok {
		return
	}
	trX, trY := p.Transforms(da)
	for i := 0; i < e.Len(); i++ {
		errlow, errhigh := yerr.YError(i)
		y := e.Y(i)
		min, max := trY(y+errlow), trY(y+errhigh)
		x := trX(e.X(i))
		da.StrokeLines(e.LineStyle, da.ClipLinesXY([]plot.Point{{x, min}, {x, max}})...)
		e.plotVerticalCap(da, plot.Point{x, min})
		e.plotVerticalCap(da, plot.Point{x, max})
	}
}

// plotVerticalCap plots a horizontal line, centered
// at the given point, capping a vertical errorbar.
func (e ErrorBars) plotVerticalCap(da *plot.DrawArea, pt plot.Point) {
	w := e.CapWidth / 2
	if da.Contains(pt) {
		da.StrokeLine2(e.LineStyle, pt.X-w, pt.Y, pt.X+w, pt.Y)
	}
}

// plotHorizontalBars plots the horizontal error bars
// if this ErrorBars implements the XErrorer interface.
func (e ErrorBars) plotHorizontalBars(da *plot.DrawArea, p *plot.Plot) {
	xerr, ok := e.XYer.(XErrorer)
	if !ok {
		return
	}
	trX, trY := p.Transforms(da)
	for i := 0; i < e.Len(); i++ {
		errlow, errhigh := xerr.XError(i)
		x := e.X(i)
		min, max := trX(x+errlow), trX(x+errhigh)
		y := trY(e.Y(i))
		da.StrokeLines(e.LineStyle, da.ClipLinesXY([]plot.Point{{min, y}, {max, y}})...)
		e.plotHorizontalCap(da, plot.Point{min, y})
		e.plotHorizontalCap(da, plot.Point{max, y})
	}
}

// plotHorizontalCap plots a vertical line, centered
// at the given point, capping a horizontal errorbar.
func (e ErrorBars) plotHorizontalCap(da *plot.DrawArea, pt plot.Point) {
	w := e.CapWidth / 2
	if da.Contains(pt) {
		da.StrokeLine2(e.LineStyle, pt.X, pt.Y-w, pt.X, pt.Y+w)
	}
}

// GlyphBoxes implements the GlyphBoxer interface,
// ensuring that the caps of the error bars are not
// clipped by the edge of the plot.
func (e ErrorBars) GlyphBoxes(p *plot.Plot) (boxes []plot.GlyphBox) {
	boxes = append(boxes, e.verticalGlyphBoxes(p)...)
	boxes = append(boxes, e.horizontalGlyphBoxes(p)...)
	return
}

// verticalGlyphBoxes returns the GlyphBoxes
// for the vertical error bar caps.
func (e ErrorBars) verticalGlyphBoxes(p *plot.Plot) (boxes []plot.GlyphBox) {
	yerr, ok := e.XYer.(YErrorer)
	if !ok {
		return
	}
	vertRect := plot.Rect{
		Min:  plot.Point{X: -e.CapWidth / 2},
		Size: plot.Point{X: e.CapWidth},
	}
	for i := 0; i < e.Len(); i++ {
		x, y := e.X(i), e.Y(i)
		errlow, errhigh := yerr.YError(i)
		min, max := p.Y.Norm(y+errlow), p.Y.Norm(y+errhigh)
		boxes = append(boxes,
			plot.GlyphBox{X: p.X.Norm(x), Y: min, Rect: vertRect},
			plot.GlyphBox{X: p.X.Norm(x), Y: max, Rect: vertRect})
	}
	return
}

// horizontalGlyphBoxes returns the GlyphBoxes
// for the horizontal error bar caps.
func (e ErrorBars) horizontalGlyphBoxes(p *plot.Plot) (boxes []plot.GlyphBox) {
	xerr, ok := e.XYer.(XErrorer)
	if !ok {
		return
	}
	horzRect := plot.Rect{
		Min:  plot.Point{Y: -e.CapWidth / 2},
		Size: plot.Point{Y: e.CapWidth},
	}
	for i := 0; i < e.Len(); i++ {
		x, y := e.X(i), e.Y(i)
		errlow, errhigh := xerr.XError(i)
		min, max := p.X.Norm(x+errlow), p.X.Norm(x+errhigh)
		boxes = append(boxes,
			plot.GlyphBox{X: min, Y: p.Y.Norm(y), Rect: horzRect},
			plot.GlyphBox{X: max, Y: p.Y.Norm(y), Rect: horzRect})
	}
	return

}

// DataRange implements the DataRanger interface,
// returning the minimum and maximum X and Y
// values of the error bars.
func (e ErrorBars) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = e.xDataRange()
	ymin, ymax = e.yDataRange()
	return
}

// xDataRange returns the range of x values
// for the error bars.
func (e ErrorBars) xDataRange() (xmin, xmax float64) {
	xmin, xmax = Range(XValues{e})
	xerr, ok := e.XYer.(XErrorer)
	if !ok {
		return
	}
	for i := 0; i < e.Len(); i++ {
		x := e.X(i)
		errlow, errhigh := xerr.XError(i)
		xmin = math.Min(xmin, x+errlow)
		xmax = math.Max(xmax, x+errlow)
		xmin = math.Min(xmin, x+errhigh)
		xmax = math.Max(xmax, x+errhigh)
	}
	return
}

// yDataRange returns the range of y values
// for the error bars.
func (e ErrorBars) yDataRange() (ymin, ymax float64) {
	ymin, ymax = Range(YValues{e})
	yerr, ok := e.XYer.(YErrorer)
	if !ok {
		return
	}
	for i := 0; i < e.Len(); i++ {
		y := e.Y(i)
		errlow, errhigh := yerr.YError(i)
		ymin = math.Min(ymin, y+errlow)
		ymax = math.Max(ymax, y+errlow)
		ymin = math.Min(ymin, y+errhigh)
		ymax = math.Max(ymax, y+errhigh)
	}
	return
}
