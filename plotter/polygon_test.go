// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"image/color"
	"log"
	"math"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/recorder"
)

// rings implements the XYers interface.
type rings []plotter.XYs

func (r rings) Len() int        { return len(r) }
func (r rings) LenAt(i int) int { return len(r[i]) }
func (r rings) XY(i, j int) (x, y float64) {
	p := r[i][j]
	return p.X, p.Y
}

// ExamplePolygon_holes draws a polygon with holes, showing how
// the different built-in vg backends render polygons with holes.
// The output of this example is at
// https://github.com/gonum/plot/blob/master/plotter/testdata/polygon_holes_golden.png,
// https://github.com/gonum/plot/blob/master/plotter/testdata/polygon_holes_golden.svg,
// https://github.com/gonum/plot/blob/master/plotter/testdata/polygon_holes_golden.pdf, and
// https://github.com/gonum/plot/blob/master/plotter/testdata/polygon_holes_golden.eps.
func ExamplePolygon_holes() {
	// Create an outer ring.
	outer1 := plotter.XYs{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}}

	// Create an inner ring with the same
	// winding order as the outer ring.
	inner1 := plotter.XYs{{X: 0.5, Y: 0.5}, {X: 1.5, Y: 0.5}, {X: 1.5, Y: 1.5}, {X: 0.5, Y: 1.5}}

	// Create an inner polygon with the opposite
	// winding order as the outer polygon.
	inner2 := plotter.XYs{{X: 3.5, Y: 2.5}, {X: 2.5, Y: 2.5}, {X: 2.5, Y: 3.5}, {X: 3.5, Y: 3.5}}

	poly, err := plotter.NewPolygon(rings{outer1, inner1, inner2})
	if err != nil {
		log.Panic(err)
	}
	poly.Color = color.NRGBA{B: 255, A: 255}

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Polygon with holes"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(poly)
	p.Legend.Add("key", poly)
	p.Legend.Font.Size = vg.Points(8)
	p.Legend.TextStyle.Color = color.White
	p.Legend.ThumbnailWidth = vg.Points(10)

	// Here we save the image in different file formats
	// to show how each back end handles polygon holes.

	// The vgimg backend treats both internal polygons
	// as holes by default.
	err = p.Save(100, 100, "testdata/polygon_holes.png")
	if err != nil {
		log.Panic(err)
	}

	// The vgsvg, vgpdf, and vgeps backgrounds all treat
	// the internal polygon with the opposite winding
	// direction as a hole but do not consider the internal
	// polygon with the same winding direction to be a hole.
	err = p.Save(100, 100, "testdata/polygon_holes.svg")
	if err != nil {
		log.Panic(err)
	}
	err = p.Save(100, 100, "testdata/polygon_holes.pdf")
	if err != nil {
		log.Panic(err)
	}
	err = p.Save(100, 100, "testdata/polygon_holes.eps")
	if err != nil {
		log.Panic(err)
	}
}

func TestPolygon_holes(t *testing.T) {
	cmpimg.CheckPlot(ExamplePolygon_holes, t, "polygon_holes.png", "polygon_holes.svg", "polygon_holes.pdf", "polygon_holes.eps")
}

// ExamplePolygon_hexagons creates a heat map with hexagon shapes.
// The output of this example is at
// https://github.com/gonum/plot/blob/master/plotter/testdata/polygon_hexagons_golden.png.
func ExamplePolygon_hexagons() {
	// hex returns a hexagon centered at (x,y) with radius r.
	hex := func(x, y, r float64) rings {
		g := make(plotter.XYs, 6)
		for i := range g {
			g[i].X = x + r*math.Cos(math.Pi/3*float64(i))
			g[i].Y = y + r*math.Sin(math.Pi/3*float64(i))
		}
		return rings{g}
	}

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Hexagons"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	colorMap := moreland.SmoothBlueRed()
	colorMap.SetMax(2)
	colorMap.SetMin(-2)
	colorMap.SetConvergePoint(0)

	const (
		r = math.Pi / 4 // r is the hexagon radius.
		// x0 and y0 are the beginning coordinates for the hexagon plot.
		x0 = 0.0
		y0 = 0.0
		// nx and ny are the number of hexagons in the x and y directions.
		nx = 5
		ny = 5
	)
	// dx and dy are the distance between hexgons.
	dx := 3 * r
	dy := r * math.Sqrt(3)

	xstart := []float64{x0, x0 - 1.5*r}
	ystart := []float64{y0, y0 - r}
	for i, xmin := range xstart {
		ymin := ystart[i]
		x := xmin
		for ix := 0; ix < nx; ix++ {
			y := ymin
			for iy := 0; iy < ny; iy++ {
				var poly *plotter.Polygon
				poly, err = plotter.NewPolygon(hex(x, y, r))
				if err != nil {
					log.Panic(err)
				}
				poly.Color, err = colorMap.At(math.Sin(x) + math.Sin(y))
				if err != nil {
					log.Panic(err)
				}
				poly.LineStyle.Width = 0
				p.Add(poly)
				y += dy
			}
			x += dx
		}
	}
	if err = p.Save(100, 100, "testdata/polygon_hexagons.png"); err != nil {
		log.Panic(err)
	}
}

func TestPolygon_hexagons(t *testing.T) {
	cmpimg.CheckPlot(ExamplePolygon_hexagons, t, "polygon_hexagons.png")
}

// This test ensures that the plotter doesn't panic if there are
// polygons wholly outside of the plotting range.
func TestPolygon_clip(t *testing.T) {
	poly, err := plotter.NewPolygon(
		rings{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}},
	)
	if err != nil {
		t.Fatal(err)
	}
	poly.Color = color.Black // Give the polygon a color to fill.
	p, err := plot.New()
	if err != nil {
		t.Fatal(err)
	}

	// Set the plotting range so that the polygon is outside of it.
	p.X.Min = 2
	p.X.Max = 5

	p.Add(poly)
	c := new(recorder.Canvas)
	dc := draw.NewCanvas(c, vg.Centimeter, vg.Centimeter)
	p.Draw(dc) // If this does not panic, then the test passes.
}
