// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"image/color"
	"math"

	xfont "golang.org/x/image/font"

	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// A Legend gives a description of the meaning of different
// data elements of the plot.  Each legend entry has a name
// and a thumbnail, where the thumbnail shows a small
// sample of the display style of the corresponding data.
type Legend struct {
	// BackgroundColor is the color behind the legend entries.
	// A nil color leaves the background transparent.
	BackgroundColor color.Color

	// TextStyle is the style given to the legend
	// entry texts.
	TextStyle text.Style

	// Padding is the amount of padding to add
	// between each entry in the legend.  If Padding
	// is zero then entries are spaced based on the
	// font size.
	Padding vg.Length

	// Top and Left specify the location of the legend.
	// If Top is true the legend is located along the top
	// edge of the plot, otherwise it is located along
	// the bottom edge.  If Left is true then the legend
	// is located along the left edge of the plot, and the
	// text is positioned after the icons, otherwise it is
	// located along the right edge and the text is
	// positioned before the icons.
	Top, Left bool

	// XOffs and YOffs are added to the legend's
	// final position.
	XOffs, YOffs vg.Length

	// YPosition specifies the vertical position of a legend entry.
	// Valid values are [-1,+1], with +1 being the top of the
	// entry vertical space, and -1 the bottom.
	YPosition float64

	// ThumbnailWidth is the width of legend thumbnails.
	ThumbnailWidth vg.Length

	// entries are all of the legendEntries described
	// by this legend.
	entries []legendEntry
}

// A legendEntry represents a single line of a legend, it
// has a name and an icon.
type legendEntry struct {
	// text is the text associated with this entry.
	text string

	// thumbs is a slice of all of the thumbnails styles
	thumbs []Thumbnailer
}

// Thumbnailer wraps the Thumbnail method, which
// draws the small image in a legend representing the
// style of data.
type Thumbnailer interface {
	// Thumbnail draws an thumbnail representing
	// a legend entry.  The thumbnail will usually show
	// a smaller representation of the style used
	// to plot the corresponding data.
	Thumbnail(c *draw.Canvas)
}

// NewLegend returns a legend with the default parameter settings.
func NewLegend() Legend {
	return newLegend(DefaultTextHandler)
}

func newLegend(hdlr text.Handler) Legend {
	return Legend{
		YPosition:      draw.PosBottom,
		ThumbnailWidth: vg.Points(20),
		TextStyle: text.Style{
			Font:    font.From(DefaultFont, 12),
			Handler: hdlr,
		},
	}
}

// Draw draws the legend to the given draw.Canvas.
func (l *Legend) Draw(c draw.Canvas) {
	if l.BackgroundColor != nil && len(l.entries) != 0 {
		r := l.Rectangle(c)
		c.SetColor(l.BackgroundColor)
		c.Fill(r.Path())
	}
	l.layout(c, func(sty text.Style, pt vg.Point, icon *draw.Canvas, e legendEntry) {
		for _, t := range e.thumbs {
			t.Thumbnail(icon)
		}
		c.FillText(sty, pt, e.text)
	})
}

// layout visits the entries at the positions used to draw them.
func (l *Legend) layout(c draw.Canvas, visit func(text.Style, vg.Point, *draw.Canvas, legendEntry)) {
	iconx := c.Min.X
	sty := l.TextStyle
	em := sty.Rectangle(" ")
	textx := iconx + l.ThumbnailWidth + em.Max.X
	if !l.Left {
		iconx = c.Max.X - l.ThumbnailWidth
		textx = iconx - em.Max.X
		sty.XAlign--
	}
	textx += l.XOffs
	iconx += l.XOffs

	descent := sty.FontExtents().Descent
	enth := l.entryHeight()
	y := c.Max.Y - enth - descent
	if !l.Top {
		y = c.Min.Y + (enth+l.Padding)*(vg.Length(len(l.entries))-1)
	}
	y += l.YOffs

	icon := &draw.Canvas{
		Canvas: c.Canvas,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: iconx, Y: y},
			Max: vg.Point{X: iconx + l.ThumbnailWidth, Y: y + enth},
		},
	}

	if l.YPosition < draw.PosBottom || draw.PosTop < l.YPosition {
		panic("plot: invalid vertical offset for the legend's entries")
	}
	yoff := vg.Length(l.YPosition-draw.PosBottom) / 2
	yoff += descent

	for _, e := range l.entries {
		yoffs := (enth - descent - sty.Rectangle(e.text).Max.Y) / 2
		yoffs += yoff
		visit(sty, vg.Point{X: textx, Y: icon.Min.Y + yoffs}, icon, e)
		icon.Min.Y -= enth + l.Padding
		icon.Max.Y -= enth + l.Padding
	}
}

// Rectangle returns the extent of the Legend.
func (l *Legend) Rectangle(c draw.Canvas) vg.Rectangle {
	r := vg.Rectangle{
		Min: vg.Point{X: vg.Length(math.Inf(1)), Y: vg.Length(math.Inf(1))},
		Max: vg.Point{X: vg.Length(math.Inf(-1)), Y: vg.Length(math.Inf(-1))},
	}
	include := func(box vg.Rectangle) {
		r.Min.X, r.Min.Y = min(r.Min.X, box.Min.X), min(r.Min.Y, box.Min.Y)
		r.Max.X, r.Max.Y = max(r.Max.X, box.Max.X), max(r.Max.Y, box.Max.Y)
	}
	l.layout(c, func(sty text.Style, pt vg.Point, icon *draw.Canvas, e legendEntry) {
		include(icon.Rectangle)
		include(legendTextRectangle(sty, e.text).Add(pt))
	})
	if len(l.entries) == 0 {
		p := c.Min
		if !l.Left {
			p.X = c.Max.X
		}
		if l.Top {
			p.Y = c.Max.Y
		}
		p = p.Add(vg.Point{X: l.XOffs, Y: l.YOffs})
		return vg.Rectangle{Min: p, Max: p}
	}
	return r
}

// legendTextRectangle includes glyph overhangs and the line positions used by
// the plain text handler, as well as the handler's layout bounds.
func legendTextRectangle(sty text.Style, txt string) vg.Rectangle {
	r := sty.Rectangle(txt)
	switch sty.Handler.(type) {
	case text.Plain, *text.Plain:
	default:
		return r
	}
	if sty.Font.Size <= 0 {
		return r
	}
	face := sty.Handler.Cache().Lookup(sty.Font, sty.Font.Size)
	// Measure at the font's design resolution to avoid pixel rounding.
	em := float64(face.Face.UnitsPerEm())
	glyphs := face.FontFace(72 * em / face.Font.Size.Points())
	defer glyphs.Close()
	scale := face.Font.Size / vg.Length(64*em)
	lines := sty.Handler.Lines(txt)
	y := sty.Height(txt)*vg.Length(sty.YAlign) - face.Extents().Ascent
	sin, cos := math.Sincos(sty.Rotation)
	for i, line := range lines {
		if line == "" {
			continue
		}
		bounds, _ := xfont.BoundString(glyphs, line)
		xoff := vg.Length(sty.XAlign) * face.Width(line)
		yoff := y + vg.Length(len(lines)-i)*sty.Font.Size
		for _, x := range []vg.Length{vg.Length(bounds.Min.X) * scale, vg.Length(bounds.Max.X) * scale} {
			for _, y := range []vg.Length{-vg.Length(bounds.Max.Y) * scale, -vg.Length(bounds.Min.Y) * scale} {
				x, y := x+xoff, y+yoff
				p := vg.Point{X: x*vg.Length(cos) - y*vg.Length(sin), Y: x*vg.Length(sin) + y*vg.Length(cos)}
				r.Min.X, r.Min.Y = min(r.Min.X, p.X), min(r.Min.Y, p.Y)
				r.Max.X, r.Max.Y = max(r.Max.X, p.X), max(r.Max.Y, p.Y)
			}
		}
	}
	return r
}

// entryHeight returns the height of the tallest legend
// entry text.
func (l *Legend) entryHeight() (height vg.Length) {
	for _, e := range l.entries {
		if h := l.TextStyle.Rectangle(e.text).Max.Y; h > height {
			height = h
		}
	}
	return
}

// Add adds an entry to the legend with the given name.
// The entry's thumbnail is drawn as the composite of all of the
// thumbnails.
func (l *Legend) Add(name string, thumbs ...Thumbnailer) {
	l.entries = append(l.entries, legendEntry{text: name, thumbs: thumbs})
}
