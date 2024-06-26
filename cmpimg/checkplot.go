// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmpimg

import (
	"bytes"
	"encoding/base64"
	"flag"
	"image"
	"image/png"
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
// If GenerateTestData = true, it regenerates the reference.
// For image.Image formats, a base64 encoded png representation is output to
// the testing log when a difference is identified.
func CheckPlot(ExampleFunc func(), t *testing.T, filenames ...string) {
	t.Helper()

	CheckPlotApprox(ExampleFunc, t, 0, filenames...)
}

// CheckPlotApprox checks a generated plot against a previously created reference.
// The normalized delta parameter describes how tight the matching should be
// performed, where delta=0 expresses a perfect match, and delta=1 a very loose match.
// If GenerateTestData = true, it regenerates the reference.
// For image.Image formats, a base64 encoded png representation is output to
// the testing log when a difference is identified.
func CheckPlotApprox(ExampleFunc func(), t *testing.T, delta float64, filenames ...string) {
	t.Helper()

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
		got, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read %s: %v", path, err)
			continue
		}
		golden := goldenPath(path)
		want, err := os.ReadFile(golden)
		if err != nil {
			t.Errorf("Failed to read golden file %s: %v", golden, err)
			continue
		}
		typ := filepath.Ext(path)[1:] // remove the dot in e.g. ".pdf"
		ok, err := EqualApprox(typ, got, want, delta)
		if err != nil {
			t.Errorf("failed to compare image for %s: %v", path, err)
			continue
		}
		if !ok {
			t.Errorf("image mismatch for %s\n", path)

			switch typ {
			case "jpeg", "jpg", "png", "tiff", "tif":
				v1, _, err := image.Decode(bytes.NewReader(got))
				if err != nil {
					t.Errorf("failed to decode %s: %v", path, err)
					continue
				}
				v2, _, err := image.Decode(bytes.NewReader(want))
				if err != nil {
					t.Errorf("failed to decode %s: %v", golden, err)
					continue
				}

				dst := image.NewRGBA64(v1.Bounds().Union(v2.Bounds()))
				rect := Diff(dst, v1, v2)
				t.Logf("image bounds union:%+v diff bounds intersection:%+v", dst.Bounds(), rect)

				var buf bytes.Buffer
				err = png.Encode(&buf, dst)
				if err != nil {
					t.Errorf("failed to encode difference png: %v", err)
					continue
				}
				t.Log("image diff:")
				t.Log("IMAGE:" + base64.StdEncoding.EncodeToString(buf.Bytes()))

				buf.Reset()
				err = png.Encode(&buf, v1)
				if err != nil {
					t.Errorf("failed to encode result png: %+v", err)
					continue
				}
				t.Log("image check:")
				t.Log("IMAGE:" + base64.StdEncoding.EncodeToString(buf.Bytes()))
			}
		}
	}
}
