// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package vgeps

import (
	"code.google.com/p/plotinum/vg"
	"testing"
)

func TestFontExtents(t *testing.T) {
	eps := New(vg.Inches(4), vg.Inches(4), "extents")
	vg.DrawFontExtents(t, eps)
	if err := eps.Save("extents.eps"); err != nil {
		t.Fatal(err)
	}
}

func TestFonts(t *testing.T) {
	eps := New(vg.Inches(4), vg.Inches(4), "fonts")
	vg.DrawFonts(t, eps)
	if err := eps.Save("fonts.eps"); err != nil {
		t.Fatal(err)
	}
}

func TestArcs(t *testing.T) {
	eps := New(vg.Inches(4), vg.Inches(4), "arcs")
	vg.DrawArcs(t, eps)
	if err := eps.Save("arcs.eps"); err != nil {
		t.Fatal(err)
	}
}
