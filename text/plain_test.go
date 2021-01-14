// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text_test

import (
	"testing"

	stdfnt "golang.org/x/image/font"

	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/font/liberation"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
)

func TestPlainText(t *testing.T) {
	type box struct{ w, h, d vg.Length }

	fonts := font.NewCache(liberation.Collection())

	tr12 := font.Font{Variant: "Serif", Size: 12}
	ti12 := font.Font{Variant: "Serif", Size: 12, Style: stdfnt.StyleItalic}
	tr42 := font.Font{Variant: "Serif", Size: 42}

	for _, tc := range []struct {
		txt string
		fnt font.Font
		box []box
		w   vg.Length
		h   vg.Length
	}{
		{
			txt: "",
			box: []box{
				{0, 10.693359375, 2.595703125},
			},
			w: 0,
			h: 13.2890625,
		},
		{
			txt: " ",
			box: []box{{3, 10.693359375, 2.595703125}},
			w:   3,
			h:   13.2890625,
		},
		{
			txt: "hello",
			box: []box{{23.994140625, 10.693359375, 2.595703125}},
			w:   23.994140625,
			h:   13.2890625,
		},
		{
			txt: "hello",
			fnt: ti12,
			box: []box{{23.994140625, 10.693359375, 2.595703125}},
			w:   23.994140625,
			h:   13.2890625,
		},
		{
			txt: "hello",
			fnt: tr42,
			box: []box{{83.9794921875, 37.4267578125, 9.0849609375}},
			w:   83.9794921875,
			h:   46.51171875,
		},
		{
			txt: "hello\n",
			box: []box{{23.994140625, 10.693359375, 2.595703125}},
			w:   23.994140625,
			h:   13.2890625,
		},
		{
			txt: "Agg",
			box: []box{{20.666015625, 10.693359375, 2.595703125}},
			w:   20.666015625,
			h:   13.2890625,
		},
		{
			txt: "Agg",
			fnt: ti12,
			box: []box{{19.330078125, 10.693359375, 2.595703125}},
			w:   19.330078125,
			h:   13.2890625,
		},
		{
			txt: "Agg",
			fnt: tr42,
			box: []box{{72.3310546875, 37.4267578125, 9.0849609375}},
			w:   72.3310546875,
			h:   46.51171875,
		},
		{
			txt: "\n",
			box: []box{
				{0, 10.693359375, 2.595703125},
			},
			w: 0,
			h: 13.2890625,
		},
		{
			txt: "\n ",
			box: []box{
				{0, 10.693359375, 2.595703125},
				{3, 10.693359375, 2.595703125},
			},
			w: 3,
			h: 27.087890625,
		},
		{
			txt: " \n ",
			box: []box{
				{3, 10.693359375, 2.595703125},
				{3, 10.693359375, 2.595703125},
			},
			w: 3,
			h: 27.087890625,
		},
		{
			txt: "hello\nworld",
			box: []box{
				{23.994140625, 10.693359375, 2.595703125},
				{27.996093750, 10.693359375, 2.595703125},
			},
			w: 27.996093750,
			h: 27.087890625,
		},
		{
			txt: "Agg\nBpp",
			box: []box{
				{20.666015625, 10.693359375, 2.595703125},
				{20.003906250, 10.693359375, 2.595703125},
			},
			w: 20.666015625,
			h: 27.087890625,
		},
	} {
		t.Run(tc.txt, func(t *testing.T) {
			fnt := tc.fnt
			if fnt == (font.Font{}) {
				fnt = tr12
			}

			sty := text.Style{
				Font:    fnt,
				Handler: &text.Plain{Fonts: fonts},
			}

			lines := sty.Handler.Lines(tc.txt)
			if got, want := len(lines), len(tc.box); got != want {
				t.Errorf("invalid number of lines: got=%d, want=%d", got, want)
			}

			for i, line := range lines {
				var b box
				b.w, b.h, b.d = sty.Handler.Box(line, sty.Font)

				if got, want := b, tc.box[i]; got != want {
					t.Errorf("invalid box[%d]: got=%v, want=%v", i, got, want)
				}
			}

			w := sty.Width(tc.txt)
			if got, want := w, tc.w; got != want {
				t.Errorf("invalid width: got=%v, want=%v", got, want)
			}

			h := sty.Height(tc.txt)
			if got, want := h, tc.h; got != want {
				t.Errorf("invalid height: got=%v, want=%v", got, want)
			}
		})
	}
}
