package plt

import (
	"code.google.com/p/plotinum/vecgfx/veceps"
	"code.google.com/p/plotinum/vecgfx/vecimg"
	"testing"
)

func TestDrawImage(t *testing.T) {
	img, err := vecimg.New(4, 4)
	if err != nil {
		t.Fatal(err)
	}

	da := &DrawArea{
		Canvas: img,
		Rect: Rect{Min: Point{0, 0},
			Size: Point{4 * img.DPI(), 4 * img.DPI()},
		},
	}
	draw(da)
	err = img.SavePNG("test.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDrawEps(t *testing.T) {
	eps := veceps.New(4, 4, "test")
	da := &DrawArea{
		Canvas: eps,
		Rect: Rect{Min: Point{0, 0},
			Size: Point{4 * eps.DPI(), 4 * eps.DPI()},
		},
	}
	draw(da)
	err := eps.Save("test.eps")
	if err != nil {
		t.Fatal(err)
	}
}

// draw draws a simple test plot
func draw(da *DrawArea) {
	plot := NewPlot()
	plot.Title = "This is a plot"
	plot.XAxis.Min = 1
	plot.XAxis.Max = 10
	plot.XAxis.Label = "X-Axis gq"
	plot.YAxis.Min = 10
	plot.YAxis.Max = 20
	plot.YAxis.Ticks.TickMarker = ConstantTicks([]Tick{{10, "ten"}, {15, ""}, {20, "twenty"}})
	plot.YAxis.Label = "Y-Axis gq"
	plot.Draw(da)
}
