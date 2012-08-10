// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package vecimg

import (
	"code.google.com/p/plotinum/vg"
	"testing"
)

func TestFontExtents(t *testing.T) {
	img, err := New(vg.Inches(4), vg.Inches(4))
	if err != nil {
		t.Fatal(err)
	}
	vg.DrawFontExtents(t, img)
	err = img.SavePNG("extents.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFonts(t *testing.T) {
	img, err := New(vg.Inches(4), vg.Inches(4))
	if err != nil {
		t.Fatal(err)
	}
	vg.DrawFonts(t, img)
	err = img.SaveJPEG("fonts.jpg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestArcs(t *testing.T) {
	img, err := New(vg.Inches(4), vg.Inches(4))
	if err != nil {
		t.Fatal(err)
	}
	vg.DrawArcs(t, img)
	err = img.SavePNG("arcs.png")
	if err != nil {
		t.Fatal(err)
	}
}
