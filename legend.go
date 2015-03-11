// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// A Legend gives a description of the meaning of different
// data elements of the plot.  Each legend entry has a name
// and a thumbnail, where the thumbnail shows a small
// sample of the display style of the corresponding data.
type Legend struct {
	// TextStyle is the style given to the legend
	// entry texts.
	draw.TextStyle

	// Padding is the amount of padding to add
	// betweeneach entry of the legend.  If Padding
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

// makeLegend returns a legend with the default
// parameter settings.
func makeLegend() (Legend, error) {
	font, err := vg.MakeFont(DefaultFont, vg.Points(12))
	if err != nil {
		return Legend{}, err
	}
	return Legend{
		ThumbnailWidth: vg.Points(20),
		TextStyle:      draw.TextStyle{Font: font},
	}, nil
}

// draw draws the legend to the given draw.Canvas.
func (l *Legend) draw(c draw.Canvas) {
	iconx := c.Min.X
	textx := iconx + l.ThumbnailWidth + l.TextStyle.Width(" ")
	xalign := 0.0
	if !l.Left {
		iconx = c.Max.X - l.ThumbnailWidth
		textx = iconx - l.TextStyle.Width(" ")
		xalign = -1
	}
	textx += l.XOffs
	iconx += l.XOffs

	enth := l.entryHeight()
	y := c.Max.Y - enth
	if !l.Top {
		y = c.Min.Y + (enth+l.Padding)*(vg.Length(len(l.entries))-1)
	}
	y += l.YOffs

	icon := &draw.Canvas{
		Canvas: c.Canvas,
		Rectangle: draw.Rectangle{
			Min: draw.Point{iconx, y},
			Max: draw.Point{iconx + l.ThumbnailWidth, y + enth},
		},
	}
	for _, e := range l.entries {
		for _, t := range e.thumbs {
			t.Thumbnail(icon)
		}
		yoffs := (enth - l.TextStyle.Height(e.text)) / 2
		c.FillText(l.TextStyle, textx, icon.Min.Y+yoffs, xalign, 0, e.text)
		icon.Min.Y -= enth + l.Padding
		icon.Max.Y -= enth + l.Padding
	}
}

// entryHeight returns the height of the tallest legend
// entry text.
func (l *Legend) entryHeight() (height vg.Length) {
	for _, e := range l.entries {
		if h := l.TextStyle.Height(e.text); h > height {
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
