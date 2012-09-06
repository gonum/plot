// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package vgeps

import (
	"code.google.com/p/plotinum/vg"
	"code.google.com/p/plotinum/vg/vgtest"
	"testing"
)

func TestFontExtents(t *testing.T) {
	eps := New(vg.Inches(4), vg.Inches(4), "extents")
	vgtest.DrawFontExtents(t, eps)
	if err := eps.Save("extents.eps"); err != nil {
		t.Fatal(err)
	}
}

func TestFonts(t *testing.T) {
	eps := New(vg.Inches(4), vg.Inches(4), "fonts")
	vgtest.DrawFonts(t, eps)
	if err := eps.Save("fonts.eps"); err != nil {
		t.Fatal(err)
	}
}

func TestArcs(t *testing.T) {
	eps := New(vg.Inches(4), vg.Inches(4), "arcs")
	vgtest.DrawArcs(t, eps)
	if err := eps.Save("arcs.eps"); err != nil {
		t.Fatal(err)
	}
}
