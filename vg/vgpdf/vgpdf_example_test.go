// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgpdf_test

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgpdf"
)

// Example_embedFonts shows how one can embed (or not) fonts inside
// a PDF plot.
func Example_embedFonts() {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("could not create plot: %v", err)
	}

	pts := plotter.XYs{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}}
	line, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatalf("could not create line: %v", err)
	}
	p.Add(line)
	p.X.Label.Text = "X axis"
	p.Y.Label.Text = "Y axis"

	c := vgpdf.New(100, 100)

	// enable/disable embedding fonts
	c.EmbedFonts(true)
	p.Draw(draw.New(c))

	f, err := os.Create("testdata/enable-embedded-fonts.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = c.WriteTo(f)
	if err != nil {
		log.Fatalf("could not write canvas: %v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could not save canvas: %v", err)
	}
}

// Example_multipage shows how one can create a PDF with multiple pages.
func Example_multipage() {
	c := vgpdf.New(5*vg.Centimeter, 5*vg.Centimeter)

	for i, col := range []color.RGBA{{B: 255, A: 255}, {R: 255, A: 255}} {
		if i > 0 {
			// Add a new page.
			c.NextPage()
		}

		p, err := plot.New()
		if err != nil {
			log.Fatalf("could not create plot: %v", err)
		}

		pts := plotter.XYs{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}}
		line, err := plotter.NewLine(pts)
		if err != nil {
			log.Fatalf("could not create line: %v", err)
		}
		line.Color = col
		p.Add(line)
		p.Title.Text = fmt.Sprintf("Plot %d", i+1)
		p.X.Label.Text = "X axis"
		p.Y.Label.Text = "Y axis"

		// Write plot to page.
		p.Draw(draw.New(c))
	}

	f, err := os.Create("testdata/multipage.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = c.WriteTo(f)
	if err != nil {
		log.Fatalf("could not write canvas: %v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could not save canvas: %v", err)
	}
}
