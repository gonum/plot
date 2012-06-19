package plot

import (
	"code.google.com/p/plotinum/vg"
)

// A Legend gives a description of the meaning of different
// data elements of the plot.
type Legend struct {
	// Entries are all of the legendEntries described
	// by this legend.
	Entries []legendEntry

	// TextStyle is the style given to the legend
	// entry texts.
	TextStyle

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

	// IconWidth is the width of legend icons.
	IconWidth vg.Length
}

// A legendEntry represents a single line of a legend, it
// has a name and an icon.
type legendEntry struct {
	// text is the text associated with this entry. 
	text string

	// thumbs is a slice of all of the icon styles
	thumbs []thumbnailer
}

// MakeLegendEntry returns a legend entry with the given
// text which draws as the thumbnail the composition
// of all of the given thumbnailers.
func MakeLegendEntry(text string, thumbs ...thumbnailer) legendEntry {
	return legendEntry{
		text: text,
		thumbs: thumbs,
	}
}

// thumbnailer wraps the DrawIcon method
type thumbnailer interface {
	// Thumbnail draws an thumbnail representing
	// a legend entry.  The thumbnail will usually show
	// a smaller representation of the style used
	// to plot the corresponding data.
	Thumbnail(da *DrawArea)
}

// NewLegend returns a new legend containing the given
// entries.
func NewLegend(entries ...legendEntry) *Legend {
	font, err := vg.MakeFont(defaultFont, vg.Points(12))
	if err != nil {
		panic(err)
	}
	return &Legend {
		Entries: entries,
		IconWidth: vg.Points(20),
		TextStyle: TextStyle{ Font:  font },
	}
}

// draw draws the legend to the given DrawArea.
func (l *Legend) draw(da *DrawArea) {
	_, enth := l.textSize()
	iconx, txtx, xalign := l.xlocs(da)
	y := l.yloc(da)

	iconSz := Point{ l.IconWidth, enth }
	for _, e := range l.Entries {
		ico := &DrawArea{
			Canvas: da.Canvas,
			Rect: Rect{ Min: Point{ iconx, y }, Size: iconSz },
		}
		for _, t := range e.thumbs {
			t.Thumbnail(ico)
		}
		da.FillText(l.TextStyle, txtx, y, xalign, 0, e.text)
		y -= enth
	}
}

// xlocs returns the x location of the legend icons, text,
// and the text alignment.
func (l *Legend) xlocs(da *DrawArea) (iconx, txtx vg.Length, xalign float64){
	if l.Left {
		iconx = da.Min.X
		txtx = iconx + l.IconWidth + l.TextStyle.Width(" ")
	} else {
		txtw, _ := l.textSize()
		iconx = da.Max().X - l.IconWidth
		txtx = iconx - txtw - l.TextStyle.Width(" ")
	}
	txtx += l.XOffs
	iconx += l.XOffs
	return
}

// yloc returns the y location of the start of the legend.
func (l *Legend) yloc(da *DrawArea) vg.Length {
	_, enth := l.textSize()
	y := da.Max().Y - enth
	if !l.Top {
		y = da.Min.Y + enth*(vg.Length(len(l.Entries)) - 1)
	}
	y += l.YOffs
	return y
}

// textSize returns the width and height of the text fields
// of a legend.
func (l *Legend) textSize() (width, height vg.Length) {
	for _, e := range l.Entries {
		if w := l.TextStyle.Width(e.text); w > width {
			width = w
		}
		if h := l.TextStyle.Height(e.text); h > height {
			height = h
		}
	}
	return
}
