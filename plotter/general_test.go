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

func goldenPath(path string) string {
	ext := filepath.Ext(path)
	noext := strings.TrimSuffix(path, ext)
	return noext + "_golden" + ext
}

// checkPlot checks a generated plot against a previously created reference.
// If generateTestData = true, it regereates the reference.
func checkPlot(ExampleFunc func(), t *testing.T, filenames ...string) {
	paths := make([]string, len(filenames))
	for i, fn := range filenames {
		paths[i] = filepath.Join("testdata", fn)
	}

	if *generateTestData {
		// Recreate Golden images and exit.
		ExampleFunc()
		for _, path := range paths {
			golden := goldenPath(path)
			_ = os.Remove(golden)
			if err := os.Rename(path, golden); err != nil {
				t.Fatal(err)
			}
		}
		return
	}

	rand.Seed(1)
	// Overwrite the Golden Images.
	ExampleFunc()

	// Read the images we've just generated and check them against the
	// Golden Images.
	for _, path := range paths {
		got, err := ioutil.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read %s: %s", path, err)
			continue
		}
		golden := goldenPath(path)
		want, err := ioutil.ReadFile(golden)
		if err != nil {
			t.Errorf("Failed to read golden file %s: %s", golden, err)
			continue
		}
		if !bytes.Equal(got, want) {
			t.Errorf("image mismatch for %s\n", path)
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
