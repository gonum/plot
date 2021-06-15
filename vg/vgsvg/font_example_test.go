// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgsvg_test

import (
	"log"
	"os"

	lmit "github.com/go-fonts/latin-modern/lmroman10italic"
	lreg "github.com/go-fonts/liberation/liberationserifregular"
	xfnt "golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
)

func Example_embedFonts() {
	// Use Latin-Modern fonts.
	cmi10 := font.Font{Typeface: "Latin-Modern", Style: xfnt.StyleItalic}
	fnt, err := opentype.Parse(lmit.TTF)
	if err != nil {
		log.Fatalf("could not parse Latin-Modern fonts: %+v", err)
	}

	font.DefaultCache.Add([]font.Face{{
		Font: cmi10,
		Face: fnt,
	}})
	plot.DefaultFont = cmi10

	p := plot.New()
	p.Title.Text = "Scatter plot"
	p.X.Label.Text = "x-Axis"
	p.Y.Label.Text = "y-Axis"

	scatter, err := plotter.NewScatter(plotter.XYs{{X: 1, Y: 1}, {X: 0, Y: 1}, {X: 0, Y: 0}})
	if err != nil {
		log.Fatalf("could not create scatter: %v", err)
	}
	p.Add(scatter)

	c := vgsvg.NewWith(
		vgsvg.UseWH(5*vg.Centimeter, 5*vg.Centimeter),
		vgsvg.EmbedFonts(true),
	)
	p.Draw(draw.New(c))

	f, err := os.Create("testdata/embed_fonts.svg")
	if err != nil {
		log.Fatalf("could not create output SVG file: %+v", err)
	}
	defer f.Close()

	_, err = c.WriteTo(f)
	if err != nil {
		log.Fatalf("could not write output SVG plot: %+v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could not close output SVG file: %v", err)
	}
}

func Example_standardFonts() {
	// Use standard fonts.
	tms := font.Font{Typeface: "Times"}
	fnt, err := opentype.Parse(lreg.TTF)
	if err != nil {
		log.Fatalf("could not parse Times fonts: %+v", err)
	}

	font.DefaultCache.Add([]font.Face{{
		Font: tms,
		Face: fnt,
	}})
	plot.DefaultFont = tms

	p := plot.New()
	p.Title.Text = "Scatter plot"
	p.X.Label.Text = "x-Axis"
	p.Y.Label.Text = "y-Axis"

	scatter, err := plotter.NewScatter(plotter.XYs{{X: 1, Y: 1}, {X: 0, Y: 1}, {X: 0, Y: 0}})
	if err != nil {
		log.Fatalf("could not create scatter: %v", err)
	}
	p.Add(scatter)

	err = p.Save(5*vg.Centimeter, 5*vg.Centimeter, "testdata/standard_fonts.svg")
	if err != nil {
		log.Fatalf("could not save SVG plot: %+v", err)
	}
}
