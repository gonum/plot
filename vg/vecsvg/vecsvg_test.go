// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package vecsvg

import (
	"code.google.com/p/plotinum/vg"
	"testing"
)

func TestFontExtents(t *testing.T) {
	img := New(vg.Inches(4), vg.Inches(4))
	vg.DrawFontExtents(t, img)
	err := img.Save("extents.svg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFonts(t *testing.T) {
	img := New(vg.Inches(4), vg.Inches(4))
	vg.DrawFonts(t, img)
	err := img.Save("fonts.svg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestArcss(t *testing.T) {
	img := New(vg.Inches(4), vg.Inches(4))
	vg.DrawArcs(t, img)
	err := img.Save("arcs.svg")
	if err != nil {
		t.Fatal(err)
	}
}
