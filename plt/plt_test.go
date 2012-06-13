package plt

import (
	"code.google.com/p/plotinum/vg"
	"code.google.com/p/plotinum/vg/veceps"
	"code.google.com/p/plotinum/vg/vecimg"
	"testing"
)

func TestDrawImage(t *testing.T) {
	da, err := NewPNGDrawArea(vg.Inches(4), vg.Inches(4))
	if err != nil {
		t.Error(err)
	}
	draw(da)
	err = da.Canvas.(*vecimg.Canvas).SavePNG("test.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDrawEps(t *testing.T) {
	da := NewEPSDrawArea(vg.Inches(4), vg.Inches(4), "test")
	draw(da)
	err := da.Canvas.(*veceps.Canvas).Save("test.eps")
	if err != nil {
		t.Fatal(err)
	}
}

// draw draws a simple test plot
func draw(da *drawArea) {
	p := New()
	p.AddData(MakeLine(DefaultLineStyle,
		Point{100000, 10},
		Point{100000.5, 30},
		Point{100001, 10}))
	p.AddData(MakeScatter(DefaultGlyphStyle,
		Point{100000, 10},
		Point{100000.5, 30},
		Point{100001, 10}))
	p.Title.Text = "This is a plot with\ntwo different lines"
	p.X.Label.Text = "X Label\ngq"
	p.Y.Min = 10
	p.Y.Max = 20
	p.Y.Tick.Marker = ConstantTicks([]Tick{{10, "Ten\n(10)"}, {15, ""}, {20, "Twenty\n(20)"}})
	p.Y.Label.Text = "Y Label\ngq"
	p.Draw(da)
}
