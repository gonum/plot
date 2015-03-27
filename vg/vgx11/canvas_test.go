// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package vgx11

import (
	"bytes"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// hasX11 checks if X11 is available by looking at the DISPLAY environment variable
func hasX11() bool {
	return os.Getenv("DISPLAY") != ""
}

func TestIssue179(t *testing.T) {
	if !hasX11() {
		t.Skip("no X11 environment")
	}

	scatter, err := plotter.NewScatter(plotter.XYs{{1, 1}, {0, 1}, {0, 0}})
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
	p, err := plot.New()
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
	p.Add(scatter)
	p.HideAxes()

	c, err := New(5.08*vg.Centimeter, 5.08*vg.Centimeter, "test issue179")
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}

	p.Draw(draw.New(c))
	b := bytes.NewBuffer([]byte{})
	err = jpeg.Encode(b, c.ximg, nil)
	if err != nil {
		t.Error(err)
	}

	f, err := os.Open(filepath.Join("..", "vgimg", "testdata", "issue179.jpg"))
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	want, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(b.Bytes(), want) {
		t.Error("Image mismatch")
	}
}
