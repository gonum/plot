// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/recorder"
)

func TestSankey_simple(t *testing.T) {
	cmpimg.CheckPlot(ExampleSankey_simple, t, "sankeySimple.png")
}

func TestSankey_grouped(t *testing.T) {
	cmpimg.CheckPlot(ExampleSankey_grouped, t, "sankeyGrouped.png")
}

// This test checks whether the Sankey plotter makes any changes to
// the input Flows.
func TestSankey_idempotent(t *testing.T) {
	flows := []plotter.Flow{
		{
			SourceCategory:   0,
			SourceLabel:      "Large",
			ReceptorCategory: 1,
			ReceptorLabel:    "Mohamed",
			Value:            5,
		},
		{
			SourceCategory:   0,
			SourceLabel:      "Small",
			ReceptorCategory: 1,
			ReceptorLabel:    "Sofia",
			Value:            5,
		},
	}
	s, err := plotter.NewSankey(flows...)
	if err != nil {
		t.Fatal(err)
	}
	p, err := plot.New()
	if err != nil {
		t.Fatal(err)
	}
	p.Add(s)
	p.HideAxes()

	// Draw the plot once.
	c1 := new(recorder.Canvas)
	dc1 := draw.NewCanvas(c1, vg.Centimeter, vg.Centimeter)
	p.Draw(dc1)

	// Draw the plot a second time.
	c2 := new(recorder.Canvas)
	dc2 := draw.NewCanvas(c2, vg.Centimeter, vg.Centimeter)
	p.Draw(dc2)

	if len(c1.Actions) != len(c2.Actions) {
		t.Errorf("inconsistent number of actions: %d != %d", len(c2.Actions), len(c1.Actions))
	}

	for i, a1 := range c1.Actions {
		if a1.Call() != c2.Actions[i].Call() {
			t.Errorf("action %d: %s\n\t!= %s", i, c2.Actions[i].Call(), a1.Call())
		}
	}
}
