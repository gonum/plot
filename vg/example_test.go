// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vg_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/golang/freetype/truetype"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func Example_addFont() {
	// download font from debian
	const url = "http://http.debian.net/debian/pool/main/f/fonts-ipafont/fonts-ipafont_00303.orig.tar.gz"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("could not download IPA font file: %+v", err)
	}
	defer resp.Body.Close()

	ttf, err := untargz("IPAfont00303/ipam.ttf", resp.Body)
	if err != nil {
		log.Fatalf("could not untar archive: %+v", err)
	}

	fontTTF, err := truetype.Parse(ttf)
	if err != nil {
		log.Fatal(err)
	}
	const fontName = "Mincho"
	vg.AddFont(fontName, fontTTF)

	plot.DefaultFont = fontName
	plotter.DefaultFont = fontName

	p, err := plot.New()
	if err != nil {
		log.Fatalf("could not create plot: %+v", err)
	}
	p.Title.Text = "Hello, 世界"
	p.X.Label.Text = "世"
	p.Y.Label.Text = "界"

	labels, err := plotter.NewLabels(
		plotter.XYLabels{
			XYs:    make(plotter.XYs, 1),
			Labels: []string{"こんにちは世界"},
		},
	)
	if err != nil {
		log.Fatalf("could not create labels: %+v", err)
	}
	p.Add(labels)
	p.Add(plotter.NewGrid())

	err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, "mincho-font.png")
	if err != nil {
		log.Fatalf("could not save plot: %+v", err)
	}
}

func untargz(name string, r io.Reader) ([]byte, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not create gzip reader: %v", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("could not find %q in tar archive", name)
			}
			return nil, fmt.Errorf("could not extract header from tar archive: %v", err)
		}

		if hdr == nil || hdr.Name != name {
			continue
		}

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, tr)
		if err != nil {
			return nil, fmt.Errorf("could not extract %q file from tar archive: %v", name, err)
		}
		return buf.Bytes(), nil
	}
}

func Example_inMemoryCanvas() {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("could not create plot: %+v", err)
	}

	p.Title.Text = "sin(x)"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "f(x)"

	p.X.Min = -2 * math.Pi
	p.X.Max = +2 * math.Pi

	fct := plotter.NewFunction(func(x float64) float64 {
		return math.Sin(x)
	})
	fct.Color = color.RGBA{R: 255, A: 255}

	p.Add(fct, plotter.NewGrid())

	c := vgimg.New(10*vg.Centimeter, 5*vg.Centimeter)
	p.Draw(draw.New(c))

	// Save image.
	f, err := os.Create("testdata/sine.png")
	if err != nil {
		log.Fatalf("could not create output image file: %+v", err)
	}
	defer f.Close()

	err = png.Encode(f, c.Image())
	if err != nil {
		log.Fatalf("could not encode image to PNG: %+v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could not close output image file: %+v", err)
	}
}

func Example_writerToCanvas() {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("could not create plot: %+v", err)
	}

	p.Title.Text = "cos(x)"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "f(x)"

	p.X.Min = -2 * math.Pi
	p.X.Max = +2 * math.Pi

	fct := plotter.NewFunction(func(x float64) float64 {
		return math.Cos(x)
	})
	fct.Color = color.RGBA{B: 255, A: 255}

	p.Add(fct, plotter.NewGrid())

	c := vgimg.PngCanvas{
		Canvas: vgimg.New(10*vg.Centimeter, 5*vg.Centimeter),
	}
	p.Draw(draw.New(c))

	// Save image.
	f, err := os.Create("testdata/cosine.png")
	if err != nil {
		log.Fatalf("could not create output image file: %+v", err)
	}
	defer f.Close()

	_, err = c.WriteTo(f)
	if err != nil {
		log.Fatalf("could not encode image to PNG: %+v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could not close output image file: %+v", err)
	}
}
