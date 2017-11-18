// +build !minimal

package plot // import "gonum.org/v1/plot"

import (
	_ "gonum.org/v1/plot/vg/vgeps"
	_ "gonum.org/v1/plot/vg/vgimg"
	_ "gonum.org/v1/plot/vg/vgpdf"
	_ "gonum.org/v1/plot/vg/vgsvg"
)
