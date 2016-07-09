// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moreland

import (
	"image/color"
	"testing"
)

// TestRgb_sRGBA tests the conversion from linear RGB space to sRGB space.
// The expected output value is from
// http://www.kennethmoreland.com/color-maps/DivergingColorMapWorkshop.xls
func TestRgb_sRGBA(t *testing.T) {
	testCases := []struct {
		l rgb
		s sRGBA
	}{
		{
			rgb{R: 0.015299702, G: 0.015299702, B: 0.015299702},
			sRGBA{R: 0.1298716701086684, G: 0.1298716701086684, B: 0.1298716701086684},
		},
	}
	for i, tc := range testCases {
		result := tc.l.sRGBA(0)
		if result != tc.s {
			t.Errorf("case %d: have %+v, want %+v", i, result, tc.s)
		}
	}
}

// TestSRGBa_rgb tests the conversion from sRGB space to linear RGB space.
// The expected output values are from
// http://www.kennethmoreland.com/color-maps/DivergingColorMapWorkshop.xls
func TestSRGBA_rgb(t *testing.T) {
	testCases := []struct {
		s sRGBA
		l rgb
	}{
		{
			sRGBA{R: 0.735356983, G: 0.735356983, B: 0.735356983},
			rgb{R: 0.499999999920366, G: 0.499999999920366, B: 0.499999999920366},
		},
		{
			sRGBA{R: 0.01292, G: 0.01292, B: 0.01292},
			rgb{R: 0.001, G: 0.001, B: 0.001},
		},
		{
			sRGBA{R: 0.759704028, G: 0.162897038, B: 0.206033415},
			rgb{R: 0.5377665307661512, G: 0.022698506403451876, B: 0.035015856125996676},
		},
	}
	for i, tc := range testCases {
		result := tc.s.rgb()
		if result != tc.l {
			t.Errorf("case %d: have %+v, want %+v", i, result, tc.l)
		}
	}
}

// TestCieXYZ_rgb tests the conversion back and forth between
// CIE XYZ space and linear RGB space.
// The expected output values are from
// http://www.kennethmoreland.com/color-maps/DivergingColorMapWorkshop.xls
func TestCieXYZ_rgb(t *testing.T) {
	xyz := cieXYZ{X: 0.128392403, Y: 0.128221351, Z: 0.408477452}
	result := xyz.rgb()
	want := rgb{R: 0.015299702837399953, G: 0.1330700251971, B: 0.4127549680071}
	if result != want {
		t.Errorf("have %+v, want %+v", result, want)
	}

	lrgb := rgb{R: 0.28909265477940005, G: 0.0663313933285, B: 0.0500602839142}
	xyz = cieXYZ{X: 0.151975056, Y: 0.112509738, Z: 0.061066471}
	if xyz.rgb() != lrgb {
		t.Errorf("rgb: have %+v, want %+v", xyz.rgb(), lrgb)
	}
	xyz = cieXYZ{X: 0.1519777983318093, Y: 0.11251566341324888, Z: 0.061068490182446714}
	if lrgb.cieXYZ() != xyz {
		t.Errorf("xyz: have %+v, want %+v", lrgb.cieXYZ(), xyz)
	}
}

// TestCieLAB_cieXYZ tests the conversion from CIE LAB space to CIE XYZ space.
// The expected output values are from
// http://www.kennethmoreland.com/color-maps/DivergingColorMapWorkshop.xls
func TestCieLAB_cieXYZ(t *testing.T) {
	lab := cieLAB{L: 42.49401592, A: 4.416911613, B: -43.38526532}
	result := lab.cieXYZ()
	want := cieXYZ{X: 0.12838835051807143, Y: 0.12822135121812256, Z: 0.40841368569543157}
	if result != want {
		t.Errorf("have %+v, want %+v", result, want)
	}
}

// TestCieLAB_cieXYZ tests the conversion from CIE LAB space to CIE XYZ space.
// The expected output values are from
// http://www.kennethmoreland.com/color-maps/DivergingColorMapWorkshop.xls
func TestMsh_cieLAB(t *testing.T) {
	c := msh{M: 80, S: 1.08, H: -1.1}
	result := c.cieLAB()
	want := cieLAB{L: 37.7062691338992, A: 32.004211237121645, B: -62.88058310076059}
	if result != want {
		t.Errorf("have %+v, want %+v", result, want)
	}
}

// TestCieLAB_sRGBA tests the conversion from CIE LAB space to sRGBA space.
// The expected output values are from
// http://www.kennethmoreland.com/color-maps/DivergingColorMapWorkshop.xls
func TestCieLAB_sRGBA(t *testing.T) {
	testCases := []struct {
		l cieLAB
		s sRGBA
	}{
		{
			cieLAB{},
			sRGBA{},
		},
		{
			cieLAB{L: 43.22418447, A: 59.07682101, B: 32.27381441},
			sRGBA{R: 0.7596910553350515, G: 0.16292472671190056, B: 0.20600836034382436},
		},
	}
	for i, tc := range testCases {
		result := tc.l.sRGBA(0)
		if result != tc.s {
			t.Errorf("case %d: have %+v, want %+v", i, result, tc.s)
		}
	}
}

func TestHueTwist(t *testing.T) {
	// The expected output values are from
	// http://www.kennethmoreland.com/color-maps/DivergingColorMapWorkshop.xls
	if hueTwist(msh{M: 80, S: 1.08, H: -1.1}, 88) != -0.5611585624524025 {
		t.Errorf("hueTwist(80, 1.08, -1.1), 88 should equal -0.561158562 but equals %g",
			hueTwist(msh{M: 80, S: 1.08, H: -1.1}, 88))
	}
}

// TestCieXYZ_cieLAB tests the conversion from CIE XYZ space to CIE LAB space.
// The expected output values are from
// http://www.kennethmoreland.com/color-maps/DivergingColorMapWorkshop.xls
func TestCieXYZ_cieLAB(t *testing.T) {
	xyz := cieXYZ{X: 0.151975056, Y: 0.112509738, Z: 0.061066471}
	lab := cieLAB{L: 40.00000000055783, A: 30.000000104296763, B: 19.99999996294335}
	if xyz.cieLAB() != lab {
		t.Errorf("lab: have %+v, want %+v", xyz.cieLAB(), lab)
	}
	xyz = cieXYZ{X: 0.15197025931227778, Y: 0.11250973800000005, Z: 0.06105693812573921}
	if lab.cieXYZ() != xyz {
		t.Errorf("xyz: have %+v, want %+v", lab.cieXYZ(), xyz)
	}
}

func TestColorToRGB(t *testing.T) {
	c := color.NRGBA{R: 194, G: 42, B: 53, A: 100}
	rgb := sRGBA{R: 0.7607782101167315, G: 0.1646692607003891, B: 0.20782101167315176, A: 0.39215686274509803}
	if colorTosRGBA(c) != rgb {
		t.Errorf("rgb: have %+v, want %+v", colorTosRGBA(c), rgb)
	}
}

// TestCieLAB_msh tests the conversion from CIE LAB space to MSH space.
// The expected output values are from
// http://www.kennethmoreland.com/color-maps/DivergingColorMapWorkshop.xls
func TestCieLAB_msh(t *testing.T) {
	lab := cieLAB{L: 43.22418447, A: 59.07682101, B: 32.27381441}
	mshVal := msh{M: 80.00000000197056, S: 1.0000000000076632, H: 0.5000000000023601}
	if lab.MSH() != mshVal {
		t.Errorf("msh: have %+v, want %+v", lab.MSH(), mshVal)
	}
}

func TestColorToMSH(t *testing.T) {
	c := color.NRGBA{B: 255, A: 255}
	result := colorToMSH(c)
	// The expected output values are from
	// http://www.kennethmoreland.com/color-maps/DivergingColorMapWorkshop.xls
	want := msh{M: 137.64998152940237, S: 1.333915268336423, H: -0.9374394027523394}
	if result != want {
		t.Errorf("want %+v but have %+v", want, result)
	}
}
