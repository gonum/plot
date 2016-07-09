// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moreland

import (
	"image/color"
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/gonum/plot/palette"
)

// bitTolerance is the precision of a uint8 when
// expressed as a uint32. This tolerance is used in tests
// where precision can be lost when coverting between
// 8-bit and 32-bit values.
const bitTolerance = 1.0 / 256.0 * 65535.0

func TestInterpolateMSHDiverging(t *testing.T) {
	type test struct {
		start, end                       msh
		convergeM, scalar, convergePoint float64
		result                           msh
	}
	tests := []test{
		test{
			start:         msh{M: 80, S: 1.08, H: -1.1},
			end:           msh{M: 80, S: 1.08, H: 0.5},
			convergeM:     88,
			convergePoint: 0.5,
			scalar:        0.125,
			result:        msh{M: 82, S: 0.81, H: -1.2402896406131008},
		},
		test{
			start:         msh{M: 80, S: 1.08, H: -1.1},
			end:           msh{M: 80, S: 1.08, H: 0.5},
			convergeM:     88,
			convergePoint: 0.5,
			scalar:        0.5,
			result:        msh{M: 88, S: 0, H: 0},
		},
		test{
			start:         msh{M: 80, S: 1.08, H: -1.1},
			end:           msh{M: 80, S: 1.08, H: 0.5},
			convergeM:     88,
			convergePoint: 0.5,
			scalar:        0.75,
			result:        msh{M: 84, S: 0.54, H: 0.7805792812262012},
		},
		test{
			start:         msh{M: 80, S: 1.08, H: -1.1},
			end:           msh{M: 80, S: 1.08, H: 0.5},
			convergeM:     88,
			convergePoint: 0.75,
			scalar:        0.7499999999999999,
			result:        msh{M: 88, S: 1.1990408665951691e-16, H: -1.6611585624524023},
		},
		test{
			start:         msh{M: 80, S: 1.08, H: -1.1},
			end:           msh{M: 80, S: 1.08, H: 0.5},
			convergeM:     88,
			convergePoint: 0.75,
			scalar:        0.75,
			result:        msh{M: 88, S: 0, H: 0},
		},
		test{
			start:         msh{M: 80, S: 1.08, H: -1.1},
			end:           msh{M: 80, S: 1.08, H: 0.5},
			convergeM:     88,
			convergePoint: 0.75,
			scalar:        0.7500000000000001,
			result:        msh{M: 88, S: 2.3980817331903383e-16, H: 1.0611585624524023},
		},
	}
	for i, test := range tests {
		p := newSmoothDiverging(test.start, test.end, test.convergeM)
		p.SetMin(0)
		p.SetMax(1)
		result := p.(*smoothDiverging).interpolateMSHDiverging(test.scalar, test.convergePoint)
		if result != test.result {
			t.Errorf("test %d: expected %v; got %v", i, test.result, result)
		}
	}
}

func TestSmoothBlueRed(t *testing.T) {
	p := SmoothBlueRed()
	wantP := []color.NRGBA{
		{59, 76, 192, 255},
		{68, 90, 204, 255},
		{77, 104, 215, 255},
		{87, 117, 225, 255},
		{98, 130, 234, 255},
		{108, 142, 241, 255},
		{119, 154, 247, 255},
		{130, 165, 251, 255},
		{141, 176, 254, 255},
		{152, 185, 255, 255},
		{163, 194, 255, 255},
		{174, 201, 253, 255},
		{184, 208, 249, 255},
		{194, 213, 244, 255},
		{204, 217, 238, 255},
		{213, 219, 230, 255},
		{221, 221, 221, 255},
		{229, 216, 209, 255},
		{236, 211, 197, 255},
		{241, 204, 185, 255},
		{245, 196, 173, 255},
		{247, 187, 160, 255},
		{247, 177, 148, 255},
		{247, 166, 135, 255},
		{244, 154, 123, 255},
		{241, 141, 111, 255},
		{236, 127, 99, 255},
		{229, 112, 88, 255},
		{222, 96, 77, 255},
		{213, 80, 66, 255},
		{203, 62, 56, 255},
		{192, 40, 47, 255},
		{180, 4, 38, 255},
	}
	c := p.Palette(33)
	if len(c.Colors()) != len(wantP) {
		t.Errorf("length doesn't match: %d != %d", len(c.Colors()), len(wantP))
	}
	for i, c := range c.Colors() {
		w := wantP[i]
		if !similar(w, c, bitTolerance) {
			t.Errorf("%d: want %+v but have %+v", i, w, c)
		}
	}
}

func TestSmoothCoolWarm(t *testing.T) {
	type test struct {
		start [3]float64
		f     func(int) palette.Palette
		end   [3]float64
	}
	tests := []test{
		{[3]float64{0.230, 0.299, 0.754}, SmoothBlueRed().Palette, [3]float64{0.706, 0.016, 0.150}},
		{[3]float64{0.436, 0.308, 0.631}, SmoothPurpleOrange().Palette, [3]float64{0.759, 0.334, 0.046}},
		{[3]float64{0.085, 0.532, 0.201}, SmoothGreenPurple().Palette, [3]float64{0.436, 0.308, 0.631}},
		{[3]float64{0.217, 0.525, 0.910}, SmoothBlueTan().Palette, [3]float64{0.677, 0.492, 0.093}},
		{[3]float64{0.085, 0.532, 0.201}, SmoothGreenRed().Palette, [3]float64{0.758, 0.214, 0.233}},
	}
	midPoint := [3]float64{0.865, 0.865, 0.865}

	for i, test := range tests {
		c := test.f(3).Colors()
		testRGB(t, i, "start", c[0], test.start)
		testRGB(t, i, "mid", c[1], midPoint)
		testRGB(t, i, "end", c[2], test.end)
	}
}

func fracToByte(v float64) uint8 {
	return uint8(v*255 + 0.5)
}

func testRGB(t *testing.T, i int, label string, c1 color.Color, c2 [3]float64) {
	c3 := color.NRGBA{
		R: fracToByte(c2[0]),
		G: fracToByte(c2[1]),
		B: fracToByte(c2[2]),
		A: 255,
	}
	if !similar(c1, c3, bitTolerance) {
		t.Errorf("%d %s: want %+v but have %+v", i, label, c1, c3)
	}
}

func TestSmoothDiverging_At(t *testing.T) {

	start := msh{M: 80, S: 1.08, H: -1.1}
	end := msh{M: 80, S: 1.08, H: 0.5}
	p := newSmoothDiverging(start, end, 88)
	p.SetMax(2)
	p.SetMin(-1)
	scalar := -1 + 3*0.125
	rgb, err := p.At(scalar)
	if err != nil {
		t.Error(err)
	}
	// The expected output values are from
	// http://www.kennethmoreland.com/color-maps/DivergingColorMapWorkshop.xls
	want := color.NRGBA{R: 98, G: 130, B: 234, A: 255}
	if !similar(want, rgb, bitTolerance) {
		t.Errorf("have %+v, want %+v", rgb, want)
	}
}

func BenchmarkSmoothDiverging_At(b *testing.B) {
	p := SmoothBlueRed()
	p.SetMax(1.0000000001)
	rand.Seed(1)
	for i := 0; i < b.N; i++ {
		if _, err := p.At(rand.Float64()); err != nil {
			b.Fatal(err)
		}
	}
}

func TestSmoothDiverging_Range(t *testing.T) {
	p := SmoothBlueRed()
	if _, err := p.At(0); !strings.Contains(err.Error(), "max == min") {
		t.Errorf("should have 'max == min' error")
	}
	p.SetMax(1)
	vals := []float64{-1, 0, 1, 2, math.Inf(1), math.Inf(-1), math.NaN()}
	errs := []error{palette.ErrUnderflow, nil, nil, palette.ErrOverflow,
		palette.ErrOverflow, palette.ErrUnderflow, palette.ErrNaN}
	for i, v := range vals {
		_, err := p.At(v)
		wantErr := errs[i]
		if wantErr == nil && err != nil {
			t.Errorf("val %g: want no error but have %v", v, err)
		} else if wantErr != nil && err == nil {
			t.Errorf("val %g: want error but have no error", v)
		} else if wantErr != err {
			t.Errorf("val %g: want error %v but have %v", v, wantErr, err)
		}
	}
	p.SetMin(2)
	if _, err := p.At(0); !strings.Contains(err.Error(), "< min") {
		t.Errorf("should have 'max < min' error")
	}
}

// similar compares whether the fields of a and b are within the
// specified tolerance of each other.
func similar(a, b color.Color, tolerance float64) bool {
	aR, aG, aB, aA := a.RGBA()
	bR, bG, bB, bA := b.RGBA()
	if math.Abs(float64(bR)-float64(aR)) > tolerance ||
		math.Abs(float64(bG)-float64(aG)) > tolerance ||
		math.Abs(float64(bB)-float64(aB)) > tolerance ||
		math.Abs(float64(bA)-float64(aA)) > tolerance {
		return false
	}
	return true
}
