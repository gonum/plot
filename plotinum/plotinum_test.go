package plotinum

import (
	"testing"
	"plotinum/vecgfx/vecimg"
)

func TestDraw(t *testing.T) {
	img, err := vecimg.New(4, 4)
	if err != nil {
		t.Fatal(err)
	}

	da := DrawArea{
		Canvas: img,
		Rect: Rect{ Min: Point{ 0, 0 },
			Max: Point{ 4*img.DPI(), 4*img.DPI() },
		},
	}
	da.Stroke(RectPath(da.Rect))
	da.Min.X += 10
	da.Max.X -= 10
	ax := MakeAxis(1, 10)
	ax.Label = "X-Axis gq"
	ax.DrawHoriz(&da, img)

	err = img.SavePNG("plot.png")
	if err != nil {
		t.Fatal(err)
	}
}