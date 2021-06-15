// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgsvg

import (
	"testing"

	"github.com/go-fonts/latin-modern/lmroman10regular"
	"github.com/go-fonts/liberation/liberationmonoregular"
	"github.com/go-fonts/liberation/liberationsansregular"
	"github.com/go-fonts/liberation/liberationserifbold"
	"github.com/go-fonts/liberation/liberationserifbolditalic"
	"github.com/go-fonts/liberation/liberationserifitalic"
	"github.com/go-fonts/liberation/liberationserifregular"
	xfnt "golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"

	"gonum.org/v1/plot/font"
)

func TestSVGFontDescr(t *testing.T) {
	newFace := func(fnt font.Font, raw []byte) font.Face {
		ttf, err := sfnt.Parse(raw)
		if err != nil {
			t.Fatalf("could not parse %q: %+v", fnt.Typeface, err)
		}
		return font.Face{Font: fnt, Face: ttf}
	}

	for i, tc := range []struct {
		fnt  font.Face
		want string
	}{
		// typefaces
		{
			fnt: newFace(
				font.Font{Typeface: "Liberation"},
				liberationserifregular.TTF,
			),
			want: "font-family:Liberation Serif;font-variant:none;font-weight:normal;font-style:normal",
		},
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "",
					Style:    xfnt.StyleNormal,
					Weight:   xfnt.WeightNormal,
				},
				liberationserifregular.TTF,
			),
			want: "font-family:Liberation Serif;font-variant:none;font-weight:normal;font-style:normal",
		},
		{
			fnt: newFace(
				font.Font{
					Typeface: "Latin Modern",
					Variant:  "",
					Style:    xfnt.StyleNormal,
					Weight:   xfnt.WeightNormal,
				},
				lmroman10regular.TTF,
			),
			want: "font-family:Latin Modern Roman;font-variant:none;font-weight:normal;font-style:normal",
		},
		// variants
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "Mono",
					Style:    xfnt.StyleNormal,
					Weight:   xfnt.WeightNormal,
				},
				liberationmonoregular.TTF,
			),
			want: "font-family:Liberation Mono;font-variant:normal;font-weight:normal;font-style:normal",
		},
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "Serif",
					Style:    xfnt.StyleNormal,
					Weight:   xfnt.WeightNormal,
				},
				liberationserifregular.TTF,
			),
			want: "font-family:Liberation Serif;font-variant:normal;font-weight:normal;font-style:normal",
		},
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "Sans",
					Style:    xfnt.StyleNormal,
					Weight:   xfnt.WeightNormal,
				},
				liberationsansregular.TTF,
			),
			want: "font-family:Liberation Sans;font-variant:normal;font-weight:normal;font-style:normal",
		},
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "SansSerif",
					Style:    xfnt.StyleNormal,
					Weight:   xfnt.WeightNormal,
				},
				liberationsansregular.TTF,
			),
			want: "font-family:Liberation Sans;font-variant:normal;font-weight:normal;font-style:normal",
		},
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "Sans-Serif",
					Style:    xfnt.StyleNormal,
					Weight:   xfnt.WeightNormal,
				},
				liberationsansregular.TTF,
			),
			want: "font-family:Liberation Sans;font-variant:normal;font-weight:normal;font-style:normal",
		},
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "Smallcaps",
					Style:    xfnt.StyleNormal,
					Weight:   xfnt.WeightNormal,
				},
				liberationserifregular.TTF,
			),
			want: "font-family:Liberation Serif;font-variant:small-caps;font-weight:normal;font-style:normal",
		},
		// styles
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "",
					Style:    xfnt.StyleItalic,
					Weight:   xfnt.WeightNormal,
				},
				liberationserifitalic.TTF,
			),
			want: "font-family:Liberation Serif;font-variant:none;font-weight:normal;font-style:italic",
		},
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "",
					Style:    xfnt.StyleOblique,
					Weight:   xfnt.WeightNormal,
				},
				liberationserifitalic.TTF,
			),
			want: "font-family:Liberation Serif;font-variant:none;font-weight:normal;font-style:oblique",
		},
		// weights
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "",
					Style:    xfnt.StyleNormal,
					Weight:   xfnt.WeightThin,
				},
				liberationserifregular.TTF,
			),
			want: "font-family:Liberation Serif;font-variant:none;font-weight:100;font-style:normal",
		},
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "",
					Style:    xfnt.StyleNormal,
					Weight:   xfnt.WeightBold,
				},
				liberationserifbold.TTF,
			),
			want: "font-family:Liberation Serif;font-variant:none;font-weight:bold;font-style:normal",
		},
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "",
					Style:    xfnt.StyleNormal,
				},
				liberationserifregular.TTF,
			),
			want: "font-family:Liberation Serif;font-variant:none;font-weight:normal;font-style:normal",
		},
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "",
					Style:    xfnt.StyleNormal,
					Weight:   xfnt.WeightExtraBold,
				},
				liberationserifbold.TTF,
			),
			want: "font-family:Liberation Serif;font-variant:none;font-weight:800;font-style:normal",
		},
		// weights+styles
		{
			fnt: newFace(
				font.Font{
					Typeface: "Liberation",
					Variant:  "",
					Style:    xfnt.StyleItalic,
					Weight:   xfnt.WeightBold,
				},
				liberationserifbolditalic.TTF,
			),
			want: "font-family:Liberation Serif;font-variant:none;font-weight:bold;font-style:italic",
		},
	} {
		got := svgFontDescr(tc.fnt)
		if got != tc.want {
			t.Errorf(
				"invalid SVG font[%d] description:\ngot= %s\nwant=%s",
				i, got, tc.want,
			)
		}
	}
}
