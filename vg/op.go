// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vg

import (
	"image"
	"image/color"
)

type Op interface {
	isOp()
}

type LineWidth struct {
	Width Length
}

type LineDash struct {
	Pattern []Length
	Offset  Length
}

type Color struct {
	Color color.Color
}

type Rotate struct {
	Radians float64
}

type Translate struct {
	Point Point
}

type Scale struct {
	X, Y float64
}

type Push struct{}
type Pop struct{}

type Stroke struct {
	Path Path
}

type Fill struct {
	Path Path
}

type FillString struct {
	Font  Font
	Point Point
	Text  string
}

type DrawImage struct {
	Rect  Rectangle
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
