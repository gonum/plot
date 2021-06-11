// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgsvg

import (
	"testing"

	stdfnt "golang.org/x/image/font"
	"gonum.org/v1/plot/font"
)

func TestSVGFontDescr(t *testing.T) {
	for i, tc := range []struct {
		fnt  font.Font
		want string
	}{
		// typefaces
		{
			fnt:  font.Font{Typeface: "Liberation"},
			want: "font-family:Liberation;font-variant:normal;font-weight:normal;font-style:normal",
		},
		{
			fnt: font.Font{
				Typeface: "Liberation",
				Variant:  "",
				Style:    stdfnt.StyleNormal,
				Weight:   stdfnt.WeightNormal,
			},
			want: "font-family:Liberation;font-variant:normal;font-weight:normal;font-style:normal",
		},
		{
			fnt: font.Font{
				Typeface: "My-Font-Name",
				Variant:  "",
				Style:    stdfnt.StyleNormal,
				Weight:   stdfnt.WeightNormal,
			},
			want: "font-family:My-Font-Name;font-variant:normal;font-weight:normal;font-style:normal",
		},
		// variants
		{
			fnt: font.Font{
				Typeface: "Liberation",
				Variant:  "Mono",
				Style:    stdfnt.StyleNormal,
				Weight:   stdfnt.WeightNormal,
			},
			want: "font-family:Liberation, monospace;font-variant:monospace;font-weight:normal;font-style:normal",
		},
		{
			fnt: font.Font{
				Typeface: "Liberation",
				Variant:  "Serif",
				Style:    stdfnt.StyleNormal,
				Weight:   stdfnt.WeightNormal,
			},
			want: "font-family:Liberation, serif;font-variant:serif;font-weight:normal;font-style:normal",
		},
		{
			fnt: font.Font{
				Typeface: "Liberation",
				Variant:  "Sans",
				Style:    stdfnt.StyleNormal,
				Weight:   stdfnt.WeightNormal,
			},
			want: "font-family:Liberation, sans-serif;font-variant:sans-serif;font-weight:normal;font-style:normal",
		},
		{
			fnt: font.Font{
				Typeface: "Liberation",
				Variant:  "SansSerif",
				Style:    stdfnt.StyleNormal,
				Weight:   stdfnt.WeightNormal,
			},
			want: "font-family:Liberation, sans-serif;font-variant:sans-serif;font-weight:normal;font-style:normal",
		},
		{
			fnt: font.Font{
				Typeface: "Liberation",
				Variant:  "Sans-Serif",
				Style:    stdfnt.StyleNormal,
				Weight:   stdfnt.WeightNormal,
			},
			want: "font-family:Liberation, sans-serif;font-variant:sans-serif;font-weight:normal;font-style:normal",
		},
		{
			fnt: font.Font{
				Typeface: "Liberation",
				Variant:  "Smallcaps",
				Style:    stdfnt.StyleNormal,
				Weight:   stdfnt.WeightNormal,
			},
			want: "font-family:Liberation, small-caps;font-variant:small-caps;font-weight:normal;font-style:normal",
		},
		// styles
		{
			fnt: font.Font{
				Typeface: "Liberation",
				Variant:  "",
				Style:    stdfnt.StyleItalic,
				Weight:   stdfnt.WeightNormal,
			},
			want: "font-family:Liberation;font-variant:normal;font-weight:normal;font-style:italic",
		},
		{
			fnt: font.Font{
				Typeface: "Liberation",
				Variant:  "",
				Style:    stdfnt.StyleOblique,
				Weight:   stdfnt.WeightNormal,
			},
			want: "font-family:Liberation;font-variant:normal;font-weight:normal;font-style:oblique",
		},
		// weights
		{
			fnt: font.Font{
				Typeface: "Liberation",
				Variant:  "",
				Style:    stdfnt.StyleNormal,
				Weight:   stdfnt.WeightThin,
			},
			want: "font-family:Liberation;font-variant:normal;font-weight:100;font-style:normal",
		},
		{
			fnt: font.Font{
				Typeface: "Liberation",
				Variant:  "",
				Style:    stdfnt.StyleNormal,
				Weight:   stdfnt.WeightBold,
			},
			want: "font-family:Liberation;font-variant:normal;font-weight:bold;font-style:normal",
		},
		{
			fnt: font.Font{
				Typeface: "Liberation",
				Variant:  "",
				Style:    stdfnt.StyleNormal,
			},
			want: "font-family:Liberation;font-variant:normal;font-weight:normal;font-style:normal",
		},
		{
			fnt: font.Font{
				Typeface: "Liberation",
				Variant:  "",
				Style:    stdfnt.StyleNormal,
				Weight:   stdfnt.WeightExtraBold,
			},
			want: "font-family:Liberation;font-variant:normal;font-weight:800;font-style:normal",
		},
		// weights+styles
		{
			fnt: font.Font{
				Typeface: "Times",
				Variant:  "",
				Style:    stdfnt.StyleItalic,
				Weight:   stdfnt.WeightBold,
			},
			want: "font-family:Times;font-variant:normal;font-weight:bold;font-style:italic",
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
