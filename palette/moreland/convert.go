// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moreland

import (
	"image/color"
	"math"
)

// rgb represents a physically linear RGB color.
type rgb struct {
	R, G, B float64
}

// cieXYZ returns a CIE XYZ color representation of the receiver.
func (c rgb) cieXYZ() cieXYZ {
	return cieXYZ{
		X: 0.4124*c.R + 0.3576*c.G + 0.1805*c.B,
		Y: 0.2126*c.R + 0.7152*c.G + 0.0722*c.B,
		Z: 0.0193*c.R + 0.1192*c.G + 0.9505*c.B,
	}
}

// sRGBA returns an sRGB color representation of the receiver using the
// provided alpha which must be in [0, 1].
func (c rgb) sRGBA(alpha float64) sRGBA {
	// f converts from a linear RGB component to an sRGB component.
	f := func(v float64) float64 {
		if v > 0.0031308 {
			return 1.055*math.Pow(v, 1/2.4) - 0.055
		}
		return 12.92 * v
	}

	return sRGBA{
		R: f(c.R),
		G: f(c.G),
		B: f(c.B),
		A: alpha,
	}
}

// cieXYZ represents a color in CIE XYZ space.
// Y must be in the range [0,1]. X and Z must be greater than 0.
type cieXYZ struct {
	X, Y, Z float64
}

// rgb returns a linear RGB representation of the receiver.
func (c cieXYZ) rgb() rgb {
	return rgb{
		R: c.X*3.2406 + c.Y*-1.5372 + c.Z*-0.4986,
		G: c.X*-0.9689 + c.Y*1.8758 + c.Z*0.0415,
		B: c.X*0.0557 + c.Y*-0.204 + c.Z*1.057,
	}
}

// cieLAB returns a CIELAB color representation of the receiver.
func (c cieXYZ) cieLAB() cieLAB {
	// f is an intermediate step in converting from CIE XYZ to CIE LAB.
	f := func(v float64) float64 {
		if v > 0.008856 {
			return math.Pow(v, 1.0/3.0)
		}
		return 7.787*v + 16.0/116.0
	}

	tempX := f(c.X / 0.9505)
	tempY := f(c.Y)
	tempZ := f(c.Z / 1.089)
	return cieLAB{
		L: (116.0 * tempY) - 16.0,
		A: 500.0 * (tempX - tempY),
		B: 200 * (tempY - tempZ),
	}
}

// sRGBA represents a color within the sRGB color space, with an alpha channel
// but not premultiplied. All values must be in the range [0,1].
type sRGBA struct {
	R, G, B, A float64
}

// rgb returns a linear RGB representation of the receiver.
func (c sRGBA) rgb() rgb {
	// f converts from an sRGB component to a linear RGB component.
	f := func(v float64) float64 {
		if v > 0.04045 {
			return math.Pow((v+0.055)/1.055, 2.4)
		}
		return v / 12.92
	}

	return rgb{
		R: f(c.R),
		G: f(c.G),
		B: f(c.B),
	}
}

// RGBA implements the color.Color interface.
func (c sRGBA) RGBA() (r, g, b, a uint32) {
	return uint32(c.R * c.A * 0xffff), uint32(c.G * c.A * 0xffff), uint32(c.B * c.A * 0xffff), uint32(c.A * 0xffff)
}

// cieLAB returns a CIE LAB representation of the receiver.
func (c sRGBA) cieLAB() cieLAB {
	return c.rgb().cieXYZ().cieLAB()
}

// colorTosRGBA converts a color to an sRGBA.
func colorTosRGBA(c color.Color) sRGBA {
	r, g, b, a := c.RGBA()
	if a == 0 {
		return sRGBA{}
	}
	return sRGBA{
		R: float64(r) / float64(a),
		G: float64(g) / float64(a),
		B: float64(b) / float64(a),
		A: float64(a) / 0xffff,
	}
}

// clamp forces all channels in c to be within the range [0, 1].
func (c *sRGBA) clamp() {
	if c.R > 1 {
		c.R = 1
	}
	if c.G > 1 {
		c.G = 1
	}
	if c.B > 1 {
		c.B = 1
	}
	if c.A > 1 {
		c.A = 1
	}
	if c.R < 0 {
		c.R = 0
	}
	if c.G < 0 {
		c.G = 0
	}
	if c.B < 0 {
		c.B = 0
	}
	if c.A < 0 {
		c.A = 0
	}
}

// cieLAB represents a color in CIE LAB space.
// L must be in the range [0, 100].
type cieLAB struct {
	L, A, B float64
}

// sRGBA return a linear RGB color representation of the receiver using the
// provided alpha which must be in [0, 1].
func (c cieLAB) sRGBA(alpha float64) sRGBA {
	return c.cieXYZ().rgb().sRGBA(alpha)
}

// cieXYZ returns a CIE XYZ color representation of the receiver.
func (c cieLAB) cieXYZ() cieXYZ {
	// f is an intermediate step in converting from CIE LAB to CIE XYZ.
	f := func(v float64) float64 {
		const (
			xlim = 0.008856
			a    = 7.787
			b    = 16. / 116.
			ylim = a*xlim + b
		)
		if v > ylim {
			return v * v * v
		}
		return (v - b) / a
	}

	// Reference white-point D65
	const xn, yn, zn = 0.95047, 1.0, 1.08883
	return cieXYZ{
		X: xn * f((c.A/500)+(c.L+16)/116),
		Y: yn * f((c.L+16)/116),
		Z: zn * f((c.L+16)/116-(c.B/200)),
	}
}

// MSH returns an MSH color representation of the receiver.
func (c cieLAB) MSH() msh {
	m := math.Pow(c.L*c.L+c.A*c.A+c.B*c.B, 0.5)
	return msh{
		M: m,
		S: math.Acos(c.L / m),
		H: math.Atan2(c.B, c.A),
	}
}

// MSH represents a color in Magnitude-Saturation-Hue color space.
type msh struct {
	M, S, H float64
}

// colorToMSH converts a color to MSH space.
// TODO: If msh ever becomes exported, change this to implment color.Model
func colorToMSH(c color.Color) msh {
	return colorTosRGBA(c).cieLAB().MSH()
}

// cieLAB returns a CIELAB representation of the receiver.
func (c msh) cieLAB() cieLAB {
	return cieLAB{
		L: c.M * math.Cos(c.S),
		A: c.M * math.Sin(c.S) * math.Cos(c.H),
		B: c.M * math.Sin(c.S) * math.Sin(c.H),
	}
}

// RGBA implements the color.Color interface.
func (c msh) RGBA() (r, g, b, a uint32) {
	return c.cieLAB().sRGBA(1.0).RGBA()
}

// hueTwist returns the hue twist between color c and converge magnitude
// convergeM.
func hueTwist(c msh, convergeM float64) float64 {
	signH := c.H / math.Abs(c.H)
	return signH * c.S * math.Sqrt(convergeM*convergeM-c.M*c.M) / (c.M * math.Sin(c.S))
}
