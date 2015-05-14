// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"math"
	"testing"
)

func TestAxisSmallTick(t *testing.T) {
	d := DefaultTicks{}
	for _, test := range []struct {
		Min, Max float64
		Labels   []string
	}{
		{
			Min:    -1.9846500878911073,
			Max:    0.4370974820125605,
			Labels: []string{"-1.6", "-0.8", "0"},
		},
		{
			Min:    -1.985e-15,
			Max:    0.4371e-15,
			Labels: []string{"-1.6e-15", "-8e-16", "0"},
		},
		{
			Min:    -1.985e15,
			Max:    0.4371e15,
			Labels: []string{"-1.6e+15", "-8e+14", "0"},
		},
		{
			Min:    math.MaxFloat64 / 4,
			Max:    math.MaxFloat64 / 3,
			Labels: []string{"4.8e+307", "5.2e+307", "5.6e+307"},
		},
	} {
		ticks := d.Ticks(test.Min, test.Max)
		var count int
		for _, tick := range ticks {
			if tick.Label != "" {
				if test.Labels[count] != tick.Label {
					t.Error("Ticks mismatch: Want", test.Labels[count], ", got", tick.Label)
				}
				count++
			}
		}
		if count != len(test.Labels) {
			t.Errorf("Too many tick labels")
		}
	}
}
