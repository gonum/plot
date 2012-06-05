package vecimg

import (
	"code.google.com/p/plotinum/vg"
	"testing"
)

func TestFontExtents(t *testing.T) {
	img, err := New(vg.Inches(4), vg.Inches(4))
	if err != nil {
		t.Fatal(err)
	}
	vg.DrawFontExtents(t, img)
	err = img.SavePNG("extents.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFonts(t *testing.T) {
	img, err := New(vg.Inches(4), vg.Inches(4))
	if err != nil {
		t.Fatal(err)
	}
	vg.DrawFonts(t, img)
	err = img.SavePNG("fonts.png")
	if err != nil {
		t.Fatal(err)
	}
}
