// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gob_test

import (
	"bytes"
	"encoding/gob"
	"image/color"
	"math/rand"
	"os"
	"testing"

	_ "github.com/gonum/plot/gob"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

func init() {
	gob.Register(commaTicks{})
}

func TestPersistency(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	// Get some random points
	n := 15
	scatterData := randomPoints(n, rnd)
	lineData := randomPoints(n, rnd)
	linePointsData := randomPoints(n, rnd)

	p, err := plot.New()
	if err != nil {
		t.Fatalf("error creating plot: %v\n", err)
	}

	p.Title.Text = "Plot Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	// Use a custom tick marker function that computes the default
	// tick marks and re-labels the major ticks with commas.
	p.Y.Tick.Marker = commaTicks{}

	// Draw a grid behind the data
	p.Add(plotter.NewGrid())
	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(scatterData)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}

	// Make a line plotter and set its style.
	l, err := plotter.NewLine(lineData)
	if err != nil {
		panic(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	// Make a line plotter with points and set its style.
	lpLine, lpPoints, err := plotter.NewLinePoints(linePointsData)
	if err != nil {
		panic(err)
	}
	lpLine.Color = color.RGBA{G: 255, A: 255}
	lpPoints.Shape = draw.PyramidGlyph{}
	lpPoints.Color = color.RGBA{R: 255, A: 255}

	// Add the plotters to the plot, with a legend
	// entry for each
	p.Add(s, l, lpLine, lpPoints)
	p.Legend.Add("scatter", s)
	p.Legend.Add("line", l)
	p.Legend.Add("line points", lpLine, lpPoints)

	// Save the plot to a PNG file.
	err = p.Save(4, 4, "test-persistency.png")
	if err != nil {
		t.Fatalf("error saving to PNG: %v\n", err)
	}
	defer os.Remove("test-persistency.png")

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(p)
	if err != nil {
		t.Fatalf("error gob-encoding plot: %v\n", err)
	}

	// TODO(sbinet): impl. BinaryMarshal for plot.Plot and vg.Font
	// {
	// 	dec := gob.NewDecoder(buf)
	// 	var p plot.Plot
	// 	err = dec.Decode(&p)
	// 	if err != nil {
	// 		t.Fatalf("error gob-decoding plot: %v\n", err)
	// 	}
	// 	// Save the plot to a PNG file.
	// 	err = p.Save(4, 4, "test-persistency-readback.png")
	// 	if err != nil {
	// 		t.Fatalf("error saving to PNG: %v\n", err)
	// 	}
	//  defer os.Remove("test-persistency-readback.png")
	// }

}

// randomPoints returns some random x, y points.
func randomPoints(n int, rnd *rand.Rand) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rnd.Float64()
		} else {
			pts[i].X = pts[i-1].X + rnd.Float64()
		}
		pts[i].Y = pts[i].X + rnd.Float64()*1e4
	}
	return pts
}

// CommaTicks computes the default tick marks, but inserts commas
// into the labels for the major tick marks.
type commaTicks struct{}

func (commaTicks) Ticks(min, max float64) []plot.Tick {
	tks := plot.DefaultTicks{}.Ticks(min, max)
	for i, t := range tks {
		if t.Label == "" { // Skip minor ticks, they are fine.
			continue
		}
		tks[i].Label = addCommas(t.Label)
	}
	return tks
}

// AddCommas adds commas after every 3 characters from right to left.
// NOTE: This function is a quick hack, it doesn't work with decimal
// points, and may have a bunch of other problems.
func addCommas(s string) string {
	rev := ""
	n := 0
	for i := len(s) - 1; i >= 0; i-- {
		rev += string(s[i])
		n++
		if n%3 == 0 {
			rev += ","
		}
	}
	s = ""
	for i := len(rev) - 1; i >= 0; i-- {
		s += string(rev[i])
	}
	return s
}
