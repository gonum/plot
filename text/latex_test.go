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

func TestLatexText(t *testing.T) {
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
			box: []box{{0, 0, 0}},
		},
		{
			txt: " ",
			box: []box{{3, 0, 0}},
			w:   3,
		},
		{
			txt: "hello",
			box: []box{{23.994140625, 8.328125, 0.634765625}},
			w:   23.994140625,
			h:   8.962890625,
		},
		{
			txt: "hello",
			fnt: ti12,
			box: []box{{23.994140625, 8.328125, 0.634765625}},
			w:   23.994140625,
			h:   8.962890625,
		},
		{
			txt: `$hello$`,
			box: []box{{23.994140625, 8.328125, 0.634765625}},
			w:   23.994140625,
			h:   8.962890625,
		},
		{
			txt: "hello",
			fnt: tr42,
			box: []box{{83.9794921875, 29.1484375, 2.2216796875}},
			w:   83.9794921875,
			h:   31.3701171875,
		},
		{
			txt: "Agg",
			box: []box{{20.666015625, 7.921875, 3.103515625}},
			w:   20.666015625,
			h:   11.025390625,
		},
		{
			txt: "Agg",
			fnt: ti12, // italics correctly ignored.
			box: []box{{20.666015625, 7.921875, 3.103515625}},
			w:   20.666015625,
			h:   11.025390625,
		},
		{
			txt: `$Agg$`,
			box: []box{{19.330078125, 7.921875, 3.072265625}},
			w:   19.330078125,
			h:   10.994140625,
		},
		{
			txt: "Agg",
			fnt: tr42,
			box: []box{{72.3310546875, 27.7265625, 10.8623046875}},
			w:   72.3310546875,
			h:   38.5888671875,
		},
		{
			txt: `VA`,
			box: []box{{13.20703125, 7.921875, 0.697265625}},
			w:   13.20703125,
			h:   8.619140625,
		},
		{
			txt: `$V\Delta$`,
			box: []box{{14.373046875, 7.921875, 0.697265625}},
			w:   14.373046875,
			h:   8.619140625,
		},
	} {
		t.Run(tc.txt, func(t *testing.T) {
			fnt := tc.fnt
			if fnt == (font.Font{}) {
				fnt = tr12
			}

			sty := text.Style{
				Font:    fnt,
				Handler: &text.Latex{Fonts: fonts},
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
