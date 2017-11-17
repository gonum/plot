// package vgall exists for one purpose only: To automatically register
// all image format implementations of vg.CanvasWriterTo.
package vgall // import "gonum.org/v1/plot/vg/vgall"

import (
	_ "gonum.org/v1/plot/vg/vgeps"
	_ "gonum.org/v1/plot/vg/vgimg"
	_ "gonum.org/v1/plot/vg/vgpdf"
	_ "gonum.org/v1/plot/vg/vgsvg"
)
