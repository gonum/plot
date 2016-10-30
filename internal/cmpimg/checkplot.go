// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmpimg

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var GenerateTestData = flag.Bool("regen", false, "Uses the current state to regenerate the test data.")

func goldenPath(path string) string {
	ext := filepath.Ext(path)
	noext := strings.TrimSuffix(path, ext)
	return noext + "_golden" + ext
}

// CheckPlot checks a generated plot against a previously created reference.
// If generateTestData = true, it regereates the reference.
func CheckPlot(ExampleFunc func(), t *testing.T, filenames ...string) {
	paths := make([]string, len(filenames))
	for i, fn := range filenames {
		paths[i] = filepath.Join("testdata", fn)
	}

	if *GenerateTestData {
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

	// Run the example.
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
		typ := filepath.Ext(path)[1:] // remove the dot in e.g. ".pdf"
		ok, err := Equal(typ, got, want)
		if err != nil {
			t.Errorf("failed to compare image for %s: %v\n", path, err)
			continue
		}
		if !ok {
			t.Errorf("image mismatch for %s\n", path)
			continue
		}
	}
}
