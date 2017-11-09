// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"math"
	"reflect"
	"testing"
)

var talbotLinHanrahanTests = []struct {
	dMin, dMax  float64
	want        int
	containment int

	wantValues    []float64
	wantStep      float64
	wantMagnitude int
}{
	// Expected values confirmed against R reference imlpementation.
	{
		dMin:        -1.9846500878911073,
		dMax:        0.4370974820125605,
		want:        3,
		containment: free,

		wantValues:    []float64{-2, -1, 0},
		wantStep:      1,
		wantMagnitude: 0,
	},
	{
		dMin:        -1.9846500878911073,
		dMax:        0.4370974820125605,
		want:        3,
		containment: containData,

		wantValues:    []float64{-2, -1, 0, 1},
		wantStep:      1,
		wantMagnitude: 0,
	},
	{
		dMin:        -1.985e-15,
		dMax:        0.4371e-15,
		want:        3,
		containment: free,

		wantValues:    []float64{-1.985e-15, -7.739500000000001e-16, 4.3709999999999994e-16},
		wantStep:      1.21105e-15,
		wantMagnitude: -16,
	},
	{
		dMin:        -1.985e-15,
		dMax:        0.4371e-15,
		want:        3,
		containment: containData,

		wantValues:    []float64{-1.985e-15, -7.739500000000001e-16, 4.3709999999999994e-16},
		wantStep:      1.21105e-15,
		wantMagnitude: -16,
	},
	{
		dMin:        -1.985e15,
		dMax:        0.4371e15,
		want:        3,
		containment: free,

		wantValues:    []float64{-2e+15, -1e+15, 0},
		wantStep:      1,
		wantMagnitude: 15,
	},
	{
		dMin:        -1.985e15,
		dMax:        0.4371e15,
		want:        3,
		containment: containData,

		wantValues:    []float64{-2e+15, -1e+15, 0, 1e+15},
		wantStep:      1,
		wantMagnitude: 15,
	},
	{
		dMin:        dlamchP * 20,
		dMax:        dlamchP * 50,
		want:        3,
		containment: free,

		wantValues:    []float64{4.440892098500626e-15, 7.771561172376096e-15, 1.1102230246251565e-14},
		wantStep:      3.3306690738754696e-15,
		wantMagnitude: -15,
	},
	{
		dMin:        dlamchP * 20,
		dMax:        dlamchP * 50,
		want:        3,
		containment: containData,

		wantValues:    []float64{4.440892098500626e-15, 7.771561172376096e-15, 1.1102230246251565e-14},
		wantStep:      3.3306690738754696e-15,
		wantMagnitude: -15,
	},
	{
		dMin:        math.MaxFloat64 / 4,
		dMax:        math.MaxFloat64 / 3,
		want:        3,
		containment: free,

		wantValues:    []float64{4.4942328371557893e+307, 5.243271643348421e+307, 5.992310449541053e+307},
		wantStep:      7.490388061926317e+306,
		wantMagnitude: 307,
	},
	{
		dMin:        math.MaxFloat64 / 4,
		dMax:        math.MaxFloat64 / 3,
		want:        3,
		containment: containData,

		wantValues:    []float64{4.4942328371557893e+307, 5.243271643348421e+307, 5.992310449541053e+307},
		wantStep:      7.490388061926317e+306,
		wantMagnitude: 307,
	},
	{
		dMin:        0.00010,
		dMax:        0.00015,
		want:        3,
		containment: free,

		wantValues:    []float64{0.0001, 0.000125, 0.00015000000000000001},
		wantStep:      2.5,
		wantMagnitude: -5,
	},
	{
		dMin:        0.00010,
		dMax:        0.00015,
		want:        3,
		containment: containData,

		wantValues:    []float64{0.0001, 0.000125, 0.00015000000000000001},
		wantStep:      2.5,
		wantMagnitude: -5,
	},
	{
		dMin:        555.6545,
		dMax:        21800.9875,
		want:        3,
		containment: free,

		wantValues:    []float64{0, 10000, 20000},
		wantStep:      1,
		wantMagnitude: 4,
	},
	{
		dMin:        555.6545,
		dMax:        21800.9875,
		want:        3,
		containment: containData,

		wantValues:    []float64{0, 12000, 24000},
		wantStep:      12,
		wantMagnitude: 3,
	},
	{
		dMin:        555.6545,
		dMax:        27800.9875,
		want:        3,
		containment: free,

		wantValues:    []float64{0, 10000, 20000, 30000},
		wantStep:      1,
		wantMagnitude: 4,
	},
	{
		dMin:        555.6545,
		dMax:        27800.9875,
		want:        3,
		containment: containData,

		wantValues:    []float64{0, 10000, 20000, 30000},
		wantStep:      1,
		wantMagnitude: 4,
	},
	{
		dMin:        55.6545,
		dMax:        1555.9875,
		want:        3,
		containment: free,

		wantValues:    []float64{0, 500, 1000, 1500},
		wantStep:      5,
		wantMagnitude: 2,
	},
	{
		dMin:        55.6545,
		dMax:        1555.9875,
		want:        3,
		containment: containData,

		wantValues:    []float64{0, 800, 1600},
		wantStep:      8,
		wantMagnitude: 2,
	},
	{
		dMin:        3.096916 - 0.125,
		dMax:        3.096916 + 0.125,
		want:        3,
		containment: free,

		wantValues:    []float64{3, 3.1, 3.2},
		wantStep:      1,
		wantMagnitude: -1,
	},
	{
		dMin:        3.096916 - 0.125,
		dMax:        3.096916 + 0.125,
		want:        3,
		containment: containData,

		wantValues:    []float64{2.9499999999999997, 3.0999999999999996, 3.2499999999999996},
		wantStep:      15,
		wantMagnitude: -2,
	},
}

func TestTalbotLinHanrahan(t *testing.T) {
	for _, test := range talbotLinHanrahanTests {
		values, step, _, magnitude := talbotLinHanrahan(test.dMin, test.dMax, test.want, test.containment, nil, nil, nil)
		if !reflect.DeepEqual(values, test.wantValues) {
			t.Errorf("unexpected values for dMin=%g, dMax=%g, want=%d, containment=%d:\ngot: %v\nwant:%v",
				test.dMin, test.dMax, test.want, test.containment, values, test.wantValues)
		}
		if step != test.wantStep {
			t.Errorf("unexpected step for dMin=%g, dMax=%g, want=%d, containment=%d: got:%v want:%v",
				test.dMin, test.dMax, test.want, test.containment, step, test.wantStep)
		}
		if magnitude != test.wantMagnitude {
			t.Errorf("unexpected magnitude for dMin=%g, dMax=%g, want=%d, containment=%t: got:%d want:%d",
				test.dMin, test.dMax, test.want, test.containment, magnitude, test.wantMagnitude)
		}
		if test.containment == containData {
			f := math.Pow10(-magnitude)
			if test.containment == containData && (test.dMin*f < values[0]*f || values[len(values)-1]*f < test.dMax*f) {
				t.Errorf("unexpected values for containment dMin=%g, dMax=%g, want=%d not containment:\ngot: %v\nwant:%v",
					test.dMin, test.dMax, test.want, values, test.wantValues)
			}
		}
	}
}
