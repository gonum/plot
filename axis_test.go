// Copyright ©2015 The gonum Authors. All rights reserved.
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
		min:        -1.985e-15,
		max:        0.4371e-15,
		wantValues: []float64{-1.985e-15, -7.739500000000001e-16, 4.3709999999999994e-16},
		wantLabels: []string{"-1.985e-15", "-7.7395e-16", "4.371e-16"},
	},
	{
		min:        -1.985e15,
		max:        0.4371e15,
		wantValues: []float64{-2e+15, -1e+15, 0},
		wantLabels: []string{"-2e+15", "-1e+15", "0"},
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
		wantValues: []float64{0.0001, 0.000125, 0.00015000000000000001},
		wantLabels: []string{"0.0001", "0.000125", "0.00015"},
	},
	{
		min:        555.6545,
		max:        21800.9875,
		wantValues: []float64{0, 10000, 20000},
		wantLabels: []string{"0", "1e+04", "2e+04"},
	},
	{
		min:        555.6545,
		max:        27800.9875,
		wantValues: []float64{0, 10000, 20000, 30000},
		wantLabels: []string{"0", "1e+04", "2e+04", "3e+04"},
	},
	{
		min:        55.6545,
		max:        1555.9875,
		wantValues: []float64{0, 500, 1000, 1500},
		wantLabels: []string{"0", "5e+02", "1e+03", "2e+03"},
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
	for _, test := range axisSmallTickTests[:1] {
		ticks := d.Ticks(test.min, test.max)
		gotLabels := labelsOf(ticks)
		gotValues := valuesOf(ticks)
		if !reflect.DeepEqual(gotValues, test.wantValues) {
			t.Errorf("tick values mismatch:\ngot: %v\nwant:%v", gotValues, test.wantValues)
		}
		if !reflect.DeepEqual(gotLabels, test.wantLabels) {
			t.Errorf("tick labels mismatch:\ngot: %q\nwant:%q", gotLabels, test.wantLabels)
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
