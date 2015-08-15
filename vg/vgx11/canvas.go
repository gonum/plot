// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

// Package vgx11 implements X-Window vg support.
package vgx11

import (
	"image"
	"image/draw"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/llgcode/draw2d/draw2dimg"

	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/vgimg"
)

// dpi is the number of dots per inch.
const dpi = 96

// Canvas implements the vg.Canvas interface,
// drawing to an image.Image using draw2d.
type Canvas struct {
	*vgimg.Canvas

	// X window values
	x    *xgbutil.XUtil
	ximg *xgraphics.Image
	wid  *xwindow.Window
}

// New returns a new image canvas with
// the size specified  rounded up to the
// nearest pixel.
func New(width, height vg.Length, name string) (*Canvas, error) {
	w := width / vg.Inch * dpi
	h := height / vg.Inch * dpi
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

	return NewImage(img, name)
}

// NewImage returns a new image canvas
// that draws to the given image.  The
// minimum point of the given image
// should probably be 0,0.
func NewImage(img draw.Image, name string) (*Canvas, error) {
	w := float64(img.Bounds().Max.X - img.Bounds().Min.X)
	h := float64(img.Bounds().Max.Y - img.Bounds().Min.Y)

	X, err := xgbutil.NewConn()
	if err != nil {
		return nil, err
	}
	keybind.Initialize(X)
	ximg := xgraphics.New(X, image.Rect(0, 0, int(w), int(h)))
	err = ximg.CreatePixmap()
	if err != nil {
		return nil, err
	}
	painter := NewPainter(ximg)
	gc := draw2dimg.NewGraphicContextWithPainter(ximg, painter)
	gc.SetDPI(dpi)
	gc.Scale(1, -1)
	gc.Translate(0, -h)

	wid := ximg.XShowExtra(name, true)
	go func() {
		xevent.Main(X)
	}()

	c := &Canvas{
		Canvas: vgimg.NewWith(vgimg.UseImageWithContext(img, gc)),
		x:      X,
		ximg:   ximg,
		wid:    wid,
	}
	vg.Initialize(c)
	return c, nil
}

func (c *Canvas) Paint() {
	c.ximg.XDraw()
	c.ximg.XPaint(c.wid.Id)
}
