// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moreland

import (
	"fmt"
	"image/color"
	"reflect"
	"testing"
)

func TestCreateLuminance(t *testing.T) {
	type testHolder struct {
		controlColors []color.Color
		want          *luminance
		name          string
	}
	tests := []testHolder{
		testHolder{
			name: "BlackBody",
			controlColors: []color.Color{
				color.NRGBA{0, 0, 0, 255},
				color.NRGBA{178, 34, 34, 255},
				color.NRGBA{227, 105, 5, 255},
				color.NRGBA{238, 210, 20, 255},
				color.NRGBA{255, 255, 255, 255},
			},
			want: BlackBody().(*luminance),
		},
		testHolder{
			name: "ExtendedBlackBody",
			controlColors: []color.Color{
				color.NRGBA{0, 0, 0, 255},
				color.NRGBA{0, 24, 168, 255},
				color.NRGBA{99, 0, 228, 255},
				color.NRGBA{220, 20, 60, 255},
				color.NRGBA{255, 117, 56, 255},
				color.NRGBA{238, 210, 20, 255},
				color.NRGBA{255, 255, 255, 255},
			},
			want: ExtendedBlackBody().(*luminance),
		},
		testHolder{
			name: "Kindlmann",
			controlColors: []color.Color{
				color.NRGBA{0, 0, 0, 255},
				color.NRGBA{46, 4, 76, 255},
				color.NRGBA{63, 7, 145, 255},
				color.NRGBA{8, 66, 165, 255},
				color.NRGBA{5, 106, 106, 255},
				color.NRGBA{7, 137, 169, 255},
				color.NRGBA{8, 168, 26, 255},
				color.NRGBA{84, 194, 9, 255},
				color.NRGBA{196, 206, 10, 255},
				color.NRGBA{252, 220, 197, 255},
				color.NRGBA{255, 255, 255, 255},
			},
			want: Kindlmann().(*luminance),
		},
		testHolder{
			name: "ExtendedKindlmann",
			controlColors: []color.Color{
				color.NRGBA{0, 0, 0, 255},
				color.NRGBA{44, 5, 103, 255},
				color.NRGBA{3, 67, 67, 255},
				color.NRGBA{5, 103, 13, 255},
				color.NRGBA{117, 124, 6, 255},
				color.NRGBA{246, 104, 74, 255},
				color.NRGBA{250, 149, 241, 255},
				color.NRGBA{232, 212, 253, 255},
				color.NRGBA{255, 255, 255, 255},
			},
			want: ExtendedKindlmann().(*luminance),
		},
	}
	for _, test := range tests {
		cmap, err := NewLuminance(test.controlColors)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(cmap, test.want) {
			fmt.Printf("%#v\n", cmap)
			t.Errorf("%s: have %#v, want %#v", test.name, cmap, test.want)
		}
	}
}

func TestExtendedBlackBody(t *testing.T) {
	scalars := []float64{0, 0.21873483862751875, 0.34506542513775906, 0.4702980511087303, 0.6517482203230537, 0.8413253643355525, 1}
	want := []color.Color{
		color.NRGBA{0, 0, 0, 255},
		color.NRGBA{0, 24, 168, 255},
		color.NRGBA{99, 0, 228, 255},
		color.NRGBA{220, 20, 60, 255},
		color.NRGBA{255, 117, 56, 255},
		color.NRGBA{238, 210, 20, 255},
		color.NRGBA{255, 255, 255, 255},
	}

	colors := ExtendedBlackBody()
	colors.SetMax(1)

	for i, scalar := range scalars {
		c, err := colors.At(scalar)
		if err != nil {
			t.Fatal(err)
		}
		if !similar(c, want[i], bitTolerance) {
			t.Errorf("color %d: have %+v, want %+v", i, c, want[i])
		}
	}
}

// floatToUint32 converts a float64 in the range [0, 1] to a uint32 in the range
// [0, 0xffff].
func floatToUint32(f float64) uint32 {
	return uint32(f * 0xffff)
}
