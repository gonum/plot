// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vg defines an interface for drawing 2D vector graphics.
// This package is designed with the hope that many different
// vector graphics back-ends can conform to the interface.
package vg

import (
	"fmt"
	"image/color"
	"io"
	"sync"
)

// A Canvas is the main drawing interface for 2D vector
// graphics.  The origin is in the bottom left corner.
type Canvas interface {
	// SetLineWidth sets the width of stroked paths.
	// If the width is not positive then stroked lines
	// are not drawn.
	//
	// The initial line width is 1 point.
	SetLineWidth(Length)

	// SetLineDash sets the dash pattern for lines.
	// The pattern slice specifies the lengths of
	// alternating dashes and gaps, and the offset
	// specifies the distance into the dash pattern
	// to start the dash.
	//
	// The initial dash pattern is a solid line.
	SetLineDash(pattern []Length, offset Length)

	// SetColor sets the current drawing color.
	// Note that fill color and stroke color are
	// the same so if you want different fill
	// and stroke colors then you must use two
	// separate calls to SetColor.
	//
	// The initial color is black.  If SetColor is
	// called with a nil color then black is used.
	SetColor(color.Color)

	// Rotate applies a rotation transform to the
	// context.  The parameter is specified in
	// radians.
	Rotate(rad float64)

	// Translate applies a translational transform
	// to the context.
	Translate(x, y Length)

	// Scale applies a scaling transform to the
	// context.
	Scale(x, y float64)

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
	FillString(f Font, x, y Length, text string)

	// DPI returns the number of canvas dots in
	// an inch.
	DPI() float64
}

// CanvasSizer is a Canvas with a defined size.
type CanvasSizer interface {
	Canvas
	Size() (x, y Length)
}

// CanvasWriterTo wraps behavior required for writing
// a Canvas representation to a concrete image.
type CanvasWriterTo interface {
	CanvasSizer
	io.WriterTo
}

var (
	formatLock sync.RWMutex
	formats    = make(map[string]func(w, h Length) CanvasWriterTo)
)

// Register allows a format to be registered with the vg package.
// Registered formats can be returned by NewCanvasWriterTo.
// Register will panic if it is called twice with the same format
// or if new is nil
func Register(format string, new func(w, h Length) CanvasWriterTo) {
	if new == nil {
		panic("vg: Register new format function is nil")
	}
	formatLock.Lock()
	defer formatLock.Unlock()
	if _, ok := formats[format]; ok {
		panic("vg: Register called twice for format " + format)
	}
	formats[format] = new
}

// Formats returns a slice of registered format names.
func Formats() []string {
	formatLock.RLock()
	f := make([]string, 0, len(formats))
	for k := range formats {
		f = append(f, k)
	}
	formatLock.RUnlock()
	return f
}

// NewCanvasWriterTo returns a concrete canvas type of the specified format
// and size.
func NewCanvasWriterTo(format string, width, height Length) (CanvasWriterTo, error) {
	formatLock.RLock()
	defer formatLock.RUnlock()
	new, ok := formats[format]
	if !ok {
		return nil, fmt.Errorf("unsupported format: %q", format)
	}
	return new(width, height), nil
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
