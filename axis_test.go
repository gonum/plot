// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/recorder"
)

var axisSmallTickTests = []struct {
	min, max   float64
	wantValues []float64
	wantLabels []string
}{
	{
		min:        -1.9846500878911073,
		max:        0.4370974820125605,
		wantValues: []float64{-1.75, -0.75, 0.25},
		wantLabels: []string{"-1.75", "-0.75", "0.25"},
	},
	{
		min:        -1.985e15,
		max:        0.4371e15,
		wantValues: []float64{-1.75e15, -7.5e14, 2.5e14},
		wantLabels: []string{"-1.75e+15", "-7.5e+14", "2.5e+14"},
	},
	{
		min:        -1.985e-15,
		max:        0.4371e-15,
		wantValues: []float64{-1.985e-15, -7.739500000000001e-16, 4.3709999999999994e-16},
		wantLabels: []string{"-1.985e-15", "-7.7395e-16", "4.371e-16"},
	},
	{
		min:        math.MaxFloat64 / 4,
		max:        math.MaxFloat64 / 3,
		wantValues: []float64{4.4942328371557893e+307, 5.243271643348421e+307, 5.992310449541053e+307},
		wantLabels: []string{"4e+307", "5e+307", "6e+307"},
	},
	{
		min:        0.00010,
		max:        0.00015,
		wantValues: []float64{0.0001, 0.00012, 0.00014000000000000001},
		wantLabels: []string{"0.0001", "0.00012", "0.00014"},
	},
	{
		min:        555.6545,
		max:        21800.9875,
		wantValues: []float64{4000, 12000, 20000},
		wantLabels: []string{"4000", "12000", "20000"},
	},
	{
		min:        555.6545,
		max:        27800.9875,
		wantValues: []float64{5000, 15000, 25000},
		wantLabels: []string{"5000", "15000", "25000"},
	},
	{
		min:        55.6545,
		max:        1555.9875,
		wantValues: []float64{300, 900, 1500},
		wantLabels: []string{"300", "900", "1500"},
	},
	{
		min:        3.096916 - 0.125,
		max:        3.096916 + 0.125,
		wantValues: []float64{3, 3.1, 3.2},
		wantLabels: []string{"3.0", "3.1", "3.2"},
	},
}

func TestAxisSmallTick(t *testing.T) {
	d := DefaultTicks{}
	for i, test := range axisSmallTickTests {
		ticks := d.Ticks(test.min, test.max)
		gotLabels := labelsOf(ticks)
		gotValues := valuesOf(ticks)
		if !reflect.DeepEqual(gotValues, test.wantValues) {
			t.Errorf("tick values mismatch %d:\ngot: %v\nwant:%v", i, gotValues, test.wantValues)
		}
		if !reflect.DeepEqual(gotLabels, test.wantLabels) {
			t.Errorf("tick labels mismatch %d:\ngot: %q\nwant:%q", i, gotLabels, test.wantLabels)
		}
	}
}

func valuesOf(ticks []Tick) []float64 {
	var values []float64
	for _, t := range ticks {
		if t.Label != "" {
			values = append(values, t.Value)
		}
	}
	return values
}

func labelsOf(ticks []Tick) []string {
	var labels []string
	for _, t := range ticks {
		if t.Label != "" {
			labels = append(labels, t.Label)
		}
	}
	return labels
}

func TestTickerFunc_Ticks(t *testing.T) {
	type args struct {
		min float64
		max float64
	}
	tests := []struct {
		name string
		args args
		want []Tick
		f    TickerFunc
	}{
		{
			name: "return exactly the same ticks as the function passed to TickerFunc",
			args: args{0, 3},
			want: []Tick{{1, "a"}, {2, "b"}},
			f: func(min, max float64) []Tick {
				return []Tick{{1, "a"}, {2, "b"}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Ticks(tt.args.min, tt.args.max); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TickerFunc.Ticks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvertedScale_Normalize(t *testing.T) {
	inverter := InvertedScale{Normalizer: LinearScale{}}
	if got := inverter.Normalize(0, 1, 1); got != 0.0 {
		t.Errorf("Expected a normalization inversion %f->%f not %f", 1.0, 0.0, got)
	}
	if got := inverter.Normalize(0, 1, .5); got != 0.5 {
		t.Errorf("Expected a normalization inversion %f->%f not %f", 0.5, 0.5, got)
	}
	if got := inverter.Normalize(0, 1, 0); got != 1.0 {
		t.Errorf("Expected a normalization inversion %f->%f not %f", 0.0, 1.0, got)
	}
}

func TestAxisPadding(t *testing.T) {
	for _, padding := range []int{0, 5, 10} {
		t.Run(fmt.Sprintf("padding-%d", padding), func(t *testing.T) {
			cmpimg.CheckPlot(func() {
				p := New()
				p.Title.Text = fmt.Sprintf("padding=%d", padding)
				p.X.Label.Text = "X-Axis"
				p.X.Label.Padding = vg.Points(float64(padding))
				p.Y.Label.Text = "Y-Axis"
				p.Y.Label.Padding = vg.Points(float64(padding))

				const size = 5 * vg.Centimeter
				err := p.Save(size, size, fmt.Sprintf("testdata/axis_padding_%02d.png", padding))
				if err != nil {
					t.Fatalf("error: %+v", err)
				}
			}, t, fmt.Sprintf("axis_padding_%02d.png", padding))
		})
	}
}

// labelBoundsHandler records the bounds requested by axis drawing.
type labelBoundsHandler struct {
	text.Handler
	bounds vg.Rectangle
}

func (h *labelBoundsHandler) Draw(_ vg.Canvas, label string, style text.Style, point vg.Point) {
	h.bounds = style.Rectangle(label).Add(point)
}

func TestCenteredAxisLabelGlyphBoxes(t *testing.T) {
	tests := []struct {
		name     string
		scale    Normalizer
		min, max float64
	}{
		{name: "linear", scale: LinearScale{}, min: -10, max: 20},
		{name: "log", scale: LogScale{}, min: 1, max: 100},
		{name: "inverted-log", scale: InvertedScale{LogScale{}}, min: 1, max: 100},
		{name: "large-linear", scale: LinearScale{}, min: math.MaxFloat64 / 2, max: math.MaxFloat64},
	}
	for _, test := range tests {
		for _, horizontal := range []bool{true, false} {
			t.Run(fmt.Sprintf("%s/horizontal=%t", test.name, horizontal), func(t *testing.T) {
				p := New()
				axis := &p.Y
				if horizontal {
					axis = &p.X
				}
				axis.Min, axis.Max, axis.Scale = test.min, test.max, test.scale
				axis.Label.Text = "Axis label"
				axis.Tick.Marker = ConstantTicks(nil)
				handler := &labelBoundsHandler{Handler: axis.Label.TextStyle.Handler}
				axis.Label.TextStyle.Handler = handler
				canvas := draw.Canvas{
					Canvas:    new(recorder.Canvas),
					Rectangle: vg.Rectangle{Min: vg.Point{X: 100, Y: 150}, Max: vg.Point{X: 500, Y: 450}},
				}
				var boxes []GlyphBox
				if horizontal {
					a := horizontalAxis{*axis}
					a.draw(canvas)
					boxes = a.GlyphBoxes(p)
				} else {
					a := verticalAxis{*axis}
					a.draw(canvas)
					boxes = a.GlyphBoxes(p)
				}
				if len(boxes) != 1 {
					t.Fatalf("got %d glyph boxes, want 1", len(boxes))
				}
				box := boxes[0]
				got := box.Rectangle.Add(vg.Point{X: canvas.X(box.X), Y: canvas.Y(box.Y)})
				want := handler.bounds
				if math.IsNaN(float64(got.Min.X)) || math.IsNaN(float64(got.Min.Y)) ||
					math.IsNaN(float64(got.Max.X)) || math.IsNaN(float64(got.Max.Y)) ||
					math.Abs(float64(got.Min.X-want.Min.X)) > 1e-9 ||
					math.Abs(float64(got.Min.Y-want.Min.Y)) > 1e-9 ||
					math.Abs(float64(got.Max.X-want.Max.X)) > 1e-9 ||
					math.Abs(float64(got.Max.Y-want.Max.Y)) > 1e-9 {
					t.Errorf("glyph box does not match drawn label bounds:\ngot: %v\nwant: %v", got, want)
				}
			})
		}
	}
}

func TestCenteredAxisLabelPadding(t *testing.T) {
	for _, horizontal := range []bool{true, false} {
		t.Run(fmt.Sprintf("horizontal=%t", horizontal), func(t *testing.T) {
			p := New()
			p.X.Tick.Marker = ConstantTicks(nil)
			p.Y.Tick.Marker = ConstantTicks(nil)
			axis := &p.Y
			if horizontal {
				axis = &p.X
			}
			axis.Min, axis.Max = 1, 1e6
			axis.Scale = LogScale{}
			axis.Label.Text = "WWWWWWWWWWWWWWWW"
			handler := &labelBoundsHandler{Handler: axis.Label.TextStyle.Handler}
			axis.Label.TextStyle.Handler = handler
			// The label fits in the canvas with 20 points to spare.
			side := axis.Label.TextStyle.Width(axis.Label.Text) + vg.Points(20)
			p.Draw(draw.Canvas{
				Canvas:    new(recorder.Canvas),
				Rectangle: vg.Rectangle{Max: vg.Point{X: side, Y: side}},
			})
			min, max := handler.bounds.Min.Y, handler.bounds.Max.Y
			if horizontal {
				min, max = handler.bounds.Min.X, handler.bounds.Max.X
			}
			if min < 0 || max > side {
				t.Errorf("label clipped: bounds [%v, %v], canvas [0, %v]", min, max, side)
			}
		})
	}
}
