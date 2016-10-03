// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package plotutil

import (
	"reflect"
	"testing"
)

var tests = []struct {
	x, y, z           []float64
	xylen, xyzlen     int
	xypanic, xyzpanic string
}{
	{
		x:        nil,
		y:        []float64{0},
		z:        []float64{0},
		xylen:    0,
		xyzlen:   0,
		xypanic:  "x is an empty or nil slice",
		xyzpanic: "x is an empty or nil slice",
	},
	{
		x:        []float64{0},
		y:        nil,
		z:        []float64{0},
		xylen:    0,
		xyzlen:   0,
		xypanic:  "y is an empty or nil slice",
		xyzpanic: "y is an empty or nil slice",
	},
	{
		x:        []float64{0},
		y:        []float64{0},
		z:        nil,
		xylen:    1,
		xyzlen:   0,
		xypanic:  "",
		xyzpanic: "z is an empty or nil slice",
	},
	{
		x:        []float64{},
		y:        []float64{0},
		z:        []float64{0},
		xylen:    0,
		xyzlen:   0,
		xypanic:  "x is an empty or nil slice",
		xyzpanic: "x is an empty or nil slice",
	},
	{
		x:        []float64{0},
		y:        []float64{},
		z:        []float64{0},
		xylen:    0,
		xyzlen:   0,
		xypanic:  "y is an empty or nil slice",
		xyzpanic: "y is an empty or nil slice",
	},
	{
		x:        []float64{0},
		y:        []float64{0},
		z:        []float64{},
		xylen:    1,
		xyzlen:   0,
		xypanic:  "",
		xyzpanic: "z is an empty or nil slice",
	},
	{
		x:        []float64{0},
		y:        []float64{0},
		z:        []float64{0},
		xylen:    1,
		xyzlen:   1,
		xypanic:  "",
		xyzpanic: "",
	},
	{
		x:        genSlice(5),
		y:        genSlice(5),
		z:        genSlice(5),
		xylen:    5,
		xyzlen:   5,
		xypanic:  "",
		xyzpanic: "",
	},
	{
		x:        genSlice(4),
		y:        genSlice(5),
		z:        genSlice(5),
		xylen:    4,
		xyzlen:   4,
		xypanic:  "",
		xyzpanic: "",
	},
	{
		x:        genSlice(5),
		y:        genSlice(4),
		z:        genSlice(5),
		xylen:    4,
		xyzlen:   4,
		xypanic:  "",
		xyzpanic: "",
	},
	{
		x:        genSlice(5),
		y:        genSlice(5),
		z:        genSlice(4),
		xylen:    5,
		xyzlen:   4,
		xypanic:  "",
		xyzpanic: "",
	},
}

func genSlice(len int) []float64 {
	sl := make([]float64, len)
	for i := range sl {
		sl[i] = float64(i)
	}
	return sl
}

func TestXYer(t *testing.T) {
	for tnum, tst := range tests {
		err := func() (err string) {
			defer func() {
				if r := recover(); r != nil {
					err = r.(string)
				}
			}()
			res := ZipXY(tst.x, tst.y)
			xy, ok := res.(*xyer)

			if !ok {
				t.Errorf("Test: %v Incorrect type got: %T expected: *xyer", tnum, res)
				return ""
			}
			if !reflect.DeepEqual(*xy, xyer{tst.x, tst.y, tst.xylen}) {
				t.Errorf("Test: %v XYer test failed", tnum)
			}
			return ""
		}()
		if err != tst.xypanic {
			t.Errorf("Test: %v Incorrect panic got: %q expected: %q", tnum, err, tst.xypanic)
		}
	}
}

func TestXYZer(t *testing.T) {
	for tnum, tst := range tests {
		err := func() (err string) {
			defer func() {
				if r := recover(); r != nil {
					err = r.(string)
				}
			}()
			res := ZipXYZ(tst.x, tst.y, tst.z)
			xyz, ok := res.(*xyzer)

			if !ok {
				t.Errorf("Test: %v Incorrect type found (expected *xyzer) %T", tnum, res)
				return ""
			}
			if !reflect.DeepEqual(*xyz, xyzer{tst.x, tst.y, tst.z, tst.xyzlen}) {
				t.Errorf("Test: %v XYZer test failed", tnum)
			}
			return ""
		}()
		if err != tst.xyzpanic {
			t.Errorf("Test: %v Incorrect panic got: %q expected: %q", tnum, err, tst.xyzpanic)
		}
	}
}
