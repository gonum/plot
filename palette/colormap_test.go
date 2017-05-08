// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package palette

import (
	"image/color"
	"math"
	"strings"
	"testing"
)

func TestColorMap_Range(t *testing.T) {
	p := FromPalette(New(color.White, color.Black))
	if _, err := p.At(0); !strings.Contains(err.Error(), "max == min") {
		t.Errorf("should have 'max == min' error")
	}
	p.SetMax(1)
	vals := []float64{-1, 0, 1, 2, math.Inf(1), math.Inf(-1), math.NaN()}
	errs := []error{ErrUnderflow, nil, nil, ErrOverflow,
		ErrOverflow, ErrUnderflow, ErrNaN}
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
