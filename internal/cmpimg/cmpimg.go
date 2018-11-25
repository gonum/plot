// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cmpimg compares the raw representation of images taking into account
// idiosyncracies related to their underlying format (SVG, PDF, PNG, ...).
package cmpimg // import "gonum.org/v1/plot/internal/cmpimg"

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"reflect"
	"strings"

	"rsc.io/pdf"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/tiff"
)

// Equal takes the raw representation of two images, raw1 and raw2,
// together with the underlying image type ("eps", "jpeg", "jpg", "pdf", "png", "svg", "tiff"),
// and returns whether the two images are equal or not.
//
// Equal may return an error if the decoding of the raw image somehow failed.
func Equal(typ string, raw1, raw2 []byte) (bool, error) {
	switch typ {
	case "svg", "tex":
		return bytes.Equal(raw1, raw2), nil

	case "eps":
		lines1, lines2 := strings.Split(string(raw1), "\n"), strings.Split(string(raw2), "\n")
		if len(lines1) != len(lines2) {
			return false, nil
		}
		for i, line1 := range lines1 {
			if strings.Contains(line1, "CreationDate") {
				continue
			}
			if line1 != lines2[i] {
				return false, nil
			}
		}
		return true, nil

	case "pdf":
		r1 := bytes.NewReader(raw1)
		pdf1, err := pdf.NewReader(r1, r1.Size())
		if err != nil {
			return false, err
		}

		r2 := bytes.NewReader(raw2)
		pdf2, err := pdf.NewReader(r2, r2.Size())
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
	return cmpPdfValues(pdf1.Trailer(), pdf2.Trailer())
}

func cmpPdfValues(v1, v2 pdf.Value) bool {
	if v1.Kind() != v2.Kind() {
		return false
	}

	switch v1.Kind() {
	case pdf.String:
		return v1.String() == v2.String()
	case pdf.Integer:
		return v1.Int64() == v2.Int64()
	case pdf.Real:
		return v1.Float64() == v2.Float64()
	case pdf.Name:
		return v1.Name() == v2.Name()
	case pdf.Stream:
		r1 := v1.Reader()
		s1, err1 := ioutil.ReadAll(r1)
		r1.Close()
		r2 := v2.Reader()
		s2, err2 := ioutil.ReadAll(r2)
		r2.Close()
		if err1 != nil || err2 != nil || len(s1) != len(s2) {
			return false
		}
		if !bytes.Equal(s1, s2) {
			return false
		}
		fallthrough
	case pdf.Dict:
		keys1, keys2 := v1.Keys(), v2.Keys()
		if len(keys1) != len(keys2) {
			return false
		}
		for i, k := range keys1 {
			if k != keys2[i] {
				return false
			}
			if k == "CreationDate" || k == "Parent" || k == "Font" {
				continue
			}
			if !cmpPdfValues(v1.Key(k), v2.Key(k)) {
				return false
			}
		}
		return true
	case pdf.Array:
		if v1.Len() != v2.Len() {
			return false
		}
		count := v1.Len()
		for i := 0; i < count; i++ {
			if !cmpPdfValues(v1.Index(i), v2.Index(i)) {
				return false
			}
		}
		return true
	}
	return false
}

// Diff calculates an intensity-scaled difference between images a and b
// and places the result in dst, returning the intersection of a, b and
// dst. It is the responsibility of the caller to construct dst so that
// it will overlap with a and b. For the purposes of Diff, alpha is not
// considered.
//
// Diff is not intended to be used for quantitative analysis of the
// difference between the input images, but rather to highlight differences
// between them for testing purposes, so the calculation is rather naive.
func Diff(dst draw.Image, a, b image.Image) image.Rectangle {
	rect := dst.Bounds().Intersect(a.Bounds()).Intersect(b.Bounds())

	// Determine greyscale dynamic range.
	min := uint16(math.MaxUint16)
	max := uint16(0)
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			p := diffColor{a.At(x, y), b.At(x, y)}
			g := color.Gray16Model.Convert(p).(color.Gray16)
			if g.Y < min {
				min = g.Y
			}
			if g.Y > max {
				max = g.Y
			}
		}
	}

	// Render intensity-scaled difference.
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			dst.Set(x, y, scaledColor{
				min: uint32(min), max: uint32(max),
				c: diffColor{a.At(x, y), b.At(x, y)},
			})
		}
	}

	return rect
}

type diffColor struct {
	a, b color.Color
}

func (c diffColor) RGBA() (r, g, b, a uint32) {
	ra, ga, ba, _ := c.a.RGBA()
	rb, gb, bb, _ := c.b.RGBA()
	return diff(ra, rb), diff(ga, gb), diff(ba, bb), math.MaxUint16
}

func diff(a, b uint32) uint32 {
	if a < b {
		return b - a
	}
	return a - b
}

type scaledColor struct {
	min, max uint32
	c        color.Color
}

func (c scaledColor) RGBA() (r, g, b, a uint32) {
	if c.max == c.min {
		return 0, 0, 0, 0
	}
	f := uint32(math.MaxUint16) / (c.max - c.min)
	r, g, b, _ = c.c.RGBA()
	r -= c.min
	r *= f
	g -= c.min
	g *= f
	b -= c.min
	b *= f
	return r, g, b, math.MaxUint16
}
