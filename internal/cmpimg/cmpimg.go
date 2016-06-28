// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cmpimg compares the raw representation of images taking into account
// idiosyncracies related to their underlying format (SVG, PDF, PNG, ...).
package cmpimg

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"reflect"

	_ "golang.org/x/image/tiff"
	"rsc.io/pdf"
)

// Equal takes the raw representation of two images, raw1 and raw2,
// together with the underlying image type ("eps", "jpeg", "jpg", "pdf", "png", "svg", "tiff"),
// and returns whether the two images are equal or not.
//
// Equal may return an error if the decoding of the raw image somehow failed.
func Equal(typ string, raw1, raw2 []byte) (bool, error) {
	switch typ {
	case "svg":
		return bytes.Equal(raw1, raw2), nil

	case "eps":
		return bytes.Equal(raw1, raw2), nil

	case "pdf":
		// TODO(sbinet): bytes.Reader.Size was introduced only after go-1.4
		// use that if/when we drop go-1.4 bwd compat.
		r1 := bytes.NewReader(raw1)
		sz1 := int64(len(raw1))
		pdf1, err := pdf.NewReader(r1, sz1)
		if err != nil {
			return false, err
		}

		// TODO(sbinet): bytes.Reader.Size was introduced only after go-1.4
		// use that if/when we drop go-1.4 bwd compat.
		r2 := bytes.NewReader(raw2)
		sz2 := int64(len(raw2))
		pdf2, err := pdf.NewReader(r2, sz2)
		if err != nil {
			return false, err
		}

		return cmpPdf(pdf1, pdf2), nil

	case "jpeg", "jpg", "png", "tiff":
		v1, _, err := image.Decode(bytes.NewReader(raw1))
		if err != nil {
			return false, err
		}
		v2, _, err := image.Decode(bytes.NewReader(raw2))
		if err != nil {
			return false, err
		}
		return reflect.DeepEqual(v1, v2), nil

	default:
		return false, fmt.Errorf("cmpimg: unknown image type %q", typ)
	}
}

func cmpPdf(pdf1, pdf2 *pdf.Reader) bool {
	n1 := pdf1.NumPage()
	n2 := pdf2.NumPage()
	if n1 != n2 {
		return false
	}

	for i := 1; i <= n1; i++ {
		p1 := pdf1.Page(i).Content()
		p2 := pdf2.Page(i).Content()
		if !reflect.DeepEqual(p1, p2) {
			return false
		}
	}

	t1 := pdf1.Trailer().String()
	t2 := pdf2.Trailer().String()
	return t1 == t2
}
