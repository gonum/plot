// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !js

package vg_test

import (
	"testing"

	"github.com/gonum/plot/vg"
)

func TestVGFONTPATH(t *testing.T) {
	if len(vg.FontDirs) == 0 {
		t.Fatalf("zero length FontDirs")
	}
}
