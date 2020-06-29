// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vggio provides a vg.Canvas implementation backed by Gioui.
package vggio // import "gonum.org/v1/plot/vg/vggio"

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op/paint"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"
)

// Canvas implements the vg.Canvas interface,
// drawing to an image.Image using vgimg and painting that image
// into a Gioui context.
type Canvas struct {
	*vgimg.Canvas
	gtx layout.Context
}

// DefaultDPI is the default dot resolution for image
// drawing in dots per inch.
const DefaultDPI = 96

// New returns a new image canvas with the provided dimensions and options.
// The currently accepted options are UseDPI and UseBackgroundColor.
// If the resolution or background color are not specified, defaults are used.
func New(gtx layout.Context, w, h vg.Length, opts ...option) *Canvas {
	cfg := &config{
		dpi: DefaultDPI,
		bkg: color.White,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	c := &Canvas{
		gtx: gtx,
		Canvas: vgimg.NewWith(
			vgimg.UseDPI(cfg.dpi),
			vgimg.UseWH(w, h),
			vgimg.UseBackgroundColor(cfg.bkg),
		),
	}
	return c
}

type config struct {
	dpi int
	bkg color.Color
}

type option func(*config)

// UseDPI sets the dots per inch of a canvas. It should only be
// used as an option argument when initializing a new canvas.
func UseDPI(dpi int) option {
	if dpi <= 0 {
		panic("DPI must be > 0.")
	}
	return func(c *config) {
		c.dpi = dpi
	}
}

// UseBackgroundColor specifies the image background color.
// Without UseBackgroundColor, the default color is white.
func UseBackgroundColor(c color.Color) option {
	return func(cfg *config) {
		cfg.bkg = c
	}
}

func (c *Canvas) pt32(p vg.Point) f32.Point {
	_, h := c.Size()
	dpi := c.DPI()
	return f32.Point{
		X: float32(p.X.Dots(dpi)),
		Y: float32(h.Dots(dpi) - p.Y.Dots(dpi)),
	}
}

// Paint paints the canvas' content on the screen.
func (c *Canvas) Paint(e system.FrameEvent) {
	w, h := c.Size()
	box := vg.Rectangle{Max: vg.Point{X: w, Y: h}}
	img := c.Canvas.Image()
	ops := c.gtx.Ops
	min := c.pt32(box.Min)
	max := c.pt32(box.Max)
	r32 := f32.Rect(min.X, min.Y, max.X, max.Y)

	paint.NewImageOp(img).Add(ops)
	paint.PaintOp{Rect: r32}.Add(ops)

	e.Frame(ops)
}
