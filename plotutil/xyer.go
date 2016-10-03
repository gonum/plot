package plotutil

import "github.com/gonum/plot/plotter"

type xyer struct {
	x, y []float64
	len  int
}

// Len returns the length of the XYer object.
// fulfills the plotter.XYer interface requirement.
func (xy *xyer) Len() int {
	return xy.len
}

// XY returns the x,y values at index idx.
// fulfills the plotter.XYer interface requirement.
func (xy *xyer) XY(idx int) (x, y float64) {
	return xy.x[idx], xy.y[idx]
}

// ZipXY is a convenience function to build a plotter.XYer from slices.
// Length will be set to the shortest slice length.
// Panics on empty or nil slices.
func ZipXY(x, y []float64) plotter.XYer {
	switch {
	case len(x) == 0:
		panic("x is an empty or nil slice")
	case len(y) == 0:
		panic("y is an empty or nil slice")
	}
	ln := len(x)
	if len(y) < ln {
		ln = len(y)
	}
	return &xyer{x: x, y: y, len: ln}
}

type xyzer struct {
	x, y, z []float64
	len     int
}

// Len returns the length of the XYZer object.
// fulfills the plotter.XYZer interface requirement.
func (xyz *xyzer) Len() int {
	return xyz.len
}

// XYZ returns the x,y,z values at index idx.
// fulfills the plotter.XYZer interface requirement.
func (xyz *xyzer) XYZ(idx int) (x, y, z float64) {
	return xyz.x[idx], xyz.y[idx], xyz.z[idx]
}

// ZipXYZ is a convenience function to build a plotter.XYZer from slices.
// Length will be set to the shortest slice length.
// Panics on empty or nil slices.
func ZipXYZ(x, y, z []float64) plotter.XYZer {
	switch {
	case len(x) == 0:
		panic("x is an empty or nil slice")
	case len(y) == 0:
		panic("y is an empty or nil slice")
	case len(z) == 0:
		panic("z is an empty or nil slice")
	}
	ln := len(x)
	if len(y) < ln {
		ln = len(y)
	}
	if len(z) < ln {
		ln = len(z)
	}

	return &xyzer{x: x, y: y, z: z, len: ln}
}
