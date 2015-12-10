// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgimg

import (
	"bufio"
	"io"

	"golang.org/x/image/tiff"
)

// A TiffCanvas is an image canvas with a WriteTo method that
// writes a tiff image.
type TiffCanvas struct {
	*Canvas
}

// WriteTo implements the io.WriterTo interface, writing a tiff image.
func (c TiffCanvas) WriteTo(w io.Writer) (int64, error) {
	wc := writerCounter{Writer: w}
	b := bufio.NewWriter(&wc)
	if err := tiff.Encode(b, c.img, nil); err != nil {
		return wc.n, err
	}
	err := b.Flush()
	return wc.n, err
}
