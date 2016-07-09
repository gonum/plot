// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package palette

import (
	"image"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// ColorMapLegend is a plot.Plotter that draws a color bar legend for a ColorMap.
type ColorMapLegend struct {
	// Vertical determines wether the legend will be
	// plotted vertically or horizontally.
	// The default is false (horizontal).
	Vertical bool

	// NColors specifies the number of colors to be
	// shown in the legend. The default is 255.
	NColors int

	cm ColorMap
}

// NewColorMapLegend creates a new legend plotter.
func NewColorMapLegend(cm ColorMap) *ColorMapLegend {
	return &ColorMapLegend{
		cm:      cm,
		NColors: 255,
	}
}

// Plot implements the Plot method of the plot.Plotter interface.
func (l *ColorMapLegend) Plot(c draw.Canvas, p *plot.Plot) {
	if l.cm.Max() == l.cm.Min() {
		panic("palette: ColorMap Max==Min")
	}
	var img *image.NRGBA64
	var xmin, xmax, ymin, ymax vg.Length
	if l.Vertical {
		trX, trY := p.Transforms(&c)
		xmin = trX(l.cm.Min())
		ymin = trY(0)
		xmax = trX(l.cm.Max())
		ymax = trY(1)
		img = image.NewNRGBA64(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: 1, Y: l.NColors},
		})
		for i := 0; i < l.NColors; i++ {
			color, err := l.cm.At(float64(i) / float64(l.NColors-1))
			if err != nil {
				panic(err)
			}
			if l.Vertical {
				img.Set(0, l.NColors-1-i, color)
			} else {
				img.Set(0, i, color)
			}
		}
	} else {
		trX, trY := p.Transforms(&c)
		ymin = trY(l.cm.Min())
		xmin = trX(0)
		ymax = trY(l.cm.Max())
		xmax = trX(1)
		img = image.NewNRGBA64(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: l.NColors, Y: 1},
		})
		for i := 0; i < l.NColors; i++ {
			color, err := l.cm.At(float64(i) / float64(l.NColors-1))
			if err != nil {
				panic(err)
			}
			img.Set(i, 0, color)
		}
	}
	rect := vg.Rectangle{
		Min: vg.Point{X: xmin, Y: ymin},
		Max: vg.Point{X: xmax, Y: ymax},
	}
	c.DrawImage(rect, img)
}

// DataRange implements the DataRange method
// of the plot.DataRanger interface.
func (l *ColorMapLegend) DataRange() (xmin, xmax, ymin, ymax float64) {
	if l.cm.Max() == l.cm.Min() {
		panic("palette: ColorMap Max==Min")
	}
	if l.Vertical {
		return 0, 1, l.cm.Min(), l.cm.Max()
	}
	return l.cm.Min(), l.cm.Max(), 0, 1
}

// SetupPlot changes the default settings of p so that
// they are appropriate for plotting a color bar legend.
func (l *ColorMapLegend) SetupPlot(p *plot.Plot) {
	if l.Vertical {
		p.HideX()
		p.Y.Padding = 0
	} else {
		p.HideY()
		p.X.Padding = 0
	}
}

// GlyphBoxes implements the GlyphBoxes method
// of the plot.GlyphBoxer interface.
func (l *ColorMapLegend) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	return nil
}
