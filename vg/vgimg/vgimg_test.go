// Copyright ©2012 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgimg_test

import (
	"bytes"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"reflect"
	"sync"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func TestIssue179(t *testing.T) {
	scatter, err := plotter.NewScatter(plotter.XYs{{1, 1}, {0, 1}, {0, 0}})
	if err != nil {
		log.Fatal(err)
	}
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.Add(scatter)
	p.HideAxes()

	c := vgimg.JpegCanvas{Canvas: vgimg.New(5.08*vg.Centimeter, 5.08*vg.Centimeter)}
	p.Draw(draw.New(c))
	b := bytes.NewBuffer([]byte{})
	if _, err = c.WriteTo(b); err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("testdata/issue179_golden.jpg")
	if err != nil {
		t.Fatal(err)
	}

	ok, err := cmpimg.Equal("jpg", b.Bytes(), want)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		ioutil.WriteFile("testdata/issue179.jpg", b.Bytes(), 0644)
		t.Fatalf("images differ")
	}
}

func TestConcurrentInit(t *testing.T) {
	ft, err := vg.MakeFont("Helvetica", 10)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
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
	p, err := plot.New()
	if err != nil {
		t.Fatal(err)
	}

	xys := plotter.XYs{
		plotter.XY{0, 0},
		plotter.XY{1, 1},
		plotter.XY{2, 2},
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

	err = p.Save(100, 100, "testdata/issue540.png")
	if err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("testdata/issue540_golden.png")
	if err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadFile("testdata/issue540.png")
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
