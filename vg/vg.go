// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

// vg defines an interface for drawing 2D vector graphics.
// This package is designed with the hope that many different
// vector graphics back-ends can conform to the interface.
package vg

import (
	"image/color"
)

// A Canvas is the main drawing interface for 2D vector
// graphics.  The origin is in the bottom left corner.
type Canvas interface {
	// SetLineWidth sets the width of stroked paths.
	// If the width is set to 0 then stroked lines are
	// not drawn.
	//
	// The initial line width is 1 point.
	SetLineWidth(Length)

	// SetLineDash sets the dash pattern for lines.
	// The first argument is the pattern (units on,
	// units off, etc.) and the second argument
	// is the initial offset into the pattern.
	//
	// The inital dash pattern is a solid line.
	SetLineDash([]Length, Length)

	// SetColor sets the current drawing color.
	// Note that fill color and stroke color are
	// the same so if you want different fill
	// and stroke colorls then you must use two
	// separate calls to SetColor.
	//
	// The initial color is black.  If SetColor is
	// called with a nil color then black is used.
	SetColor(color.Color)

	// Rotate applies a rotation transform to the
	// context.  The parameter is specified in
	// radians.
	Rotate(float64)

	// Translate applies a translational transform
	// to the context.
	Translate(Length, Length)

	// Scale applies a scaling transform to the
	// context.
	Scale(float64, float64)

	// Push saves the current line width, the
	// current dash pattern, the current
	// transforms, and the current color
	// onto a stack so that the state can later
	// be restored by calling Pop().
	Push()

	// Pop restores the context saved by the
	// corresponding call to Push().
	Pop()

	// Stroke strokes the given path.
	Stroke(Path)

	// Fill fills the given path.
	Fill(Path)

	// FillString fills in text at the specified
	// location using the given font.
	FillString(Font, Length, Length, string)

	// DPI returns the number of canvas dots in
	// an inch.
	DPI() float64
}

// Initialize sets all of the canvas's values to their
// initial values.
func Initialize(c Canvas) {
	c.SetLineWidth(Points(1))
	c.SetLineDash([]Length{}, 0)
	c.SetColor(color.Black)
}

type Path []PathComp

// Move moves the current location of the path to
// the given point.
func (p *Path) Move(x, y Length) {
	*p = append(*p, PathComp{Type: MoveComp, X: x, Y: y})
}

// Line draws a line from the current point to the
// given point.
func (p *Path) Line(x, y Length) {
	*p = append(*p, PathComp{Type: LineComp, X: x, Y: y})
}

// Arc adds an arc to the path defined by the center
// point of the arc's circle, the radius of the circle
// and the start and sweep angles.
func (p *Path) Arc(x, y, rad Length, s, a float64) {
	*p = append(*p, PathComp{
		Type:   ArcComp,
		X:      x,
		Y:      y,
		Radius: rad,
		Start:  s,
		Angle:  a,
	})
}

// Close closes the path by connecting the current
// location to the start location with a line.
func (p *Path) Close() {
	*p = append(*p, PathComp{Type: CloseComp})
}

// Constants that tag the type of each path
// component.
const (
	MoveComp = iota
	LineComp
	ArcComp
	CloseComp
)

// A PathComp is a component of a path structure.
type PathComp struct {
	// Type is the type of a particluar component.
	// Based on the type, each of the following
	// fields may have a different meaning or may
	// be meaningless.
	Type int

	// The X and Y fields are used as the destination
	// of a MoveComp or LineComp and are the center
	// point of an ArcComp.  They are not used in
	// the CloseComp.
	X, Y Length

	// Radius is only used for ArcComps, it is
	// the radius of the circle defining the arc.
	Radius Length

	// Start and Angle are only used for ArcComps.
	// They define the start angle and sweep angle of
	// the arc around the circle.  The units of the
	// angle are radians.
	Start, Angle float64
}
