// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotutil

import (
	"math"

	"github.com/gonum/plot/vg"
)

// phi is the golden ratio, as defined by:
//  https://en.wikipedia.org/wiki/Golden_ratio
var phi = vg.Length(1+math.Sqrt(5)) / 2

// GoldenRatio returns the width and the height such as:
//  width = size
//  height = size/phi
// so that width and height follow the golden proportions.
func GoldenRatio(size vg.Length) (vg.Length, vg.Length) {
	w := size
	h := size / phi
	return w, h
}
