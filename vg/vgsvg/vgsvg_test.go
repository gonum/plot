// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package vgsvg

import (
	"code.google.com/p/plotinum/vg"
	"code.google.com/p/plotinum/vg/vgtest"
	"testing"
)

func TestFontExtents(t *testing.T) {
	img := New(vg.Inches(4), vg.Inches(4))
	vgtest.DrawFontExtents(t, img)
	err := img.Save("extents.svg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFonts(t *testing.T) {
	img := New(vg.Inches(4), vg.Inches(4))
	vgtest.DrawFonts(t, img)
	err := img.Save("fonts.svg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestArcss(t *testing.T) {
	img := New(vg.Inches(4), vg.Inches(4))
	vgtest.DrawArcs(t, img)
	err := img.Save("arcs.svg")
	if err != nil {
		t.Fatal(err)
	}
}
