// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgsvg_test

import (
	"bytes"
	"os"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
)

func TestSVG(t *testing.T) {
	cmpimg.CheckPlot(Example, t, "scatter.svg")
}

func TestEmbedFonts(t *testing.T) {
	fnt := plot.DefaultFont
	defer func() {
		plot.DefaultFont = fnt
	}()

	cmpimg.CheckPlot(Example_embedFonts, t, "embed_fonts.svg")
}

func TestStandardFonts(t *testing.T) {
	fnt := plot.DefaultFont
	defer func() {
		plot.DefaultFont = fnt
	}()

	cmpimg.CheckPlot(Example_standardFonts, t, "standard_fonts.svg")
}

func TestNewWith(t *testing.T) {
	p := plot.New()
	p.Title.Text = "Scatter plot"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	scatter, err := plotter.NewScatter(plotter.XYs{{X: 1, Y: 1}, {X: 0, Y: 1}, {X: 0, Y: 0}})
	if err != nil {
		t.Fatalf("could not create scatter: %v", err)
	}
	p.Add(scatter)

	c := vgsvg.NewWith(vgsvg.UseWH(5*vg.Centimeter, 5*vg.Centimeter))
	p.Draw(draw.New(c))

	b := new(bytes.Buffer)
	if _, err = c.WriteTo(b); err != nil {
		t.Fatal(err)
	}

	want, err := os.ReadFile("testdata/scatter_golden.svg")
	if err != nil {
		t.Fatal(err)
	}

	ok, err := cmpimg.Equal("svg", b.Bytes(), want)
	if err != nil {
		t.Fatalf("could not compare images: %v", err)
	}
	if !ok {
		t.Fatalf("images differ:\ngot:\n%s\nwant:\n%s\n", b.Bytes(), want)
	}
}

func TestHtmlEscape(t *testing.T) {
	p := plot.New()
	p.Title.Text = "Scatter & line plot"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	scatter, err := plotter.NewScatter(plotter.XYs{{X: 1, Y: 1}, {X: 0, Y: 1}, {X: 0, Y: 0}})
	if err != nil {
		t.Fatalf("could not create scatter: %v", err)
	}
	p.Add(scatter)

	line, err := plotter.NewLine(plotter.XYs{{X: 1, Y: 1}, {X: 0, Y: 1}, {X: 0, Y: 0}})
	if err != nil {
		t.Fatalf("could not create scatter: %v", err)
	}
	line.Width = 0.5
	p.Add(line)

	c := vgsvg.NewWith(vgsvg.UseWH(5*vg.Centimeter, 5*vg.Centimeter))
	p.Draw(draw.New(c))

	b := new(bytes.Buffer)
	if _, err = c.WriteTo(b); err != nil {
		t.Fatal(err)
	}

	if *cmpimg.GenerateTestData {
		// Recreate Golden images and exit.
		err = os.WriteFile("testdata/scatter_line_golden.svg", b.Bytes(), 0o644)
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	want, err := os.ReadFile("testdata/scatter_line_golden.svg")
	if err != nil {
		t.Fatal(err)
	}

	ok, err := cmpimg.Equal("svg", b.Bytes(), want)
	if err != nil {
		t.Fatalf("could not compare images: %v", err)
	}
	if !ok {
		t.Fatalf("images differ:\ngot:\n%s\nwant:\n%s\n", b.Bytes(), want)
	}
}
