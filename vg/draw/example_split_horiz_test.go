// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw_test

import (
	"fmt"

	"gonum.org/v1/plot/vg/draw"
)

func ExampleCrop_splitHorizontal() {
	var c draw.Canvas
	// Split c along a vertical line centered on the canvas.
	left, right := SplitHorizontal(c, c.Size().X/2)
	fmt.Println(left.Rectangle.Size(), right.Rectangle.Size())
}
