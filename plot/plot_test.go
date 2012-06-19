package plot

import (
	"code.google.com/p/plotinum/vg"
	"code.google.com/p/plotinum/vg/veceps"
	"code.google.com/p/plotinum/vg/vecimg"
	"math/rand"
	//	"time"
	"testing"
)

var seed = int64(0) // time.Now().UnixNano()

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
	w, h := vg.Inches(4), vg.Inches(2)
	da := NewDrawArea(veceps.New(w, h, "test"), w, h)
	draw(da)
	err := da.Canvas.(*veceps.Canvas).Save("test.eps")
	if err != nil {
		t.Fatal(err)
	}
}

// draw draws a simple test plot
func draw(da *DrawArea) {
	rand.Seed(seed)
	n := 10
	uniform := make(Ys, n)
	normal := make(Ys, n)
	expon := make(Ys, n)
	for i := 0; i < n; i++ {
		uniform[i] = rand.Float64()
		normal[i] = rand.NormFloat64()
		expon[i] = rand.ExpFloat64()
	}
	p, err := New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Plot Title"
	p.X.Label.Text = "Values"

	b0 := MakeHorizBox(vg.Points(20), 0, uniform)
	b1 := MakeHorizBox(vg.Points(20), 1, normal)
	b2 := MakeHorizBox(vg.Points(20), 2, expon)
	p.AddData(b0, b1, b2)
	p.Legend.AddEntry("outliers", b0.GlyphStyle)
	p.NominalY("Uniform\nDistribution", "Normal\nDistribution",
		"Exponential\nDistribution")

	_, med0, _, _ := b0.Statistics()
	_, med1, _, _ := b1.Statistics()
	_, med2, _, _ := b2.Statistics()
	meds := XYs{{med0, b0.X}, {med1, b1.X}, {med2, b2.X}}
	l := Line{meds, DefaultLineStyle}
	s := Scatter{meds, GlyphStyle{Shape: CircleGlyph, Radius: vg.Points(2)}}
	p.AddData(l, s)
	p.Legend.AddEntry("median", l, s)
	p.Draw(da)
}
