// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"reflect"
	"testing"
)

func TestAxisSmallTick(t *testing.T) {
	d := DefaultTicks{}
	for _, test := range []struct {
		min, max float64
		want     []string
	}{
		{
			min:  -1.9846500878911073,
			max:  0.4370974820125605,
			want: []string{"-1.6", "-0.8", "0"},
		},
		//TODO(kortschak) Fix precision to make these tests pass.
		/*{
			min:  -1.985e-15,
			max:  0.4371e-15,
			want: []string{"-1.6e-15", "-8e-16", "0"},
		},*/
		{
			min:  -1.985e15,
			max:  0.4371e15,
			want: []string{"-1.6e+15", "-8e+14", "0"},
		},
		/*{
			min:  math.MaxFloat64 / 4,
			max:  math.MaxFloat64 / 3,
			want: []string{"4.8e+307", "5.2e+307", "5.6e+307"},
		},*/
		{
			min:  0.00010,
			max:  0.00015,
			want: []string{"0.0001", "0.00011", "0.00012", "0.00013", "0.00014"},
		},
		{
			min:  555.6545,
			max:  21800.9875,
			want: []string{"6000", "12000", "18000"},
		},
		{
			min:  555.6545,
			max:  27800.9875,
			want: []string{"8000", "16000", "24000"},
		},
		{
			min:  55.6545,
			max:  1555.9875,
			want: []string{"500", "1000", "1500"},
		},
	} {
		ticks := d.Ticks(test.min, test.max)
		got := labelsOf(ticks)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("tick labels mismatch:\ngot: %q\nwant:%q", got, test.want)
		}
	}
}

func TestPreciseMajorTicks(t *testing.T) {
	p := PreciseTicks{}
	for _, test := range []struct {
		min, max  float64
		valueWant []float64
		labelWant []string
	}{
		{
			min:       3.096916 - 0.125,
			max:       3.096916 + 0.125,
			valueWant: []float64{3., 3.1, 3.2},
			labelWant: []string{"3", "3.1", "3.2"},
		},
	} {
		ticks := p.Ticks(test.min, test.max)
		labelGot := labelsOf(ticks)
		valueGot := valuesOf(ticks)
		if !reflect.DeepEqual(labelGot, test.labelWant) {
			t.Errorf("tick labels mismatch:\ngot: %q\nwant:%q", labelGot, test.labelWant)
		}
		if !reflect.DeepEqual(valueGot, test.valueWant) {
			t.Errorf("tick values mismatch:\ngot: %q\nwant:%q", valueGot, test.valueWant)
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
