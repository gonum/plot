// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

// Some of this code (namely the code for computing the
// width of a string in a given font) was copied from
// code.google.com/p/freetype-go/freetype/ which includes
// the following copyright notice:
// Copyright 2010 The Freetype-Go Authors. All rights reserved.

package vg

import (
	"errors"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

const (
	// importString is the import string expected for
	// this package.  It is used to find the font
	// directory included with the package source.
	importString = "code.google.com/p/plotinum/vg"
)

var (
	// FontMap maps Postscript/PDF font names to compatible
	// free fonts (TrueType converted ghostscript fonts).
	// Fonts that are not keys of this map are not supported.
	FontMap = map[string]string{

		// At the moment, we use fonts from GNU's freefont
		// project.

		"Courier":             "NimbusMonL-Regu",
		"Courier-Bold":        "NimbusMonL-Bold",
		"Courier-Oblique":     "NimbusMonL-ReguObli",
		"Courier-BoldOblique": "NimbusMonL-BoldObli",

		"Helvetica":             "NimbusSanL-Regu",
		"Helvetica-Bold":        "NimbusSanL-Bold",
		"Helvetica-Oblique":     "NimbusSanL-ReguItal",
		"Helvetica-BoldOblique": "NimbusSanL-BoldItal",

		"Times-Roman":      "NimbusRomNo9L-Regu",
		"Times-Bold":       "NimbusRomNo9L-Medi",
		"Times-Italic":     "NimbusRomNo9L-ReguItal",
		"Times-BoldItalic": "NimbusRomNo9L-MediItal",
	}

	// loadedFonts is indexed by a font name and it
	// caches the associated *truetype.Font.
	loadedFonts = make(map[string]*truetype.Font)
)

// A Font represents one of the supported font
// faces.
type Font struct {
	// Size is the size of the font.  The font size can
	// be used as a reasonable value for the horizontal
	// distance between two successive lines of text.
	Size Length

	// name is the name of this font.
	name string

	// This is a little bit of a hack, but the truetype
	// font is currently only needed to determine the
	// dimensions of strings drawn in this font.
	// The actual drawing of the strings is handled
	// separately by different back-ends:
	// Both Postscript and PDF are capable of drawing
	// their own fonts and draw2d loads its own copy of
	// the truetype fonts for its own output.
	//
	// This isn't a necessity--some future backend is
	// free to use this field--however it is a consequence
	// of the fact that the current backends were
	// developed independently of this package.

	// font is the truetype font pointer for this
	// font.
	font *truetype.Font
}

// MakeFont returns a font object.  The name of the font must
// be a key of the FontMap.  The font file is located by searching
// the FontDirs slice for a directory containing the relevant font
// file.  The font file name is name mapped by FontMap with the
// .ttf extension.  For example, the font file for the font name
// Courier is NimbusMonL-Regu.ttf.
func MakeFont(name string, size Length) (font Font, err error) {
	font.Size = size
	font.name = name
	font.font, err = getFont(name)
	return
}

// Name returns the name of the font.
func (f *Font) Name() string {
	return f.name
}

// Font returns the corresponding truetype.Font.
func (f *Font) Font() *truetype.Font {
	return f.font
}

// SetName sets the name of the font, effectively
// changing the font.  If an error is returned then
// the font is left unchanged.
func (f *Font) SetName(name string) error {
	font, err := getFont(name)
	if err != nil {
		return err
	}
	f.name = name
	f.font = font
	return nil
}

// FontExtents contains font metric information.
type FontExtents struct {
	// Ascent is the distance that the text
	// extends above the baseline.
	Ascent Length

	// Descent is the distance that the text
	// extends below the baseline.  The descent
	// is given as a negative value.
	Descent Length

	// Height is the distance from the lowest
	// descending point to the highest ascending
	// point.
	Height Length
}

// Extents returns the FontExtents for a font.
func (f *Font) Extents() FontExtents {
	bounds := f.font.Bounds(f.Font().FUnitsPerEm())
	scale := f.Size / Points(float64(f.Font().FUnitsPerEm()))
	return FontExtents{
		Ascent:  Points(float64(bounds.YMax)) * scale,
		Descent: Points(float64(bounds.YMin)) * scale,
		Height:  Points(float64(bounds.YMax-bounds.YMin)) * scale,
	}
}

// Width returns width of a string when drawn using the font.
func (f *Font) Width(s string) Length {
	// scale converts truetype.FUnit to float64
	scale := f.Size / Points(float64(f.font.FUnitsPerEm()))

	width := 0
	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := f.font.Index(rune)
		if hasPrev {
			width += int(f.font.Kerning(f.font.FUnitsPerEm(), prev, index))
		}
		width += int(f.font.HMetric(f.font.FUnitsPerEm(), index).AdvanceWidth)
		prev, hasPrev = index, true
	}
	return Points(float64(width)) * scale
}

// getFont returns the truetype.Font for the given font name or an error.
func getFont(name string) (*truetype.Font, error) {
	if f, ok := loadedFonts[name]; ok {
		return f, nil
	}

	path, err := fontPath(name)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("Failed to open font file: " + err.Error())
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.New("Failed to read font file: " + err.Error())
	}

	font, err := freetype.ParseFont(bytes)
	if err == nil {
		loadedFonts[name] = font
	} else {
		err = errors.New("Failed to parse font file: " + err.Error())
	}

	return font, err
}

// FontPath returns the path for a font name or an error if it is not found.
func fontPath(name string) (string, error) {
	fname, err := fontFile(name)
	if err != nil {
		return "", err
	}

	for _, d := range FontDirs {
		p := filepath.Join(d, fname)
		if _, err := os.Stat(p); err != nil {
			continue
		}
		return p, nil
	}

	return "", errors.New("Failed to locate a font file " + fname + " for font name " + name)
}

// FontDirs is a slice of directories searched for font data files.
// If the first font file found is unreadable or cannot be parsed, then
// subsequent directories are not tried, and the font will fail to load.
//
// The default slice contains, in the following order, the values of the
// environment variable VGFONTPATH if it is defined, then the vg
// source fonts directory if it is found (i.e., if vg was installed by
// go get).  If the resulting FontDirs slice is empty then the current
// directory is added to it.  This slice may be changed to load fonts
// from different locations.
var FontDirs = initFontDirs()

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

// FontFile returns the font file name for a font name or an error
// if it is an unknown font (i.e., not in the FontMap).
func fontFile(name string) (string, error) {
	var err error
	n, ok := FontMap[name]
	if !ok {
		errStr := "Unknown font: " + name + ".  Available fonts are:"
		for n := range FontMap {
			errStr += " " + n
		}
		err = errors.New(errStr)
	}
	return n + ".ttf", err
}
