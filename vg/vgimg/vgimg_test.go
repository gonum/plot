// Copyright ©2012 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgimg_test

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/font/liberation"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

var cache = font.NewCache(liberation.Collection())

func TestIssue179(t *testing.T) {
	scatter, err := plotter.NewScatter(plotter.XYs{
		{X: 1, Y: 1}, {X: 0, Y: 1}, {X: 0, Y: 0},
	})
	if err != nil {
		log.Fatal(err)
	}
	p := plot.New()
	p.Add(scatter)
	p.HideAxes()

	c := vgimg.JpegCanvas{Canvas: vgimg.New(5.08*vg.Centimeter, 5.08*vg.Centimeter)}
	p.Draw(draw.New(c))
	b := bytes.NewBuffer([]byte{})
	if _, err = c.WriteTo(b); err != nil {
		t.Fatal(err)
	}

	want, err := os.ReadFile("testdata/issue179_golden.jpg")
	if err != nil {
		t.Fatal(err)
	}

	ok, err := cmpimg.Equal("jpg", b.Bytes(), want)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		_ = os.WriteFile("testdata/issue179.jpg", b.Bytes(), 0644)
		t.Fatalf("images differ")
	}
}

func TestConcurrentInit(t *testing.T) {
	var (
		ft = cache.Lookup(font.Font{Variant: "Sans"}, 10)
		wg sync.WaitGroup
	)
	wg.Add(2)
	go func() {
		c := vgimg.New(215, 215)
		c.FillString(ft, vg.Point{}, "hi")
		wg.Done()
	}()
	go func() {
		c := vgimg.New(215, 215)
		c.FillString(ft, vg.Point{}, "hi")
		wg.Done()
	}()
	wg.Wait()
}

func TestUseBackgroundColor(t *testing.T) {
	colors := []color.Color{color.Transparent, color.NRGBA{R: 255, A: 255}}
	for i, col := range colors {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			c := vgimg.NewWith(vgimg.UseWH(1, 1), vgimg.UseBackgroundColor(col))
			img := c.Image()
			wantCol := color.RGBAModel.Convert(col)
			haveCol := img.At(0, 0)
			if !reflect.DeepEqual(haveCol, wantCol) {
				t.Fatalf("color should be %#v but is %#v", wantCol, haveCol)
			}
		})
	}
}

func TestIssue540(t *testing.T) {
	p := plot.New()

	xys := plotter.XYs{
		plotter.XY{X: 0, Y: 0},
		plotter.XY{X: 1, Y: 1},
		plotter.XY{X: 2, Y: 2},
	}

	p.Title.Text = "My title"
	p.X.Tick.Label.Font.Size = 0 // hide X-axis labels
	p.Y.Tick.Label.Font.Size = 0 // hide Y-axis labels

	lines, points, err := plotter.NewLinePoints(xys)
	if err != nil {
		log.Fatal(err)
	}
	lines.Color = color.RGBA{B: 255, A: 255}

	p.Add(lines, points)
	p.Add(plotter.NewGrid())

	if *cmpimg.GenerateTestData {
		// Recreate Golden images and exit.
		err = p.Save(100, 100, "testdata/issue540_golden.png")
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	err = p.Save(100, 100, "testdata/issue540.png")
	if err != nil {
		t.Fatal(err)
	}

	want, err := os.ReadFile("testdata/issue540_golden.png")
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile("testdata/issue540.png")
	if err != nil {
		t.Fatal(err)
	}

	ok, err := cmpimg.Equal("png", got, want)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("images differ")
	}
}

func TestIssue687(t *testing.T) {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	const (
		fname = "testdata/issue687.png"
		size  = 500
	)
	cmpimg.CheckPlot(func() {
		p := plot.New()
		p.Title.Text = "Issue 687"
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		data := make(plotter.XYs, 4779) // Lower values still lose part of horizontal.
		for i := range data {
			data[i] = plotter.XY{X: float64(i) / float64(len(data)*2), Y: float64(len(data) - min(i, len(data)/2))}
		}
		lines, err := plotter.NewLine(data)
		lines.Color = color.RGBA{R: 0xff, A: 0xff}
		if err != nil {
			t.Fatal(err)
		}
		p.Add(lines)

		err = p.Save(size, size, fname)
		if err != nil {
			t.Fatal(err)
		}

	}, t, filepath.Base(fname))
}
