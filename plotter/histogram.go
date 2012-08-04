// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
	"math"
	"image/color"
	"fmt"
)

// Histogram implements the Plotter interface,
// drawing a histogram of the data.
type Histogram struct {
	XYer

	// BinWidth is the width of each histogram
	// bin.  If BinWidth is non-positive then a
	// reasonable default is used.
	BinWidth float64

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
		XYer: xy,
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
			{ trX(bin.xMin), trY(0) },
			{ trX(bin.xMax), trY(0) },
			{ trX(bin.xMax), trY(bin.height) },
			{ trX(bin.xMin), trY(bin.height) },
		}
		if l.FillColor != nil {
			da.FillPolygon(l.FillColor, da.ClipPolygonXY(pts))
		}
		pts = append(pts, plot.Point{ trX(bin.xMin), trY(0) })
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
	xmin, xmax := xDataRange(h.XYer)
	w := h.UsedBinWidth()
	n := int(math.Ceil((xmax - xmin) / w))
	bins := make([]histBin, n)

	for i := range bins {
		bins[i].xMin = xmin + float64(i)*w
		bins[i].xMax = xmin + float64(i+1)*w
	}

	for i := 0; i < h.Len(); i++ {
		x := h.X(i)
		bin := int((x - xmin) / w)
		if bin >= n && x == xmax {
			bin = n-1
		}
		if bin < 0 || bin >= n || bins[bin].xMin > x || bins[bin].xMax < x {
			panic(fmt.Sprintf("%g, xmin=%g, xmax=%g, w=%g, bin=%d, n=%d\n",
				h.X(i), xmin, xmax, w, bin, n))
		}
		bins[bin].height += h.Y(i)
	}
	return bins
}

// histBin is a histogram bin.
type histBin struct {
	xMin, xMax float64
	height float64
}

// UsedBinWidth returns the bin width being used to
// draw the histogram.  If the histogram's BinWidth
// field is positive then it is returned, otherwise the
// default value is computed and returned.
//
// The default number of bins is is the square root of
// the number of samples.  According to Wikipedia,
// this is what MS Excel uses.
func (h *Histogram) UsedBinWidth() float64 {
	if h.BinWidth > 0 {
		return h.BinWidth
	}
	n := 0.0
	for i := 0; i < h.Len(); i++ {
		n += math.Max(h.Y(i), 1.0)
	}
	xmin, xmax := xDataRange(h.XYer)
	return (xmax - xmin) / math.Sqrt(n)
}