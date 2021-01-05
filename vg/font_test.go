// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vg

import (
	"errors"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

func TestFontExtends(t *testing.T) {
	for _, tc := range []struct {
		font string
		want map[Length]FontExtents
	}{
		// values obtained when gonum/plot used the package
		// github.com/freetype/truetype for handling fonts.
		{
			font: "Times-Roman",
			want: map[Length]FontExtents{
				10: FontExtents{
					Ascent:  8.9111328125,
					Descent: 2.1630859375,
					Height:  11.4990234375,
				},
				12: FontExtents{
					Ascent:  10.693359375,
					Descent: 2.595703125,
					Height:  13.798828125,
				},
				24: FontExtents{
					Ascent:  21.38671875,
					Descent: 5.19140625,
					Height:  27.59765625,
				},
			},
		},
		{
			font: "Times-Bold",
			want: map[Length]FontExtents{
				10: FontExtents{
					Ascent:  8.9111328125,
					Descent: 2.1630859375,
					Height:  11.4990234375,
				},
				12: FontExtents{
					Ascent:  10.693359375,
					Descent: 2.595703125,
					Height:  13.798828125,
				},
				24: FontExtents{
					Ascent:  21.38671875,
					Descent: 5.19140625,
					Height:  27.59765625,
				},
			},
		},
		{
			font: "Times-Italic",
			want: map[Length]FontExtents{
				10: FontExtents{
					Ascent:  8.9111328125,
					Descent: 2.1630859375,
					Height:  11.4990234375,
				},
				12: FontExtents{
					Ascent:  10.693359375,
					Descent: 2.595703125,
					Height:  13.798828125,
				},
				24: FontExtents{
					Ascent:  21.38671875,
					Descent: 5.19140625,
					Height:  27.59765625,
				},
			},
		},
		{
			font: "Times-BoldItalic",
			want: map[Length]FontExtents{
				10: FontExtents{
					Ascent:  8.9111328125,
					Descent: 2.1630859375,
					Height:  11.4990234375,
				},
				12: FontExtents{
					Ascent:  10.693359375,
					Descent: 2.595703125,
					Height:  13.798828125,
				},
				24: FontExtents{
					Ascent:  21.38671875,
					Descent: 5.19140625,
					Height:  27.59765625,
				},
			},
		},
	} {
		for _, size := range []Length{10, 12, 24} {
			fnt, err := MakeFont(tc.font, size)
			if err != nil {
				t.Fatal(err)
			}

			got := fnt.Extents()
			if got, want := got, tc.want[size]; got != want {
				t.Errorf(
					"invalid font extents for %q, size=%v:\ngot= %#v\nwant=%#v",
					tc.font, size, got, want,
				)
			}
		}
	}
}

func TestFontWidth(t *testing.T) {
	fnt, err := MakeFont("Times-Roman", 12)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		txt  string
		want Length
	}{
		// values obtained when gonum/plot used the package
		// github.com/freetype/truetype for handling fonts.
		{" ", 3},
		{"i", 3.333984375},
		{"q", 6},
		{"A", 8.666015625},
		{"F", 6.673828125},
		{"T", 7.330078125},
		{"V", 8.666015625},
		{"Δ", 7.716796875},
		{"∇", 9.333984375},
		{"  ", 6},
		{"AA", 17.33203125},
		{"Aq", 14.666015625},
		{"Fi", 10.0078125},
		{"Ti", 10.2421875},
		{"AV", 15.78515625},
		{"A∇", 18},
		{"VA", 15.78515625},
		{"VΔ", 16.3828125},
		{"∇Δ", 17.05078125},
		{"   ", 9},
		{"T T", 17.2265625},
	} {
		t.Run(tc.txt, func(t *testing.T) {
			got := fnt.Width(tc.txt)
			if got, want := got, tc.want; got != want {
				t.Fatalf(
					"invalid width: got=%v, want=%v",
					got, want,
				)
			}
		})
	}
}

func TestFontKern(t *testing.T) {
	fnt, err := MakeFont("Times-Roman", 12)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		txt  string
		want fixed.Int26_6
	}{
		// values obtained when gonum/plot used the package
		// github.com/freetype/truetype for handling fonts.
		{"AV", -264},
		{"A∇", 0}, // Liberation has no kerning information for greek symbols
		{"AΔ", 0},
		{"AA", 0},
		{"VA", -264},
		{"∇A", 0}, // Liberation has no kerning information for greek symbols
		{"∇Δ", 0}, // Liberation has no kerning information for greek symbols
	} {
		t.Run(tc.txt, func(t *testing.T) {
			var (
				t0 = rune(tc.txt[0])
				t1 = rune(tc.txt[1])

				buf  sfnt.Buffer
				ppem = fixed.Int26_6(fnt.Font().UnitsPerEm())
			)

			i0, err := fnt.Font().GlyphIndex(&buf, t0)
			if err != nil {
				t.Fatalf("could not find glyph %q: %+v", t0, err)
			}
			i1, err := fnt.Font().GlyphIndex(&buf, t1)
			if err != nil {
				t.Fatalf("could not find glyph %q: %+v", t1, err)
			}
			kern, err := fnt.Font().Kern(&buf, i0, i1, ppem, font.HintingNone)
			switch {
			case err == nil:
				// ok
			case errors.Is(err, sfnt.ErrNotFound):
				kern = 0
			default:
				t.Fatalf("could not find kerning for %q/%q: %+v", t0, t1, err)
			}

			if got, want := kern, tc.want; got != want {
				t.Fatalf("invalid kerning: got=%v, want=%v", got, want)
			}
		})
	}
}
