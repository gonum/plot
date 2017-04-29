// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw

import (
	"fmt"

	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/vgeps"
	"github.com/gonum/plot/vg/vgimg"
	"github.com/gonum/plot/vg/vgpdf"
	"github.com/gonum/plot/vg/vgsvg"
)

// NewFormattedCanvas creates a new vg.CanvasWriterTo with the specified
// image format.
//
// Supported formats are:
//
//  eps, jpg|jpeg, pdf, png, svg, and tif|tiff.
func NewFormattedCanvas(w, h vg.Length, format string) (vg.CanvasWriterTo, error) {
	var c vg.CanvasWriterTo
	switch format {
	case "eps":
		c = vgeps.New(w, h)

	case "jpg", "jpeg":
		c = vgimg.JpegCanvas{Canvas: vgimg.New(w, h)}

	case "pdf":
		c = vgpdf.New(w, h)

	case "png":
		c = vgimg.PngCanvas{Canvas: vgimg.New(w, h)}

	case "svg":
		c = vgsvg.New(w, h)

	case "tif", "tiff":
		c = vgimg.TiffCanvas{Canvas: vgimg.New(w, h)}

	default:
		return nil, fmt.Errorf("unsupported format: %q", format)
	}
	return c, nil
}
