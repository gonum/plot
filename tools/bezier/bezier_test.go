// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bezier

import (
	"math"
	"reflect"
	"testing"

	"gonum.org/v1/plot/vg"
)

const tol = 1e-12

func approxEqual(a, b vg.Point, tol float64) bool {
	return (math.Abs(float64(a.X-b.X)) <= tol || (math.IsNaN(float64(a.X)) && math.IsNaN(float64(b.X)))) &&
		(math.Abs(float64(a.Y-b.Y)) <= tol || (math.IsNaN(float64(a.Y)) && math.IsNaN(float64(b.Y))))
}

func TestNew(t *testing.T) {
	for i, test := range []struct {
		ctrls []vg.Point
		curve Curve
	}{
		{
			ctrls: nil,
			curve: nil,
		},
		{
			ctrls: []vg.Point{{X: 1, Y: 2}, {X: 3, Y: 4}, {X: 5, Y: 6}, {X: 7, Y: 8}},
			curve: Curve{
				{Point: vg.Point{X: 1, Y: 2}, Control: vg.Point{X: 1, Y: 2}},
				{Point: vg.Point{X: 3, Y: 4}, Control: vg.Point{X: 9, Y: 12}},
				{Point: vg.Point{X: 5, Y: 6}, Control: vg.Point{X: 15, Y: 18}},
				{Point: vg.Point{X: 7, Y: 8}, Control: vg.Point{X: 7, Y: 8}},
			},
		},
		{
			ctrls: []vg.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 0}},
			curve: Curve{
				{Point: vg.Point{X: 0, Y: 0}, Control: vg.Point{X: 0, Y: 0}},
				{Point: vg.Point{X: 0, Y: 1}, Control: vg.Point{X: 0, Y: 3}},
				{Point: vg.Point{X: 1, Y: 1}, Control: vg.Point{X: 3, Y: 3}},
				{Point: vg.Point{X: 1, Y: 0}, Control: vg.Point{X: 1, Y: 0}},
			},
		},
		{
			ctrls: []vg.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}},
			curve: Curve{
				{Point: vg.Point{X: 0, Y: 0}, Control: vg.Point{X: 0, Y: 0}},
				{Point: vg.Point{X: 0, Y: 1}, Control: vg.Point{X: 0, Y: 3}},
				{Point: vg.Point{X: 1, Y: 0}, Control: vg.Point{X: 3, Y: 0}},
				{Point: vg.Point{X: 1, Y: 1}, Control: vg.Point{X: 1, Y: 1}},
			},
		},
	} {
		bc := New(test.ctrls...)
		if !reflect.DeepEqual(bc, test.curve) {
			t.Errorf("unexpected result for test %d:\ngot: %+v\nwant:%+v", i, bc, test.ctrls)
		}
	}
}

func TestPoint(t *testing.T) {
	type tPoints []struct {
		t     float64
		point vg.Point
	}
	for i, test := range []struct {
		ctrls []vg.Point
		tPoints
	}{
		{
			ctrls: []vg.Point{{X: 1, Y: 2}, {X: 3, Y: 4}, {X: 5, Y: 6}, {X: 7, Y: 8}},
			tPoints: tPoints{
				{t: 0, point: vg.Point{X: 1, Y: 2}},
				{t: 0.1, point: vg.Point{X: 1.6, Y: 2.6}},
				{t: 0.2, point: vg.Point{X: 2.2, Y: 3.2}},
				{t: 0.3, point: vg.Point{X: 2.8, Y: 3.8}},
				{t: 0.4, point: vg.Point{X: 3.4, Y: 4.4}},
				{t: 0.5, point: vg.Point{X: 4, Y: 5}},
				{t: 0.6, point: vg.Point{X: 4.6, Y: 5.6}},
				{t: 0.7, point: vg.Point{X: 5.2, Y: 6.2}},
				{t: 0.8, point: vg.Point{X: 5.8, Y: 6.8}},
				{t: 0.9, point: vg.Point{X: 6.4, Y: 7.4}},
				{t: 1, point: vg.Point{X: 7, Y: 8}},
			},
		},
		{
			ctrls: []vg.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 0}},
			tPoints: tPoints{
				{t: 0, point: vg.Point{X: 0, Y: 0}},
				{t: 0.1, point: vg.Point{X: 0.028, Y: 0.27}},
				{t: 0.2, point: vg.Point{X: 0.104, Y: 0.48}},
				{t: 0.3, point: vg.Point{X: 0.216, Y: 0.63}},
				{t: 0.4, point: vg.Point{X: 0.352, Y: 0.72}},
				{t: 0.5, point: vg.Point{X: 0.5, Y: 0.75}},
				{t: 0.6, point: vg.Point{X: 0.648, Y: 0.72}},
				{t: 0.7, point: vg.Point{X: 0.784, Y: 0.63}},
				{t: 0.8, point: vg.Point{X: 0.896, Y: 0.48}},
				{t: 0.9, point: vg.Point{X: 0.972, Y: 0.27}},
				{t: 1, point: vg.Point{X: 1, Y: 0}},
			},
		},
		{
			ctrls: []vg.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}},
			tPoints: tPoints{
				{t: 0, point: vg.Point{X: 0, Y: 0}},
				{t: 0.1, point: vg.Point{X: 0.028, Y: 0.244}},
				{t: 0.2, point: vg.Point{X: 0.104, Y: 0.392}},
				{t: 0.3, point: vg.Point{X: 0.216, Y: 0.468}},
				{t: 0.4, point: vg.Point{X: 0.352, Y: 0.496}},
				{t: 0.5, point: vg.Point{X: 0.5, Y: 0.5}},
				{t: 0.6, point: vg.Point{X: 0.648, Y: 0.504}},
				{t: 0.7, point: vg.Point{X: 0.784, Y: 0.532}},
				{t: 0.8, point: vg.Point{X: 0.896, Y: 0.608}},
				{t: 0.9, point: vg.Point{X: 0.972, Y: 0.756}},
				{t: 1, point: vg.Point{X: 1, Y: 1}},
			},
		},
	} {
		bc := New(test.ctrls...)
		for j, tPoint := range test.tPoints {
			got := bc.Point(tPoint.t)
			want := test.tPoints[j].point
			if !approxEqual(got, want, tol) {
				t.Errorf("unexpected point for test %d part %d %+v: got:%+v want:%+v", i, j, test.ctrls, got, want)
			}
		}
	}
}

func TestCurve(t *testing.T) {
	for i, test := range []struct {
		ctrls  []vg.Point
		points []vg.Point
	}{
		{
			ctrls: []vg.Point{{X: 1, Y: 2}, {X: 3, Y: 4}, {X: 5, Y: 6}, {X: 7, Y: 8}},
			points: []vg.Point{
				{X: 1, Y: 2},
				{X: 1.6, Y: 2.6},
				{X: 2.2, Y: 3.2},
				{X: 2.8, Y: 3.8},
				{X: 3.4, Y: 4.4},
				{X: 4, Y: 5},
				{X: 4.6, Y: 5.6},
				{X: 5.2, Y: 6.2},
				{X: 5.8, Y: 6.8},
				{X: 6.4, Y: 7.4},
				{X: 7, Y: 8},
			},
		},
		{
			ctrls: []vg.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 0}},
			points: []vg.Point{
				{X: 0, Y: 0},
				{X: 0.028, Y: 0.27},
				{X: 0.104, Y: 0.48},
				{X: 0.216, Y: 0.63},
				{X: 0.352, Y: 0.72},
				{X: 0.5, Y: 0.75},
				{X: 0.648, Y: 0.72},
				{X: 0.784, Y: 0.63},
				{X: 0.896, Y: 0.48},
				{X: 0.972, Y: 0.27},
				{X: 1, Y: 0},
			},
		},
		{
			ctrls: []vg.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}},
			points: []vg.Point{
				{X: 0, Y: 0},
				{X: 0.028, Y: 0.244},
				{X: 0.104, Y: 0.392},
				{X: 0.216, Y: 0.468},
				{X: 0.352, Y: 0.496},
				{X: 0.5, Y: 0.5},
				{X: 0.648, Y: 0.504},
				{X: 0.784, Y: 0.532},
				{X: 0.896, Y: 0.608},
				{X: 0.972, Y: 0.756},
				{X: 1, Y: 1},
			},
		},
		{
			ctrls:  []vg.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}},
			points: []vg.Point{},
		},
		{
			ctrls: []vg.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}},
			points: []vg.Point{
				{X: vg.Length(math.NaN()), Y: vg.Length(math.NaN())},
			},
		}, {
			ctrls: []vg.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}},
			points: []vg.Point{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
			},
		},
	} {
		bc := New(test.ctrls...).Curve(make([]vg.Point, len(test.points)))
		for j, got := range bc {
			want := test.points[j]
			if !approxEqual(got, want, tol) {
				t.Errorf("unexpected point for test %d part %d %+v: got:%+v want:%+v", i, j, test.ctrls, got, want)
			}
		}
	}
}
