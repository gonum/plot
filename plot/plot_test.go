package plot

import (
	"code.google.com/p/plotinum/vg"
	"code.google.com/p/plotinum/vg/veceps"
	"code.google.com/p/plotinum/vg/vecimg"
	"testing"
)

func TestDrawImage(t *testing.T) {
	w, h := vg.Inches(3), vg.Inches(2)
	img, err := vecimg.New(w, h)
	if err != nil {
		t.Error(err)
	}
	da := NewDrawArea(img, w, h)
	draw(da)
	err = img.SavePNG("test.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDrawEps(t *testing.T) {
	w, h := vg.Inches(3), vg.Inches(2)
	da := NewDrawArea(veceps.New(w, h, "test"), w, h)
	draw(da)
	err := da.Canvas.(*veceps.Canvas).Save("test.eps")
	if err != nil {
		t.Fatal(err)
	}
}

// draw draws a simple test plot
func draw(da *DrawArea) {
	Example_boxPlot(da)
}
