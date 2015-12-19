// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
)

var generateTestData = flag.Bool("regen", false, "Uses the current state to regenerate the test data.")

// checkPlot checks a generated plot against a previously created reference.
// If generateTestData = true, it regereates the reference.
func checkPlot(ExampleFunc func(), t *testing.T, paths ...string) {
	if *generateTestData {
		// Recreate Golden images and exit.
		ExampleFunc()
		return
	}

	want := make([][]byte, len(paths))
	// Read Golden Images before overwriting them.
	for i, path := range paths {
		f, err := os.Open(filepath.Join("testdata", path))
		if err != nil {
			t.Fatal(err)
		}
		want[i], err = ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
		err = f.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

	rand.Seed(1)
	// Overwrite the Golden Images.
	ExampleFunc()

	// Read the images we've just generated and check them against the
	// Golden Images.
	for i, path := range paths {
		errored := false
		f, err := os.Open(filepath.Join("testdata", path))
		if err != nil {
			t.Error(err)
			errored = true
		}
		have, err := ioutil.ReadAll(f)
		if err != nil {
			t.Error(err)
			errored = true
		}
		err = f.Close()
		if err != nil {
			t.Error(err)
			errored = true
		}
		if !bytes.Equal(have, want[i]) {
			t.Errorf("image mismatch for %s\n", path)
			errored = true
		}
		if errored {
			// If there has been an error, write out the golden image for inspection.
			ext := filepath.Ext(path)
			noext := strings.TrimSuffix(path, ext)
			goldenPath := noext + "_golden" + ext
			f, err := os.Create(filepath.Join("testdata", goldenPath))
			if err != nil {
				t.Fatal(err)
			}
			_, err = f.Write(want[i])
			if err != nil {
				t.Fatal(err)
			}
			err = f.Close()
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("golden image has been written out as %s", goldenPath)
		}
	}
}

// Draw the plot logo.
func Example() {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}

	DefaultLineStyle.Width = vg.Points(1)
	DefaultGlyphStyle.Radius = vg.Points(3)

	p.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{0, "0"}, {0.25, ""}, {0.5, "0.5"}, {0.75, ""}, {1, "1"},
	})
	p.X.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{0, "0"}, {0.25, ""}, {0.5, "0.5"}, {0.75, ""}, {1, "1"},
	})

	pts := XYs{{0, 0}, {0, 1}, {0.5, 1}, {0.5, 0.6}, {0, 0.6}}
	line, err := NewLine(pts)
	if err != nil {
		log.Panic(err)
	}
	scatter, err := NewScatter(pts)
	if err != nil {
		log.Panic(err)
	}
	p.Add(line, scatter)

	pts = XYs{{1, 0}, {0.75, 0}, {0.75, 0.75}}
	line, err = NewLine(pts)
	if err != nil {
		log.Panic(err)
	}
	scatter, err = NewScatter(pts)
	if err != nil {
		log.Panic(err)
	}
	p.Add(line, scatter)

	pts = XYs{{0.5, 0.5}, {1, 0.5}}
	line, err = NewLine(pts)
	if err != nil {
		log.Panic(err)
	}
	scatter, err = NewScatter(pts)
	if err != nil {
		log.Panic(err)
	}
	p.Add(line, scatter)

	err = p.Save(100, 100, "testdata/plotLogo.png")
	if err != nil {
		log.Panic(err)
	}
}

func TestMainExample(t *testing.T) {
	checkPlot(Example, t, "plotLogo.png")
}
