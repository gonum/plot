// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot_test

import (
	"bytes"
	"fmt"
	"image/color"
	"reflect"
	"testing"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/recorder"
)

func TestLegendAlignment(t *testing.T) {
	font, err := vg.MakeFont(plot.DefaultFont, 10.822510822510822) // This font size gives an entry height of 10.
	if err != nil {
		t.Fatalf("failed to create font: %v", err)
	}
	l := plot.Legend{
		ThumbnailWidth: vg.Points(20),
		TextStyle:      draw.TextStyle{Font: font},
	}
	for _, n := range []string{"A", "B", "C", "D"} {
		b, err := plotter.NewBarChart(plotter.Values{0}, 1)
		if err != nil {
			t.Fatalf("failed to create bar chart %q: %v", n, err)
		}
		l.Add(n, b)
	}

	var r recorder.Canvas
	c := draw.NewCanvas(&r, 100, 100)
	l.Draw(draw.Canvas{
		Canvas: c.Canvas,
		Rectangle: draw.Rectangle{
			Min: draw.Point{0, 0},
			Max: draw.Point{100, 100},
		},
	})

	got := r.Actions

	// want is a snapshot of the actions for the code above when the
	// graphical output has been visually confirmed to be correct for
	// the bar charts example show in gonum/plot#25.
	want := []recorder.Action{
		&recorder.SetColor{
			Color: color.Gray16{},
		},
		&recorder.Fill{
			Path: vg.Path{
				{Type: vg.MoveComp, X: 80, Y: 30},
				{Type: vg.LineComp, X: 80, Y: 40},
				{Type: vg.LineComp, X: 100, Y: 40},
				{Type: vg.LineComp, X: 100, Y: 30},
				{Type: vg.CloseComp},
			},
		},
		&recorder.SetColor{
			Color: color.Gray16{},
		},
		&recorder.SetLineWidth{
			Width: 1,
		},
		&recorder.SetLineDash{},
		&recorder.Stroke{
			Path: vg.Path{
				{Type: vg.MoveComp, X: 80, Y: 30},
				{Type: vg.LineComp, X: 80, Y: 40},
				{Type: vg.LineComp, X: 100, Y: 40},
				{Type: vg.LineComp, X: 100, Y: 30},
				{Type: vg.LineComp, X: 80, Y: 30},
			},
		},
		&recorder.SetColor{},
		&recorder.FillString{
			Font:   string("Times-Roman"),
			Size:   10.822510822510822,
			X:      69.48051948051948,
			Y:      30.82251082251082,
			String: "A",
		},
		&recorder.SetColor{
			Color: color.Gray16{},
		},
		&recorder.Fill{
			Path: vg.Path{
				{Type: vg.MoveComp, X: 80, Y: 20},
				{Type: vg.LineComp, X: 80, Y: 30},
				{Type: vg.LineComp, X: 100, Y: 30},
				{Type: vg.LineComp, X: 100, Y: 20},
				{Type: vg.CloseComp},
			},
		},
		&recorder.SetColor{
			Color: color.Gray16{},
		},
		&recorder.SetLineWidth{
			Width: 1,
		},
		&recorder.SetLineDash{},
		&recorder.Stroke{
			Path: vg.Path{
				{Type: vg.MoveComp, X: 80, Y: 20},
				{Type: vg.LineComp, X: 80, Y: 30},
				{Type: vg.LineComp, X: 100, Y: 30},
				{Type: vg.LineComp, X: 100, Y: 20},
				{Type: vg.LineComp, X: 80, Y: 20},
			},
		},
		&recorder.SetColor{},
		&recorder.FillString{
			Font:   string("Times-Roman"),
			Size:   10.822510822510822,
			X:      70.07575757575758,
			Y:      20.82251082251082,
			String: "B",
		},
		&recorder.SetColor{
			Color: color.Gray16{
				Y: uint16(0),
			},
		},
		&recorder.Fill{
			Path: vg.Path{
				{Type: vg.MoveComp, X: 80, Y: 10},
				{Type: vg.LineComp, X: 80, Y: 20},
				{Type: vg.LineComp, X: 100, Y: 20},
				{Type: vg.LineComp, X: 100, Y: 10},
				{Type: vg.CloseComp},
			},
		},
		&recorder.SetColor{
			Color: color.Gray16{},
		},
		&recorder.SetLineWidth{
			Width: 1,
		},
		&recorder.SetLineDash{},
		&recorder.Stroke{
			Path: vg.Path{
				{Type: vg.MoveComp, X: 80, Y: 10},
				{Type: vg.LineComp, X: 80, Y: 20},
				{Type: vg.LineComp, X: 100, Y: 20},
				{Type: vg.LineComp, X: 100, Y: 10},
				{Type: vg.LineComp, X: 80, Y: 10},
			},
		},
		&recorder.SetColor{},
		&recorder.FillString{
			Font:   string("Times-Roman"),
			Size:   10.822510822510822,
			X:      70.07575757575758,
			Y:      10.822510822510822,
			String: "C",
		},
		&recorder.SetColor{
			Color: color.Gray16{},
		},
		&recorder.Fill{
			Path: vg.Path{
				{Type: vg.MoveComp, X: 80, Y: 0},
				{Type: vg.LineComp, X: 80, Y: 10},
				{Type: vg.LineComp, X: 100, Y: 10},
				{Type: vg.LineComp, X: 100, Y: 0},
				{Type: vg.CloseComp},
			},
		},
		&recorder.SetColor{
			Color: color.Gray16{},
		},
		&recorder.SetLineWidth{
			Width: 1,
		},
		&recorder.SetLineDash{},
		&recorder.Stroke{
			Path: vg.Path{
				{Type: vg.MoveComp, X: 80, Y: 0},
				{Type: vg.LineComp, X: 80, Y: 10},
				{Type: vg.LineComp, X: 100, Y: 10},
				{Type: vg.LineComp, X: 100, Y: 0},
				{Type: vg.LineComp, X: 80, Y: 0},
			},
		},
		&recorder.SetColor{},
		&recorder.FillString{
			Font:   string("Times-Roman"),
			Size:   10.822510822510822,
			X:      69.48051948051948,
			Y:      0.8225108225108215,
			String: "D",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected legend actions:\ngot:\n%s\nwant:\n%s", formatActions(got), formatActions(want))
		t.Errorf("First diff:\n%s", printActionDiff(got, want))
	}
}

func formatActions(actions []recorder.Action) string {
	var buf bytes.Buffer
	for _, a := range actions {
		fmt.Fprintf(&buf, "\t%s\n", a.Call())
	}
	return buf.String()
}

// printActionDiff prints the first line that is different between two actions.
func printActionDiff(got, want []recorder.Action) string {
	var buf bytes.Buffer
	for i, g := range got {
		if i >= len(want) {
			fmt.Fprintf(&buf, "line %d:\n\tgot: %s\n\twant is empty", i, g.Call())
			break
		}
		w := want[i]
		if w.Call() != g.Call() {
			fmt.Fprintf(&buf, "line %d:\n\tgot: %s\n\twant: %s", i, g.Call(), w.Call())
			break
		}
	}
	if len(want) > len(got) {
		fmt.Fprintf(&buf, "line %d:\n\twant: %s\n\tgot is empty", len(got), want[len(got)].Call())
	}
	return buf.String()
}
