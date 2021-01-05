// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text_test

import (
	"testing"

	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func TestLatexText(t *testing.T) {
	type box struct{ w, h, d vg.Length }

	tr12, err := vg.MakeFont("Times-Roman", 12)
	if err != nil {
		t.Fatalf("could not make font: %+v", err)
	}

	ti12, err := vg.MakeFont("Times-Italic", 12)
	if err != nil {
		t.Fatalf("could not make font: %+v", err)
	}

	tr42, err := vg.MakeFont("Times-Roman", 42)
	if err != nil {
		t.Fatalf("could not make font: %+v", err)
	}

	for _, tc := range []struct {
		txt string
		fnt vg.Font
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
			box: []box{{23.994140625, 8.328125, 0.125}},
			w:   23.994140625,
			h:   8.453125,
		},
		{
			txt: "hello",
			fnt: ti12,
			box: []box{{23.994140625, 8.328125, 0.125}},
			w:   23.994140625,
			h:   8.453125,
		},
		{
			txt: "hello",
			fnt: tr42,
			box: []box{{83.9794921875, 29.1484375, 0.4375}},
			w:   83.9794921875,
			h:   29.5859375,
		},
		{
			txt: "Agg",
			box: []box{{20.666015625, 7.921875, 2.59375}},
			w:   20.666015625,
			h:   10.515625,
		},
		{
			txt: "Agg",
			fnt: ti12,
			box: []box{{19.330078125, 7.921875, 2.5625}},
			w:   19.330078125,
			h:   10.484375,
		},
		{
			txt: "Agg",
			fnt: tr42,
			box: []box{{72.3310546875, 27.7265625, 9.078125}},
			w:   72.3310546875,
			h:   36.8046875,
		},
		{
			txt: `$Agg$`, // FIXME(sbinet): should be italicized
			box: []box{{20.666015625, 7.921875, 2.59375}},
			w:   20.666015625,
			h:   10.515625,
		},
		{
			txt: `VA`,
			box: []box{{13.20703125, 7.921875, 0.1875}},
			w:   13.20703125,
			h:   8.109375,
		},
		{
			txt: `$V\Delta$`,
			box: []box{{16.3828125, 7.921875, 0.1875}},
			w:   16.3828125,
			h:   8.109375,
		},
	} {
		t.Run(tc.txt, func(t *testing.T) {
			fnt := tc.fnt
			if fnt == (vg.Font{}) {
				fnt = tr12
			}

			sty := draw.TextStyle{
				Font:    fnt,
				Handler: text.Latex{},
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
