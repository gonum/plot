// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw_test

import (
	"fmt"

	"gonum.org/v1/plot/vg/draw"
)

func ExampleCrop_splitVertical() {
	var c draw.Canvas
	// Split c along a horizontal line centered on the canvas.
	bottom, top := SplitHorizontal(c, c.Size().Y/2)
	fmt.Println(bottom.Rectangle.Size(), top.Rectangle.Size())
}
