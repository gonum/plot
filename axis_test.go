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
