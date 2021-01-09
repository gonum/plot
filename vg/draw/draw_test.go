// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw_test

import (
	"fmt"
	"image/color"
	"reflect"
	"testing"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/recorder"

	_ "gonum.org/v1/plot/vg/vgeps"
	_ "gonum.org/v1/plot/vg/vgimg"
	_ "gonum.org/v1/plot/vg/vgpdf"
	_ "gonum.org/v1/plot/vg/vgsvg"
	_ "gonum.org/v1/plot/vg/vgtex"
)

func TestCrop(t *testing.T) {
	ls := draw.LineStyle{
		Color: color.NRGBA{0, 20, 0, 123},
		Width: 0.1 * vg.Inch,
	}
	var r1 recorder.Canvas
	c1 := draw.NewCanvas(&r1, 6, 3)
	c11 := draw.Crop(c1, 0, -3, 0, 0)
	c12 := draw.Crop(c1, 3, 0, 0, 0)

	var r2 recorder.Canvas
	c2 := draw.NewCanvas(&r2, 6, 3)
	c21 := draw.Canvas{
		Canvas: c2.Canvas,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: 0},
			Max: vg.Point{X: 3, Y: 3},
		},
	}
	c22 := draw.Canvas{
		Canvas: c2.Canvas,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 3, Y: 0},
			Max: vg.Point{X: 6, Y: 3},
		},
	}
	str := "unexpected result: %+v != %+v"
	if c11.Rectangle != c21.Rectangle {
		t.Errorf(str, c11.Rectangle, c21.Rectangle)
	}
	if c12.Rectangle != c22.Rectangle {
		t.Errorf(str, c11.Rectangle, c21.Rectangle)
	}

	c11.StrokeLine2(ls, c11.Min.X, c11.Min.Y,
		c11.Min.X+3*vg.Inch, c11.Min.Y+3*vg.Inch)
	c12.StrokeLine2(ls, c12.Min.X, c12.Min.Y,
		c12.Min.X+3*vg.Inch, c12.Min.Y+3*vg.Inch)
	c21.StrokeLine2(ls, c21.Min.X, c21.Min.Y,
		c21.Min.X+3*vg.Inch, c21.Min.Y+3*vg.Inch)
	c22.StrokeLine2(ls, c22.Min.X, c22.Min.Y,
		c22.Min.X+3*vg.Inch, c22.Min.Y+3*vg.Inch)

	if !reflect.DeepEqual(r1.Actions, r2.Actions) {
		t.Errorf(str, r1.Actions, r2.Actions)
	}
}

func TestTile(t *testing.T) {
	var r recorder.Canvas
	c := draw.NewCanvas(&r, 13, 7)
	const (
		rows = 2
		cols = 3
		pad  = 1
	)
	tiles := draw.Tiles{
		Rows: rows, Cols: cols,
		PadTop: pad, PadBottom: pad,
		PadRight: pad, PadLeft: pad,
		PadX: pad, PadY: pad,
	}
	rectangles := [][]vg.Rectangle{
		{
			vg.Rectangle{
				Min: vg.Point{X: 1, Y: 4},
				Max: vg.Point{X: 4, Y: 6},
			},
			vg.Rectangle{
				Min: vg.Point{X: 5, Y: 4},
				Max: vg.Point{X: 8, Y: 6},
			},
			vg.Rectangle{
				Min: vg.Point{X: 9, Y: 4},
				Max: vg.Point{X: 12, Y: 6},
			},
		},
		{
			vg.Rectangle{
				Min: vg.Point{X: 1, Y: 1},
				Max: vg.Point{X: 4, Y: 3},
			},
			vg.Rectangle{
				Min: vg.Point{X: 5, Y: 1},
				Max: vg.Point{X: 8, Y: 3},
			},
			vg.Rectangle{
				Min: vg.Point{X: 9, Y: 1},
				Max: vg.Point{X: 12, Y: 3},
			},
		},
	}
	for j := 0; j < rows; j++ {
		for i := 0; i < cols; i++ {
			str := "row %d col %d unexpected result: %+v != %+v"
			tile := tiles.At(c, i, j)
			if tile.Rectangle != rectangles[j][i] {
				t.Errorf(str, j, i, tile.Rectangle, rectangles[j][i])
			}
		}
	}
}

func TestFormattedCanvas(t *testing.T) {
	for _, test := range []struct {
		format string
		err    error
	}{
		{format: "eps"},
		{format: "jpg"},
		{format: "jpeg"},
		{format: "pdf"},
		{format: "png"},
		{format: "svg"},
		{format: "tex"},
		{format: "tiff"},
		{format: "tif"},
		{
			format: "",
			err:    fmt.Errorf("unsupported format: \"\""),
		},
		{
			format: "abc",
			err:    fmt.Errorf("unsupported format: \"abc\""),
		},
	} {
		t.Run(test.format, func(t *testing.T) {
			_, err := draw.NewFormattedCanvas(10, 10, test.format)
			switch {
			case err != nil && test.err != nil:
				if got, want := err.Error(), test.err.Error(); got != want {
					t.Fatalf("invalid error.\ngot= %v\nwant=%v", got, want)
				}
			case err != nil && test.err == nil:
				t.Fatalf("unexpected error: %+v", err)
			case err == nil && test.err != nil:
				t.Fatalf("expected an error (got=%v)", err)
			}
		})
	}
}
