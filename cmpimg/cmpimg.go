// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cmpimg compares the raw representation of images taking into account
// idiosyncracies related to their underlying format (SVG, PDF, PNG, ...).
package cmpimg // import "gonum.org/v1/plot/cmpimg"

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
	return EqualApprox(typ, raw1, raw2, 0)
}

// EqualApprox takes the raw representation of two images, raw1 and raw2,
// together with the underlying image type ("eps", "jpeg", "jpg", "pdf", "png", "svg", "tiff"),
// a normalized delta parameter to describe how close the matching should be
// performed (delta=0: perfect match, delta=1, loose match)
// and returns whether the two images are equal or not.
//
// EqualApprox may return an error if the decoding of the raw image somehow failed.
// EqualApprox only uses the normalized delta parameter for "jpeg", "jpg", "png",
// and "tiff" images. It ignores that parameter for other document types.
func EqualApprox(typ string, raw1, raw2 []byte, delta float64) (bool, error) {
	switch {
	case delta < 0:
		delta = 0
	case delta > 1:
		delta = 1
	}

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
		pdf1, err := pdf.NewReader(bytes.NewReader(raw1), int64(len(raw1)))
		if err != nil {
			return false, err
		}

		pdf2, err := pdf.NewReader(bytes.NewReader(raw2), int64(len(raw2)))
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
		return cmpImg(v1, v2, delta), nil

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

func cmpImg(v1, v2 image.Image, delta float64) bool {
	img1, ok := v1.(*image.RGBA)
	if !ok {
		img1 = newRGBAFrom(v1)
	}

	img2, ok := v2.(*image.RGBA)
	if !ok {
		img2 = newRGBAFrom(v2)
	}

	if len(img1.Pix) != len(img2.Pix) {
		return false
	}

	max := delta * delta
	bnd := img1.Bounds()
	for x := bnd.Min.X; x < bnd.Max.X; x++ {
		for y := bnd.Min.Y; y < bnd.Max.Y; y++ {
			c1 := img1.RGBAAt(x, y)
			c2 := img2.RGBAAt(x, y)
			if !yiqEqApprox(c1, c2, max) {
				return false
			}
		}
	}

	return true
}

// yiqEqApprox compares the colors of 2 pixels, in the NTSC YIQ color space,
// as described in:
//
//   Measuring perceived color difference using YIQ NTSC
//   transmission color space in mobile applications.
//   Yuriy Kotsarenko, Fernando Ramos.
//
// An electronic version is available at:
//
// - http://www.progmat.uaem.mx:8080/artVol2Num2/Articulo3Vol2Num2.pdf
func yiqEqApprox(c1, c2 color.RGBA, d2 float64) bool {
	const max = 35215.0 // difference between 2 maximally different pixels.

	var (
		r1 = float64(c1.R)
		g1 = float64(c1.G)
		b1 = float64(c1.B)

		r2 = float64(c2.R)
		g2 = float64(c2.G)
		b2 = float64(c2.B)

		y1 = r1*0.29889531 + g1*0.58662247 + b1*0.11448223
		i1 = r1*0.59597799 - g1*0.27417610 - b1*0.32180189
		q1 = r1*0.21147017 - g1*0.52261711 + b1*0.31114694

		y2 = r2*0.29889531 + g2*0.58662247 + b2*0.11448223
		i2 = r2*0.59597799 - g2*0.27417610 - b2*0.32180189
		q2 = r2*0.21147017 - g2*0.52261711 + b2*0.31114694

		y = y1 - y2
		i = i1 - i2
		q = q1 - q2

		diff = 0.5053*y*y + 0.299*i*i + 0.1957*q*q
	)
	return diff <= max*d2
}

func newRGBAFrom(src image.Image) *image.RGBA {
	var (
		bnds = src.Bounds()
		dst  = image.NewRGBA(bnds)
	)
	draw.Draw(dst, bnds, src, image.Point{}, draw.Src)
	return dst
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
