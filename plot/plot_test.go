package plot

import (
	"code.google.com/p/plotinum/plt"
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
	da := plt.NewDrawArea(img, w, h)
	draw(da)
	err = img.SavePNG("test.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDrawEps(t *testing.T) {
	w, h := vg.Inches(4), vg.Inches(4)
	da := plt.NewDrawArea(veceps.New(w, h, "test"), w, h)
	draw(da)
	err := da.Canvas.(*veceps.Canvas).Save("test.eps")
	if err != nil {
		t.Fatal(err)
	}
}

// draw draws a simple test plot
func draw(da *plt.DrawArea) {
	p := plt.New()
	p.Title.Text = "Title"
	p.X.Label.Text = "X Label"
	p.Y.Label.Text = "Y Label"
	p.AddData(MakeBox(vg.Points(12), 0, Values{
		-50, -45, 5, 10, 15, 20, 25, 30, 35, 40, 80,
	}))
	p.AddData(MakeBox(vg.Points(12), 1, Values{
		-50, -45, 5, 10, 15, 20, 25, 30, 35, 40, 80,
	}))
	p.Draw(da)
}
