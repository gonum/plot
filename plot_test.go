// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot_test

import (
	"bytes"
	"fmt"
	"image/color"
	"io/ioutil"
	"math"
	"testing"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/recorder"
	"gonum.org/v1/plot/vg/vgimg"
)

func TestLegendAlignment(t *testing.T) {
	font, err := vg.MakeFont(plot.DefaultFont, 10)
	if err != nil {
		t.Fatalf("failed to create font: %v", err)
	}
	l := plot.Legend{
		ThumbnailWidth: vg.Points(20),
		TextStyle: draw.TextStyle{
			Font:    font,
			Handler: plot.DefaultTextHandler,
		},
	}
	for i, n := range []string{"A", "B", "C", "D"} {
		b, err := plotter.NewBarChart(plotter.Values{0}, 1)
		if err != nil {
			t.Fatalf("failed to create bar chart %q: %v", n, err)
		}
		b.Color = color.Gray{byte(i+1)*64 - 1}
		l.Add(n, b)
	}

	c := vgimg.PngCanvas{Canvas: vgimg.New(5*vg.Centimeter, 5*vg.Centimeter)}
	l.Draw(draw.New(c))
	var buf bytes.Buffer
	if _, err = c.WriteTo(&buf); err != nil {
		t.Fatal(err)
	}

	if *cmpimg.GenerateTestData {
		// Recreate Golden images and exit.
		err = ioutil.WriteFile("testdata/legendAlignment_golden.png", buf.Bytes(), 0o644)
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	err = ioutil.WriteFile("testdata/legendAlignment.png", buf.Bytes(), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("testdata/legendAlignment_golden.png")
	if err != nil {
		t.Fatal(err)
	}

	ok, err := cmpimg.Equal("png", buf.Bytes(), want)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("images differ")
	}

}

func TestIssue514(t *testing.T) {
	for _, ulp := range []int{
		0,
		+1, +2, +3, +4, +5, +6, +7, +8, +9, +10, +11, +12, +13, +14, +15, +16, +17, +18, +19, +20, +21, +22,
		-1, -2, -3, -4, -5, -6, -7, -8, -9, -10, -11, -12, -13, -14, -15, -16, -17, -18, -19, -20, -21, -22,
	} {
		t.Run(fmt.Sprintf("ulps%+02d", ulp), func(t *testing.T) {
			done := make(chan int)
			go func() {
				defer close(done)

				p, err := plot.New()
				if err != nil {
					t.Errorf("could not create plot: %v", err)
					return
				}

				var (
					y1 = 100.0
					y2 = y1
				)

				switch {
				case ulp < 0:
					y2 = math.Float64frombits(math.Float64bits(y1) - uint64(-ulp))
				case ulp > 0:
					y2 = math.Float64frombits(math.Float64bits(y1) + uint64(ulp))
				}

				pts, err := plotter.NewScatter(plotter.XYs{
					{X: 1, Y: y1},
					{X: 1, Y: y2},
				})
				if err != nil {
					t.Errorf("could not create scatter: %v", err)
					return
				}

				p.Add(pts)

				c := draw.NewCanvas(&recorder.Canvas{}, 100, 100)
				p.Draw(c)
			}()

			timeout := time.NewTimer(100 * time.Millisecond)
			defer timeout.Stop()

			select {
			case <-done:
			case <-timeout.C:
				t.Fatalf("could not create plot with small axis range within allotted time")
			}
		})
	}
}
