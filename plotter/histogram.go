// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
	"fmt"
	"image/color"
	"math"
)

// Histogram implements the Plotter interface,
// drawing a histogram of the data.
type Histogram struct {
	XYer

	// If Normalize is positive then the mass under
	// the histogram is normalized to sum to the
	// given value.
	Normalize float64

	// NumBins is the number of bins.
	// If NumBins is non-positive then a reasonable
	// default is used.
	// 
	// The default number of bins is is the square root of
	// the number of samples.  According to Wikipedia,
	// this is what MS Excel uses.
	NumBins int

	// FillColor is the color used to fill each
	// bar of the histogram.  If the color is nil
	// then the bars are not filled.
	FillColor color.Color

	// LineStyle is the style of the outline of each
	// bar of the histogram.
	plot.LineStyle
}

// NewHistogram returns a new histogram that represents
// the distribution of the given points.
func NewHistogram(xy XYer) *Histogram {
	return &Histogram{
		XYer:      xy,
		FillColor: color.Gray{128},
		LineStyle: DefaultLineStyle,
	}
}

// Plot implements the Plotter interface, drawing a line
// that connects each point in the Line.
func (l *Histogram) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)

	for _, bin := range l.bins() {
		pts := []plot.Point{
			{trX(bin.xMin), trY(0)},
			{trX(bin.xMax), trY(0)},
			{trX(bin.xMax), trY(bin.height)},
			{trX(bin.xMin), trY(bin.height)},
		}
		if l.FillColor != nil {
			da.FillPolygon(l.FillColor, da.ClipPolygonXY(pts))
		}
		pts = append(pts, plot.Point{trX(bin.xMin), trY(0)})
		da.StrokeLines(l.LineStyle, da.ClipLinesXY(pts)...)
	}
}

// DataRange returns the minimum and maximum X and Y values
func (l *Histogram) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = math.Inf(1)
	xmax = math.Inf(-1)
	ymax = math.Inf(-1)
	for _, bin := range l.bins() {
		if bin.xMax > xmax {
			xmax = bin.xMax
		}
		if bin.xMin < xmin {
			xmin = bin.xMin
		}
		if bin.height > ymax {
			ymax = bin.height
		}
	}
	return
}

// bins returns the histogram's bins.
func (h *Histogram) bins() []histBin {
	xmin, xmax := Range(XValues{h.XYer})
	n := h.numBins(xmin, xmax)
	bins := make([]histBin, n)

	w := (xmax - xmin) / float64(n-1)
	for i := range bins {
		bins[i].xMin = xmin + float64(i)*w
		bins[i].xMax = xmin + float64(i+1)*w
	}

	sum := 0.0
	for i := 0; i < h.Len(); i++ {
		x := h.X(i)
		bin := int((x - xmin) / w)
		if bin < 0 || bin >= n {
			panic(fmt.Sprintf("%g, xmin=%g, xmax=%g, w=%g, bin=%d, n=%d\n",
				h.X(i), xmin, xmax, w, bin, n))
		}
		y := h.Y(i)
		bins[bin].height += y
		sum += y
	}

	if h.Normalize > 0 {
		for i := range bins {
			bins[i].height = (bins[i].height / w / sum) * h.Normalize
		}
	}

	return bins
}

// numBins returns the actual number of bins
// used by the histogram.
func (h *Histogram) numBins(xmin, xmax float64) int {
	n := h.NumBins
	if h.NumBins <= 0 {
		m := 0.0
		for i := 0; i < h.Len(); i++ {
			m += math.Max(h.Y(i), 1.0)
		}
		n = int(math.Ceil(math.Sqrt(m)))
	}
	if n < 1 || xmax <= xmin {
		n = 1
	}
	return n
}

// histBin is a histogram bin.
type histBin struct {
	xMin, xMax float64
	height     float64
}
