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
