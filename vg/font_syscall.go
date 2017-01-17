// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !js

package vg

import (
	"go/build"
	"os"
	"path/filepath"
)

func init() {
	initFontDirs()
}

// InitFontDirs returns the initial value for the FontDirectories variable.
func initFontDirs() []string {
	dirs := filepath.SplitList(os.Getenv("VGFONTPATH"))

	if pkg, err := build.Import(importString, "", build.FindOnly); err == nil {
		p := filepath.Join(pkg.Dir, "fonts")
		if _, err := os.Stat(p); err == nil {
			dirs = append(dirs, p)
		}
	}

	if len(dirs) == 0 {
		dirs = []string{"./fonts"}
	}

	return dirs
}
