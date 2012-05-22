package vecimg

import (
	"testing"
	"code.google.com/p/plotinum/vecgfx"
)

func TestFontExtents(t *testing.T) {
	img, err := New(4, 4)
	if err != nil {
		t.Fatal(err)
	}
	vecgfx.DrawFontExtents(t, img)
	err = img.SavePNG("extents.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFonts(t *testing.T) {
	img, err := New(4, 4)
	if err != nil {
		t.Fatal(err)
	}
	vecgfx.DrawFonts(t, img)
	err = img.SavePNG("fonts.png")
	if err != nil {
		t.Fatal(err)
	}
}