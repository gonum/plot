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

	da := &DrawArea{
		Canvas: img,
		Rect: Rect{ Min: Point{ 0, 0 },
			Sz: Point{ 4*img.DPI(), 4*img.DPI() },
		},
	}
	da.Stroke(RectPath(da.Rect))

	plot := NewPlot()
	plot.Title = "This is a plot"
	plot.XAxis = MakeAxis(1, 10)
	plot.XAxis.Label = "X-Axis gq"
	plot.YAxis = MakeAxis(1, 10)
	plot.YAxis.Label = "Y-Axis gq"
	plot.Draw(da)

	err = img.SavePNG("plot.png")
	if err != nil {
		t.Fatal(err)
	}
}