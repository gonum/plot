// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moreland

import (
	"log"
	"math"
	"os"
	"testing"

	"github.com/gonum/matrix/mat64"
	"github.com/gonum/plot"
	"github.com/gonum/plot/internal/cmpimg"
	"github.com/gonum/plot/palette"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgimg"
)

type offsetUnitGrid struct {
	XOffset, YOffset float64

	Data *mat64.Dense
}

func (g offsetUnitGrid) Dims() (c, r int)   { r, c = g.Data.Dims(); return c, r }
func (g offsetUnitGrid) Z(c, r int) float64 { return g.Data.At(r, c) }
func (g offsetUnitGrid) X(c int) float64 {
	_, n := g.Data.Dims()
	if c < 0 || n <= c {
		panic("index out of range")
	}
	return float64(c) + g.XOffset
}
func (g offsetUnitGrid) Y(r int) float64 {
	m, _ := g.Data.Dims()
	if r < 0 || m <= r {
		panic("index out of range")
	}
	return float64(r) + g.YOffset
}

// This Example gives examples of plots using the palettes in this package.
// The output can be found at
// https://github.com/gonum/plot/blob/master/palette/moreland/testdata/moreland_golden.png.
func Example() {
	m := offsetUnitGrid{
		XOffset: -50,
		YOffset: -50,
		Data:    mat64.NewDense(100, 100, nil),
	}
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			x := float64(i-50) / 10
			y := float64(j-50) / 10
			v := math.Sin(x*x+y*y) / (x*x + y*y)
			m.Data.Set(i, j, v)
		}
	}

	const (
		rows = 3
		cols = 3
	)
	c := vgimg.New(vg.Points(800), vg.Points(800))
	dc := draw.New(c)
	tiles := draw.Tiles{
		Rows: rows,
		Cols: cols,
	}
	type paletteHolder struct {
		name string
		cmap palette.Palette
	}
	palettes := []paletteHolder{
		paletteHolder{
			name: "SmoothBlueRed",
			cmap: SmoothBlueRed().Palette(255),
		},
		paletteHolder{
			name: "SmoothBlueTan",
			cmap: SmoothBlueTan().Palette(255),
		},
		paletteHolder{
			name: "SmoothGreenPurple",
			cmap: SmoothGreenPurple().Palette(255),
		},
		paletteHolder{
			name: "SmoothGreenRed",
			cmap: SmoothGreenRed().Palette(255),
		},
		paletteHolder{
			name: "SmoothPurpleOrange",
			cmap: SmoothPurpleOrange().Palette(255),
		},
		paletteHolder{
			name: "BlackBody",
			cmap: BlackBody().Palette(255),
		},
		paletteHolder{
			name: "ExtendedBlackBody",
			cmap: ExtendedBlackBody().Palette(255),
		},
		paletteHolder{
			name: "Kindlmann",
			cmap: Kindlmann().Palette(255),
		},
		paletteHolder{
			name: "ExtendedKindlmann",
			cmap: ExtendedKindlmann().Palette(255),
		},
	}

	for i, plte := range palettes {

		h := plotter.NewHeatMap(m, plte.cmap)

		p, err := plot.New()
		if err != nil {
			log.Panic(err)
		}
		p.Title.Text = plte.name

		p.Add(h)

		p.X.Padding = 0
		p.Y.Padding = 0
		p.Draw(tiles.At(dc, i%cols, i/cols))
	}

	pngimg := vgimg.PngCanvas{Canvas: c}
	f, err := os.Create("testdata/moreland.png")
	if err != nil {
		log.Panic(err)
	}
	if _, err = pngimg.WriteTo(f); err != nil {
		log.Panic(err)
	}
}

func TestHeatMap(t *testing.T) {
	cmpimg.CheckPlot(Example, t, "moreland.png")
}
