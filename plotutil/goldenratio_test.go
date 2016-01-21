// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotutil

import (
	"testing"

	"github.com/gonum/plot/vg"
)

func TestGoldenRatio(t *testing.T) {
	sz := 10 * vg.Centimeter
	w, h := GoldenRatio(sz)
	o := w/h*w/h - w/h - 1
	if o != 0 {
		t.Fatalf(
			"w and h do not follow the golden proportions. got=%v want=0.\n",
			o,
		)
	}
}
