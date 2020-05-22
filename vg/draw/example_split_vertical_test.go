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

func ExampleCrop_splitVertical() {
	c := draw.New(vgimg.New(vg.Points(10), vg.Points(16)))

	// Split c along a horizontal line centered on the canvas.
	bottom, top := SplitVertical(c, c.Size().Y/2)
	fmt.Printf("top:    %#v\n", top.Rectangle)
	fmt.Printf("bottom: %#v\n", bottom.Rectangle)

	// Output:
	// top:    vg.Rectangle{Min:vg.Point{X:0, Y:8}, Max:vg.Point{X:10, Y:16}}
	// bottom: vg.Rectangle{Min:vg.Point{X:0, Y:0}, Max:vg.Point{X:10, Y:8}}
}
