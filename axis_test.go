// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"math"
	"reflect"
	"strconv"
	"testing"
)

func TestAxisSmallTick(t *testing.T) {
	d := DefaultTicks{}
	for _, test := range []struct {
		min, max float64
		format   func(v float64, prec int) string
		want     []string
	}{
		{
			min:    -1.9846500878911073,
			max:    0.4370974820125605,
			format: DefaultTickFormat,
			want:   []string{"-1.6", "-0.8", "0"},
		},
		{
			min:    -1.985e-15,
			max:    0.4371e-15,
			format: DefaultTickFormat,
			want:   []string{"-1.6e-15", "-8e-16", "0"},
		},
		{
			min:    -1.985e15,
			max:    0.4371e15,
			format: DefaultTickFormat,
			want:   []string{"-1.6e+15", "-8e+14", "0"},
		},
		{
			min:    math.MaxFloat64 / 4,
			max:    math.MaxFloat64 / 3,
			format: DefaultTickFormat,
			want:   []string{"4.8e+307", "5.2e+307", "5.6e+307"},
		},
		{
			min:    0.00010,
			max:    0.00015,
			format: DefaultTickFormat,
			want:   []string{"0.0001", "0.00011", "0.00012", "0.00013", "0.00014"},
		},
		{
			min: 0.0001,
			max: 0.0005,
			format: func(v float64, prec int) string {
				return strconv.FormatFloat(v, 'e', 1, 64)
			},
			want: []string{"1.0e-04", "2.0e-04", "3.0e-04", "4.0e-04", "5.0e-04"},
		},
	} {
		ticks := d.Ticks(test.min, test.max, test.format)
		got := labelsOf(ticks)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("tick labels mismatch:\ngot: %q\nwant:%q", got, test.want)
		}
	}
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
