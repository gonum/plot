// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot_test

import (
	"image/color"
	"math"
	"reflect"
	"testing"

	xfont "golang.org/x/image/font"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/recorder"
)

func TestLegend_standalone(t *testing.T) {
	cmpimg.CheckPlot(ExampleLegend_standalone, t, "legend_standalone.png")
}

func TestLegend_background(t *testing.T) {
	cmpimg.CheckPlot(ExampleLegend_background, t, "legend_background.png")
}

func TestLegendBackground(t *testing.T) {
	for _, left := range []bool{false, true} {
		for _, top := range []bool{false, true} {
			var canvas recorder.Canvas
			c := draw.Canvas{Canvas: &canvas, Rectangle: vg.Rectangle{Min: vg.Point{X: 10, Y: 20}, Max: vg.Point{X: 210, Y: 120}}}
			l := plot.NewLegend()
			l.Add("first")
			l.Add("second")
			l.Left, l.Top = left, top
			l.XOffs, l.YOffs = 3, 4
			l.BackgroundColor = color.NRGBA{R: 230, G: 240, B: 250, A: 255}
			l.Draw(c)
			if len(canvas.Actions) < 2 {
				t.Fatalf("no background drawn: %v", canvas.Actions)
			}
			clr, ok := canvas.Actions[0].(*recorder.SetColor)
			if !ok || clr.Color != l.BackgroundColor {
				t.Fatalf("first action = %v, want background color", canvas.Actions[0])
			}
			fill, ok := canvas.Actions[1].(*recorder.Fill)
			if !ok {
				t.Fatalf("second action = %v, want background fill", canvas.Actions[1])
			}
			r := l.Rectangle(c)
			if !reflect.DeepEqual(fill.Path, r.Path()) {
				t.Errorf("left=%t top=%t: background %v does not match legend rectangle %v", left, top, fill.Path, r)
			}
			if left && r.Min.X != c.Min.X+l.XOffs || !left && r.Max.X != c.Max.X+l.XOffs {
				t.Errorf("left=%t: legend rectangle has wrong horizontal anchor: %v", left, r)
			}
		}
	}
}

type legendBoundsThumbnail struct{ bounds *[]vg.Rectangle }

func (t legendBoundsThumbnail) Thumbnail(c *draw.Canvas) {
	*t.bounds = append(*t.bounds, c.Rectangle)
}

func TestLegendBackgroundBounds(t *testing.T) {
	for _, left := range []bool{false, true} {
		for _, top := range []bool{false, true} {
			for _, yposition := range []float64{draw.PosBottom, draw.PosCenter, draw.PosTop} {
				for _, yalign := range []text.YAlignment{text.YBottom, text.YCenter, text.YTop} {
					for _, rotation := range []float64{0, math.Pi / 6} {
						var canvas recorder.Canvas
						c := draw.Canvas{Canvas: &canvas, Rectangle: vg.Rectangle{Min: vg.Point{X: 10, Y: 20}, Max: vg.Point{X: 210, Y: 120}}}
						l := plot.NewLegend()
						l.Left, l.Top = left, top
						l.XOffs, l.YOffs = 3, -4
						l.YPosition = yposition
						l.TextStyle.YAlign = yalign
						l.TextStyle.Rotation = rotation
						l.TextStyle.Font.Style = xfont.StyleItalic
						l.Padding = vg.Points(2)
						l.BackgroundColor = color.White
						var icons []vg.Rectangle
						thumb := legendBoundsThumbnail{bounds: &icons}
						l.Add("Ågyp", thumb)
						l.Add("first\n\nsecond", thumb)
						l.Draw(c)
						r := l.Rectangle(c)
						check := func(p vg.Point) {
							t.Helper()
							const tol vg.Length = 1.0 / 64
							if p.X < r.Min.X-tol || p.X > r.Max.X+tol || p.Y < r.Min.Y-tol || p.Y > r.Max.Y+tol {
								t.Fatalf("left=%t top=%t yposition=%g yalign=%g rotation=%g: point %v outside background %v", left, top, yposition, yalign, rotation, p, r)
							}
						}
						for _, icon := range icons {
							check(icon.Min)
							check(icon.Max)
						}
						var angle float64
						var stack []float64
						for _, action := range canvas.Actions {
							switch a := action.(type) {
							case *recorder.Push:
								stack = append(stack, angle)
							case *recorder.Pop:
								angle, stack = stack[len(stack)-1], stack[:len(stack)-1]
							case *recorder.Rotate:
								angle += a.Angle
							case *recorder.FillString:
								if a.String == "" {
									continue
								}
								fnt := l.TextStyle.Handler.Cache().Lookup(a.Font, a.Size)
								face := fnt.FontFace(72)
								bounds, _ := xfont.BoundString(face, a.String)
								face.Close()
								for _, x := range []int32{int32(bounds.Min.X), int32(bounds.Max.X)} {
									for _, y := range []int32{int32(bounds.Min.Y), int32(bounds.Max.Y)} {
										p := a.Point.Add(vg.Point{X: vg.Points(float64(x) / 64), Y: -vg.Points(float64(y) / 64)})
										sin, cos := math.Sincos(angle)
										check(vg.Point{X: p.X*vg.Length(cos) - p.Y*vg.Length(sin), Y: p.X*vg.Length(sin) + p.Y*vg.Length(cos)})
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func TestLegendNoBackground(t *testing.T) {
	for _, empty := range []bool{false, true} {
		var canvas recorder.Canvas
		c := draw.Canvas{Canvas: &canvas, Rectangle: vg.Rectangle{Max: vg.Point{X: 100, Y: 100}}}
		l := plot.NewLegend()
		if empty {
			l.BackgroundColor = color.White
		} else {
			l.Add("entry")
		}
		l.Draw(c)
		for _, action := range canvas.Actions {
			if _, ok := action.(*recorder.Fill); ok {
				t.Errorf("unexpected background fill for empty=%t", empty)
			}
		}
	}
}

func TestLegendBackgroundZeroFont(t *testing.T) {
	var canvas recorder.Canvas
	c := draw.Canvas{Canvas: &canvas, Rectangle: vg.Rectangle{Max: vg.Point{X: 100, Y: 100}}}
	l := plot.NewLegend()
	l.TextStyle.Font.Size = 0
	l.BackgroundColor = color.White
	l.Add("entry")
	l.Draw(c)
	r := l.Rectangle(c)
	for _, v := range []vg.Length{r.Min.X, r.Min.Y, r.Max.X, r.Max.Y} {
		if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
			t.Fatalf("non-finite zero-font legend bounds: %v", r)
		}
	}
}

func TestLegendRectangleEmpty(t *testing.T) {
	c := draw.Canvas{Rectangle: vg.Rectangle{Min: vg.Point{X: 10, Y: 20}, Max: vg.Point{X: 210, Y: 120}}}
	for _, left := range []bool{false, true} {
		for _, top := range []bool{false, true} {
			l := plot.NewLegend()
			l.Left, l.Top = left, top
			l.XOffs, l.YOffs = 3, -4
			want := vg.Point{X: 213, Y: 16}
			if left {
				want.X = 13
			}
			if top {
				want.Y = 116
			}
			if got := l.Rectangle(c); got.Min != want || got.Max != want {
				t.Errorf("left=%t top=%t: empty legend bounds=%v, want point %v", left, top, got, want)
			}
		}
	}
}
