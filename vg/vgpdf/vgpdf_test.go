// Copyright ©2018 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgpdf_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/internal/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgpdf"
)

// ExampleEmbedFonts shows how one can embed (or not) fonts inside
// a PDF plot.
func ExampleEmbedFonts() {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("could not create plot: %v", err)
	}

	pts := plotter.XYs{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
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

func TestEmbedFonts(t *testing.T) {
	for _, tc := range []struct {
		name  string
		embed bool
	}{
		{
			name:  "testdata/disable-embedded-fonts_golden.pdf",
			embed: false,
		},
		{
			name:  "testdata/enable-embedded-fonts_golden.pdf",
			embed: true,
		},
	} {
		t.Run(fmt.Sprintf("embed=%v", tc.embed), func(t *testing.T) {
			p, err := plot.New()
			if err != nil {
				t.Fatalf("could not create plot: %v", err)
			}

			pts := plotter.XYs{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
			line, err := plotter.NewLine(pts)
			if err != nil {
				t.Fatalf("could not create line: %v", err)
			}
			p.Add(line)
			p.X.Label.Text = "X axis"
			p.Y.Label.Text = "Y axis"

			c := vgpdf.New(100, 100)

			// enable/disable embedding fonts
			c.EmbedFonts(tc.embed)
			p.Draw(draw.New(c))

			var buf bytes.Buffer
			_, err = c.WriteTo(&buf)
			if err != nil {
				t.Fatalf("could not write canvas: %v", err)
			}

			want, err := ioutil.ReadFile(tc.name)
			if err != nil {
				t.Fatalf("failed to read golden plot: %v", err)
			}

			ok, err := cmpimg.Equal("pdf", buf.Bytes(), want)
			if err != nil {
				t.Fatalf("failed to run cmpimg test: %v", err)
			}

			if !ok {
				t.Fatalf("plot mismatch: %v", tc.name)
			}
		})
	}
}

func TestArc(t *testing.T) {
	pts := plotter.XYs{{1, 1}, {2, 2}}
	scat, err := plotter.NewScatter(pts)
	if err != nil {
		t.Fatal(err)
	}
	p, err := plot.New()
	if err != nil {
		t.Fatal(err)
	}
	p.Add(scat)

	c := vgpdf.New(100, 100)

	c.EmbedFonts(false)
	p.Draw(draw.New(c))

	f, err := os.Create("testdata/arc.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = c.WriteTo(f)
	if err != nil {
		t.Fatalf("could not write canvas: %v", err)
	}

	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("testdata/arc_golden.pdf")
	if err != nil {
		t.Fatalf("failed to read golden plot: %v", err)
	}

	got, err := ioutil.ReadFile("testdata/arc.pdf")
	if err != nil {
		t.Fatalf("failed to read plot: %v", err)
	}

	ok, err := cmpimg.Equal("pdf", got, want)
	if err != nil {
		t.Fatalf("failed to run cmpimg test: %v", err)
	}

	if !ok {
		t.Fatalf("plot mismatch")
	}
}
