// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vggio_test

import (
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vggio"
)

func ExampleCanvas() {
	const (
		w   = 20 * vg.Centimeter
		h   = 15 * vg.Centimeter
		dpi = 96
	)
	go func(w, h vg.Length) {
		win := app.NewWindow(
			app.Title("Gonum"),
			app.Size(
				unit.Px(float32(w.Dots(dpi))),
				unit.Px(float32(h.Dots(dpi))),
			),
		)
		defer win.Close()

		done := time.NewTimer(2 * time.Second)
		defer done.Stop()

		for {
			select {
			case e := <-win.Events():
				switch e := e.(type) {
				case system.FrameEvent:
					p, err := plot.New()
					if err != nil {
						log.Fatalf("could not create plot: %+v", err)
					}
					p.Title.Text = "My title"
					p.X.Label.Text = "X"
					p.Y.Label.Text = "Y"

					cnv := vggio.New(e, w, h, vggio.UseDPI(dpi))
					p.Draw(draw.New(cnv))
					cnv.Paint(e)

				case key.Event:
					switch e.Name {
					case "Q", key.NameEscape:
						os.Exit(0)
					}
				}
			case <-done.C:
				os.Exit(0)
			}
		}
	}(w, h)

	app.Main()

	// Output:
}
