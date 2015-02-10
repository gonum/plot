// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package recorder

import (
	"image/color"
	"strings"
	"testing"

	"github.com/gonum/plot/vg"
)

func TestRecorder(t *testing.T) {
	rec := New(72)
	rec.Actions = append(rec.Actions, &FillString{Font: "Times-Roman", Size: 12, X: 0, Y: 10, String: "Text"})
	rec.Comment("End of preamble")
	rec.Scale(1, 2)
	rec.Rotate(0.72)
	rec.KeepCaller = true
	rec.Stroke(vg.Path{{Type: vg.MoveComp, X: 3, Y: 4}})
	rec.Push()
	rec.Pop()
	rec.Translate(3, 4)
	rec.KeepCaller = false
	rec.DPI()
	rec.SetLineWidth(100)
	rec.SetLineDash([]vg.Length{2, 5}, 6)
	rec.SetColor(color.RGBA{R: 0x65, G: 0x23, B: 0xf2})
	rec.Fill(vg.Path{{Type: vg.MoveComp, X: 3, Y: 4}, {Type: vg.LineComp, X: 2, Y: 3}, {Type: vg.CloseComp}})
	if len(rec.Actions) != len(want) {
		t.Fatalf("unexpected number of actions recorded: got:%d want:%d", len(rec.Actions), len(want))
	}
	for i, a := range rec.Actions {
		if got := a.Call(); !strings.HasSuffix(got, want[i]) {
			t.Errorf("unexpected action:\n\tgot: %#v\n\twant: %#v", got, want[i])
		}
	}

	var replay Canvas
	err := rec.ReplayOn(&replay)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	for i, a := range rec.Actions {
		got := replay.Actions[i].Call()
		want := a.Call()
		if !strings.HasSuffix(want, got) {
			t.Errorf("unexpected action:\n\tgot: %#v\n\twant: %#v", got, want)
		}
	}

	replay.Reset()
	rec.Actions = append(rec.Actions, &FillString{Font: "Foo", Size: 12, X: 0, Y: 10, String: "Bar"})
	err = rec.ReplayOn(&replay)
	if !strings.HasPrefix(err.Error(), "Unknown font: Foo.") {
		t.Errorf("unexpected error: %v", err)
	}
}

var want = []string{
	`FillString("Times-Roman", 12, 0, 10, "Text")`,
	`Comment("End of preamble")`,
	`Scale(1, 2)`,
	`Rotate(0.72)`,
	`github.com/gonum/plot/vg/recorder/recorder_test.go:22 Stroke(vg.Path{vg.PathComp{Type:0, X:3, Y:4, Radius:0, Start:0, Angle:0}})`,
	`github.com/gonum/plot/vg/recorder/recorder_test.go:23 Push()`,
	`github.com/gonum/plot/vg/recorder/recorder_test.go:24 Pop()`,
	`github.com/gonum/plot/vg/recorder/recorder_test.go:25 Translate(3, 4)`,
	`DPI()`,
	`SetLineWidth(100)`,
	`SetLineDash([]vg.Length{2, 5}, 6)`,
	`SetColor(color.RGBA{R:0x65, G:0x23, B:0xf2, A:0x0})`,
	`Fill(vg.Path{vg.PathComp{Type:0, X:3, Y:4, Radius:0, Start:0, Angle:0}, vg.PathComp{Type:1, X:2, Y:3, Radius:0, Start:0, Angle:0}, vg.PathComp{Type:3, X:0, Y:0, Radius:0, Start:0, Angle:0}})`,
}
