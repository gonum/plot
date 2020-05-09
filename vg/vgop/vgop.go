// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgop // import "gonum.org/v1/plot/vg/vgop"

import (
	"image"
	"image/color"

	"gonum.org/v1/plot/vg"
)

type Op interface {
	isOp()
}

type LineWidth struct {
	Width vg.Length
}

type LineDash struct {
	Pattern []vg.Length
	Offset  vg.Length
}

type Color struct {
	Color color.Color
}

type Rotate struct {
	Radians float64
}

type Translate struct {
	Point vg.Point
}

type Scale struct {
	X, Y float64
}

type Push struct{}
type Pop struct{}

type Stroke struct {
	Path vg.Path
}

type Fill struct {
	Path vg.Path
}

type FillString struct {
	Font  vg.Font
	Point vg.Point
	Text  string
}

type DrawImage struct {
	Rect  vg.Rectangle
	Image image.Image
}

func (LineWidth) isOp()  {}
func (LineDash) isOp()   {}
func (Color) isOp()      {}
func (Rotate) isOp()     {}
func (Translate) isOp()  {}
func (Scale) isOp()      {}
func (Push) isOp()       {}
func (Pop) isOp()        {}
func (Stroke) isOp()     {}
func (Fill) isOp()       {}
func (FillString) isOp() {}
func (DrawImage) isOp()  {}

var (
	_ Op = (*LineWidth)(nil)
	_ Op = (*LineDash)(nil)
	_ Op = (*Color)(nil)
	_ Op = (*Rotate)(nil)
	_ Op = (*Translate)(nil)
	_ Op = (*Scale)(nil)
	_ Op = (*Push)(nil)
	_ Op = (*Pop)(nil)
	_ Op = (*Stroke)(nil)
	_ Op = (*Fill)(nil)
	_ Op = (*FillString)(nil)
	_ Op = (*DrawImage)(nil)
)
