// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
	"math"
)

// DefaultCapWidth is the default width of error bar caps.
var DefaultCapWidth = vg.Points(5)

// YErrorBars implements the plot.Plotter, plot.DataRanger,
// and plot.GlyphBoxer interfaces, drawing vertical error
// bars, denoting error in Y values.
type YErrorBars struct {
	XYs

	// YErrors is a copy of the Y errors for each point.
	YErrors

	// LineStyle is the style used to draw the error bars.
	plot.LineStyle

	// CapWidth is the width of the caps drawn at the top
	// of each error bar.
	CapWidth vg.Length
}

func NewYErrorBars(yerrs interface {
	XYer
	YErrorer
}) (*YErrorBars, error) {

	errors := make(YErrors, yerrs.Len())
	for i := range errors {
		errors[i].Low, errors[i].High = yerrs.YError(i)
		if err := CheckFloats(errors[i].Low, errors[i].High); err != nil {
			return nil, err
		}
	}
	xys, err := CopyXYs(yerrs)
	if err != nil {
		return nil, err
	}

	return &YErrorBars{
		XYs:       xys,
		YErrors:   errors,
		LineStyle: DefaultLineStyle,
		CapWidth:  DefaultCapWidth,
	}, nil
}

// Plot implements the Plotter interface, drawing labels.
func (e *YErrorBars) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)
	for i, err := range e.YErrors {
		x := trX(e.XYs[i].X)
		ylow := trY(e.XYs[i].Y - err.Low)
		yhigh := trY(e.XYs[i].Y + err.High)

		bar := da.ClipLinesY([]plot.Point{{x, ylow}, {x, yhigh}})
		da.StrokeLines(e.LineStyle, bar...)
		e.drawCap(&da, x, ylow)
		e.drawCap(&da, x, yhigh)
	}
}

// drawCap draws the cap if it is not clipped.
func (e *YErrorBars) drawCap(da *plot.DrawArea, x, y vg.Length) {
	if !da.Contains(plot.Pt(x, y)) {
		return
	}
	da.StrokeLine2(e.LineStyle, x-e.CapWidth/2, y, x+e.CapWidth/2, y)
}

// DataRange implements the plot.DataRanger interface.
func (e *YErrorBars) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = Range(XValues{e})
	ymin = math.Inf(1)
	ymax = math.Inf(-1)
	for i, err := range e.YErrors {
		y := e.XYs[i].Y
		ylow := y - err.Low
		yhigh := y + err.High
		ymin = math.Min(math.Min(math.Min(ymin, y), ylow), yhigh)
		ymax = math.Max(math.Max(math.Max(ymax, y), ylow), yhigh)
	}
	return
}

// GlyphBoxes implements the plot.GlyphBoxer interface.
func (e *YErrorBars) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	rect := plot.Rect{
		Min: plot.Point{
			X: -e.CapWidth / 2,
			Y: -e.LineStyle.Width / 2,
		},
		Size: plot.Point{
			X: e.CapWidth,
			Y: e.LineStyle.Width,
		},
	}
	var bs []plot.GlyphBox
	for i, err := range e.YErrors {
		x := plt.X.Norm(e.XYs[i].X)
		y := e.XYs[i].Y
		bs = append(bs,
			plot.GlyphBox{X: x, Y: plt.Y.Norm(y - err.Low), Rect: rect},
			plot.GlyphBox{X: x, Y: plt.Y.Norm(y + err.High), Rect: rect})
	}
	return bs
}

// XErrorBars implements the plot.Plotter, plot.DataRanger,
// and plot.GlyphBoxer interfaces, drawing horizontal error
// bars, denoting error in Y values.
type XErrorBars struct {
	XYs

	// XErrors is a copy of the X errors for each point.
	XErrors

	// LineStyle is the style used to draw the error bars.
	plot.LineStyle

	// CapWidth is the width of the caps drawn at the top
	// of each error bar.
	CapWidth vg.Length
}

func NewXErrorBars(xerrs interface {
	XYer
	XErrorer
}) (*XErrorBars, error) {

	errors := make(XErrors, xerrs.Len())
	for i := range errors {
		errors[i].Low, errors[i].High = xerrs.XError(i)
		if err := CheckFloats(errors[i].Low, errors[i].High); err != nil {
			return nil, err
		}
	}
	xys, err := CopyXYs(xerrs)
	if err != nil {
		return nil, err
	}

	return &XErrorBars{
		XYs:       xys,
		XErrors:   errors,
		LineStyle: DefaultLineStyle,
		CapWidth:  DefaultCapWidth,
	}, nil
}

// Plot implements the Plotter interface, drawing labels.
func (e *XErrorBars) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)
	for i, err := range e.XErrors {
		y := trY(e.XYs[i].Y)
		xlow := trX(e.XYs[i].X - err.Low)
		xhigh := trX(e.XYs[i].X + err.High)

		bar := da.ClipLinesX([]plot.Point{{xlow, y}, {xhigh, y}})
		da.StrokeLines(e.LineStyle, bar...)
		e.drawCap(&da, xlow, y)
		e.drawCap(&da, xhigh, y)
	}
}

// drawCap draws the cap if it is not clipped.
func (e *XErrorBars) drawCap(da *plot.DrawArea, x, y vg.Length) {
	if !da.Contains(plot.Pt(x, y)) {
		return
	}
	da.StrokeLine2(e.LineStyle, x, y-e.CapWidth/2, x, y+e.CapWidth/2)
}

// DataRange implements the plot.DataRanger interface.
func (e *XErrorBars) DataRange() (xmin, xmax, ymin, ymax float64) {
	ymin, ymax = Range(YValues{e})
	xmin = math.Inf(1)
	xmax = math.Inf(-1)
	for i, err := range e.XErrors {
		x := e.XYs[i].X
		xlow := x - err.Low
		xhigh := x + err.High
		xmin = math.Min(math.Min(math.Min(xmin, x), xlow), xhigh)
		xmax = math.Max(math.Max(math.Max(xmax, x), xlow), xhigh)
	}
	return
}

// GlyphBoxes implements the plot.GlyphBoxer interface.
func (e *XErrorBars) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	rect := plot.Rect{
		Min: plot.Point{
			X: -e.LineStyle.Width / 2,
			Y: -e.CapWidth / 2,
		},
		Size: plot.Point{
			X: e.LineStyle.Width,
			Y: e.CapWidth,
		},
	}
	var bs []plot.GlyphBox
	for i, err := range e.XErrors {
		x := e.XYs[i].X
		y := plt.Y.Norm(e.XYs[i].Y)
		bs = append(bs,
			plot.GlyphBox{X: plt.X.Norm(x - err.Low), Y: y, Rect: rect},
			plot.GlyphBox{X: plt.X.Norm(x + err.High), Y: y, Rect: rect})
	}
	return bs
}
