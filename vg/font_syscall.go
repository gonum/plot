// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !js

package vg

import (
	"os"
	"path"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

func init() {
	FontDirs = initFontDirs()
}

// initFontDirs returns the initial value for the FontDirs variable.
func initFontDirs() []string {
	const fonts = "gonum.org/v1/plot/vg/fonts"

	dirs := filepath.SplitList(os.Getenv("VGFONTPATH"))

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles,
	}
	pkgs, err := packages.Load(cfg, fonts)
	if err == nil {
		dirs = append(dirs, path.Dir(pkgs[0].GoFiles[0]))
	}

	if len(dirs) == 0 {
		dirs = []string{"./fonts"}
	}

	return dirs
}
