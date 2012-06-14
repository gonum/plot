package plt

import (
	"code.google.com/p/plotinum/vg"
	"code.google.com/p/plotinum/vg/veceps"
	"code.google.com/p/plotinum/vg/vecimg"
	"testing"
)

func TestDrawImage(t *testing.T) {
	w, h := vg.Inches(4), vg.Inches(4)
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
	w, h := vg.Inches(4), vg.Inches(4)
	da := NewDrawArea(veceps.New(w, h, "test"), w, h)
	draw(da)
	err := da.Canvas.(*veceps.Canvas).Save("test.eps")
	if err != nil {
		t.Fatal(err)
	}
}

// draw draws a simple test plot
func draw(da *DrawArea) {
	p := New()
	p.Title.Text = "This is a plot with\ntwo different lines"
	p.X.Label.Text = "X Label\ngq"
	p.X.Min = 100000
	p.X.Max = 100001
	p.Y.Label.Text = "Y Label\ngq"
	p.Y.Min = 10
	p.Y.Max = 20
	p.Draw(da)
}
