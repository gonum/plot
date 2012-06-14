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
	p.AddData(MakeLine(DefaultLineStyle, DataPoints{ {100000, 10}, {100000.5, 30}, {100001, 10}} ))
	p.AddData(MakeScatter(DefaultGlyphStyle, DataPoints{ {100000, 10}, {100000.5, 30}, {100001, 10}} ))
	gsty := DefaultGlyphStyle
	gsty.Shape = RingGlyph
	gsty.Radius = vg.Points(18)
	p.Title.Text = "This is a plot with\ntwo different lines"
	p.X.Label.Text = "X Label\ngq"
	p.Y.Min = 10
	p.Y.Max = 20
	p.Y.Tick.Marker = ConstantTicks([]Tick{
		{10, "Ten\n(10)"}, {15, ""},
		{20, "Twenty\n(20)"}})
	p.Y.Label.Text = "Y Label\ngq"
	p.Draw(da)
}
