// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw_test

import (
	"fmt"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func ExampleCrop_splitHorizontal() {
	c := draw.New(vgimg.New(vg.Points(10), vg.Points(16)))

	// Split c along a vertical line centered on the canvas.
	left, right := SplitHorizontal(c, c.Size().X/2)
	fmt.Printf("left:  %#v\n", left.Rectangle)
	fmt.Printf("right: %#v\n", right.Rectangle)

	// Output:
	// left:  vg.Rectangle{Min:vg.Point{X:0, Y:0}, Max:vg.Point{X:5, Y:16}}
	// right: vg.Rectangle{Min:vg.Point{X:5, Y:0}, Max:vg.Point{X:10, Y:16}}
}
