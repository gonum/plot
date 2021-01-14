// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Some of this code (namely the code for computing the
// width of a string in a given font) was copied from
// github.com/golang/freetype/ which includes
// the following copyright notice:
// Copyright 2010 The Freetype-Go Authors. All rights reserved.

package vg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"gonum.org/v1/plot/vg/fonts"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

var (
	// FontMap maps Postscript/PDF font names to compatible
	// free fonts (OpenType converted ghostscript fonts).
	// Fonts that are not keys of this map are not supported.
	FontMap = map[string]string{

		// We use fonts from RedHat's Liberation project:
		//  https://fedorahosted.org/liberation-fonts/

		"Courier":             "LiberationMono-Regular",
		"Courier-Bold":        "LiberationMono-Bold",
		"Courier-Oblique":     "LiberationMono-Italic",
		"Courier-BoldOblique": "LiberationMono-BoldItalic",

		"Helvetica":             "LiberationSans-Regular",
		"Helvetica-Bold":        "LiberationSans-Bold",
		"Helvetica-Oblique":     "LiberationSans-Italic",
		"Helvetica-BoldOblique": "LiberationSans-BoldItalic",

		"Times-Roman":      "LiberationSerif-Regular",
		"Times-Bold":       "LiberationSerif-Bold",
		"Times-Italic":     "LiberationSerif-Italic",
		"Times-BoldItalic": "LiberationSerif-BoldItalic",

		"LiberationMono-Regular":    "LiberationMono-Regular",
		"LiberationMono-Bold":       "LiberationMono-Bold",
		"LiberationMono-Italic":     "LiberationMono-Italic",
		"LiberationMono-BoldItalic": "LiberationMono-BoldItalic",

		"LiberationSans-Regular":    "LiberationSans-Regular",
		"LiberationSans-Bold":       "LiberationSans-Bold",
		"LiberationSans-Italic":     "LiberationSans-Italic",
		"LiberationSans-BoldItalic": "LiberationSans-BoldItalic",

		"LiberationSerif-Regular":    "LiberationSerif-Regular",
		"LiberationSerif-Bold":       "LiberationSerif-Bold",
		"LiberationSerif-Italic":     "LiberationSerif-Italic",
		"LiberationSerif-BoldItalic": "LiberationSerif-BoldItalic",
	}

	// loadedFonts is indexed by a font name and it
	// caches the associated *opentype.Font.
	loadedFonts = make(map[string]*opentype.Font)

	// FontLock protects access to the loadedFonts map.
	fontLock sync.RWMutex

	// default hinting for OpenType fonts
	defaultHinting = font.HintingNone
)

// A Font represents one of the supported font
// faces.
type Font struct {
	// Size is the size of the font.  The font size can
	// be used as a reasonable value for the vertical
	// distance between two successive lines of text.
	Size Length

	// name is the name of this font.
	name string

	// This is a little bit of a hack, but the opentype
	// font is currently only needed to determine the
	// dimensions of strings drawn in this font.
	// The actual drawing of the strings is handled
	// separately by different back-ends:
	// Both Postscript and PDF are capable of drawing
	// their own fonts and Gio loads its own copy of
	// the opentype fonts for its own output.
	//
	// This isn't a necessity--some future backend is
	// free to use this field--however it is a consequence
	// of the fact that the current backends were
	// developed independently of this package.

	// font is the opentype font pointer for this font.
	font *opentype.Font

	// hinting specifies how a vector font's glyph nodes is quantized.
	hinting font.Hinting
}

// MakeFont returns a font object.  The name of the font must
// be a key of the FontMap.  The font file is located by searching
// the FontDirs slice for a directory containing the relevant font
// file.  The font file name is name mapped by FontMap with the
// .ttf extension.  For example, the font file for the font name
// Courier is LiberationMono-Regular.ttf.
func MakeFont(name string, size Length) (font Font, err error) {
	font.Size = size
	font.name = name
	font.font, err = getFont(name)
	font.hinting = defaultHinting
	return
}

// Name returns the name of the font.
func (f *Font) Name() string {
	return f.name
}

// Font returns the corresponding opentype.Font.
func (f *Font) Font() *opentype.Font {
	return f.font
}

func (f *Font) FontFace(dpi float64) font.Face {
	face, err := opentype.NewFace(f.font, &opentype.FaceOptions{
		Size: f.Size.Points(),
		DPI:  dpi,
	})
	if err != nil {
		panic(err)
	}
	return face
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
	// is given as a positive value.
	Descent Length

	// Height is the distance from the lowest
	// descending point to the highest ascending
	// point.
	Height Length
}

// Extents returns the FontExtents for a font.
func (f *Font) Extents() FontExtents {
	var (
		// TODO(sbinet): re-use a Font-level sfnt.Buffer instead?
		buf  sfnt.Buffer
		ppem = fixed.Int26_6(f.font.UnitsPerEm())
	)

	met, err := f.font.Metrics(&buf, ppem, f.hinting)
	if err != nil {
		panic(fmt.Errorf("could not extract font extents: %v", err))
	}
	scale := f.Size / Points(float64(ppem))
	return FontExtents{
		Ascent:  Points(float64(met.Ascent)) * scale,
		Descent: Points(float64(met.Descent)) * scale,
		Height:  Points(float64(met.Height)) * scale,
	}
}

// Width returns width of a string when drawn using the font.
func (f *Font) Width(s string) Length {
	var (
		pixelsPerEm = fixed.Int26_6(f.font.UnitsPerEm())

		// scale converts sfnt.Unit to float64
		scale = f.Size / Points(float64(pixelsPerEm))

		width     = 0
		hasPrev   = false
		buf       sfnt.Buffer
		prev, idx sfnt.GlyphIndex
	)
	for _, rune := range s {
		var err error
		idx, err = f.font.GlyphIndex(&buf, rune)
		if err != nil {
			panic(fmt.Errorf("could not get glyph index: %v", err))
		}
		if hasPrev {
			kern, err := f.font.Kern(&buf, prev, idx, pixelsPerEm, f.hinting)
			switch {
			case err == nil:
				width += int(kern)
			case errors.Is(err, sfnt.ErrNotFound):
				// no-op
			default:
				panic(fmt.Errorf("could not get kerning: %v", err))
			}
		}
		adv, err := f.font.GlyphAdvance(&buf, idx, pixelsPerEm, f.hinting)
		if err != nil {
			panic(fmt.Errorf("could not retrieve glyph's advance: %v", err))
		}
		width += int(adv)
		prev, hasPrev = idx, true
	}
	return Points(float64(width)) * scale
}

// AddFont associates an opentype.Font with the given name.
func AddFont(name string, font *opentype.Font) {
	fontLock.Lock()
	loadedFonts[name] = font
	fontLock.Unlock()
}

// getFont returns the opentype.Font for the given font name or an error.
func getFont(name string) (*opentype.Font, error) {
	fontLock.RLock()
	f, ok := loadedFonts[name]
	fontLock.RUnlock()
	if ok {
		return f, nil
	}

	bytes, err := fontData(name)
	if err != nil {
		return nil, err
	}

	font, err := sfnt.Parse(bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse font file %q: %w", name, err)
	}

	fontLock.Lock()
	loadedFonts[name] = font
	fontLock.Unlock()

	return font, err
}

// fontData returns the []byte data for a font name or an error if it is not found.
func fontData(name string) ([]byte, error) {
	fname, err := fontFile(name)
	if err != nil {
		return nil, err
	}

	for _, d := range FontDirs {
		p := filepath.Join(d, fname)
		data, err := ioutil.ReadFile(p)
		if err != nil {
			continue
		}
		return data, nil
	}

	data, err := fonts.Asset(fname)
	if err == nil {
		return data, nil
	}

	return nil, errors.New("vg: failed to locate a font file " + fname + " for font name " + name)
}

// FontDirs is a slice of directories searched for font data files.
// If the first font file found is unreadable or cannot be parsed, then
// subsequent directories are not tried, and the font will fail to load.
//
// The default slice is initialised with the contents of the VGFONTPATH
// environment variable if it is defined.
// This slice may be changed to load fonts from different locations.
var FontDirs []string

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
