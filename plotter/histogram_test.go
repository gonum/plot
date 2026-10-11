// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"testing"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func TestHistogram(t *testing.T) {
	cmpimg.CheckPlot(ExampleHistogram, t, "histogram.png")
}

func TestHistogramWideRangeDoesNotPanic(t *testing.T) {
	data := plotter.Values{203.4, -9.223372036854776e17, 220.4}
	h, err := plotter.NewHist(data, 10)
	if err != nil {
		t.Fatalf("NewHist: %v", err)
	}
	var weight float64
	for _, bin := range h.Bins {
		weight += bin.Weight
	}
	if weight != 3 {
		t.Fatalf("total weight %v, want 3", weight)
	}
}

func TestSingletonHistogram(t *testing.T) {
	done := make(chan struct{}, 1)
	go func() {
		defer close(done)
		p := plot.New()
		hist, err := plotter.NewHist(plotter.Values([]float64{1.0}), 60)
		if err != nil {
			t.Errorf("unexpected error from NewHist: %v", err)
			return
		}
		hist.Normalize(1)

		p.Add(hist)

		_, err = p.WriterTo(4*vg.Inch, 4*vg.Inch, "png")
		if err != nil {
			t.Errorf("unexpected error from WriterTo: %v", err)
			return
		}
	}()

	select {
	case <-time.After(10 * time.Second):
		t.Error("histogram timed out")
	case <-done:
	}
}
func TestHistogramLogScale(t *testing.T) {
	cmpimg.CheckPlot(ExampleHistogram_logScaleY, t, "histogram_logy.png")
}
