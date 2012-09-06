// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package vgimg

import (
	"code.google.com/p/plotinum/vg"
	"code.google.com/p/plotinum/vg/vgtest"
	"image/png"
	"testing"
)

func TestFontExtents(t *testing.T) {
	img, err := New(vg.Inches(4), vg.Inches(4))
	if err != nil {
		t.Fatal(err)
	}
	vgtest.DrawFontExtents(t, img)
	err = img.Save("extents.png", png.Encode)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFonts(t *testing.T) {
	img, err := New(vg.Inches(4), vg.Inches(4))
	if err != nil {
		t.Fatal(err)
	}
	vgtest.DrawFonts(t, img)
	err = img.Save("fonts.png", png.Encode)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArcs(t *testing.T) {
	img, err := New(vg.Inches(4), vg.Inches(4))
	if err != nil {
		t.Fatal(err)
	}
	vgtest.DrawArcs(t, img)
	err = img.Save("arcs.png", png.Encode)
	if err != nil {
		t.Fatal(err)
	}
}
