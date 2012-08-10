// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package vecpdf

import (
	"code.google.com/p/plotinum/vg"
	"testing"
)

func TestFontExtents(t *testing.T) {
	pdf := New(vg.Inches(4), vg.Inches(4))
	vg.DrawFontExtents(t, pdf)
	if err := pdf.Save("extents.pdf"); err != nil {
		t.Fatal(err)
	}
}

func TestFonts(t *testing.T) {
	pdf := New(vg.Inches(4), vg.Inches(4))
	vg.DrawFonts(t, pdf)
	if err := pdf.Save("fonts.pdf"); err != nil {
		t.Fatal(err)
	}
}

func TestArcs(t *testing.T) {
	pdf := New(vg.Inches(4), vg.Inches(4))
	vg.DrawArcs(t, pdf)
	if err := pdf.Save("arcs.pdf"); err != nil {
		t.Fatal(err)
	}
}
