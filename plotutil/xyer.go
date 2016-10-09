// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotutil

import (
	"fmt"

	"github.com/gonum/plot/plotter"
)

// VecXY implements plotter.XYer.
type VecXY struct {
	X, Y []float64
}

// ZipXY is a convenience function to build a VecXY from slices with length checks.
// Panics on length mismatch.
func ZipXY(x, y []float64) plotter.XYer {
	if len(y) != len(x) {
		panic(fmt.Sprintf("VecXY length mismatch %d != %d", len(x), len(y)))
	}
	return &VecXY{X: x, Y: y}
}

// Len returns the length of the VecXY object.
// Panics on slice length mismatch.
// Fulfills the plotter.XYer interface requirement.
func (xy *VecXY) Len() int {
	if len(xy.X) != len(xy.Y) {
		panic(fmt.Sprintf("VecXY length mismatch %d != %d", len(xy.X), len(xy.Y)))
	}
	return len(xy.X)
}

// XY returns the x,y values at index idx.
// fulfills the plotter.XYer interface requirement.
func (xy *VecXY) XY(idx int) (x, y float64) {
	return xy.X[idx], xy.Y[idx]
}

// xyzer implements plotter.XYZer for use by ZipXYZ.
type VecXYZ struct {
	X, Y, Z []float64
}

// ZipXYZ is a convenience function to build a plotter.XYZer from slices.
// Length will be set to the shortest slice length.
// Panics on empty or nil slices.
func ZipXYZ(x, y, z []float64) plotter.XYZer {
	if len(y) != len(x) || len(y) != len(z) {
		panic(fmt.Sprintf("VecXYZ length mismatch %d != %d != %d", len(x), len(y), len(z)))
	}
	return &VecXYZ{X: x, Y: y, Z: z}
}

// Len returns the length of the XYZer object.
// fulfills the plotter.XYZer interface requirement.
func (xyz *VecXYZ) Len() int {
	if len(xyz.Y) != len(xyz.X) || len(xyz.Y) != len(xyz.Z) {
		panic(fmt.Sprintf("VecXYZ length mismatch %d != %d != %d", len(xyz.X), len(xyz.Y), len(xyz.Z)))
	}
	return len(xyz.X)
}

// XYZ returns the x,y,z values at index idx.
// fulfills the plotter.XYZer interface requirement.
func (xyz *VecXYZ) XYZ(idx int) (x, y, z float64) {
	return xyz.X[idx], xyz.Y[idx], xyz.Z[idx]
}
