// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"math"
	"reflect"
	"testing"
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
