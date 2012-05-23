// vecimg implements the vecgfx.Canvas interface
// using the draw2d package as a backend to output
// raster images.
package vecimg

import (
	"bufio"
	"code.google.com/p/draw2d/draw2d"
	"code.google.com/p/plotinum/vecgfx"
	"fmt"
	"go/build"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

const (
	// dpi is the number of dots per inch.
	dpi = 96

	// importString is the current package import string.
	importString = "code.google.com/p/plotinum/vecgfx/vecimg"
)

type ImageCanvas struct {
	gc    draw2d.GraphicContext
	img   image.Image
	color color.Color
}

// New returns a new image canvas with the size specified in inches,
// rounded up to the nearest pixel.
func New(w, h float64) (*ImageCanvas, error) {
	pkg, err := build.Import(importString, "", build.FindOnly)
	if err != nil {
		return nil, err
	}
	draw2d.SetFontFolder(filepath.Join(pkg.Dir, "fonts"))

	w *= dpi
	h *= dpi
	img := image.NewRGBA(image.Rect(0, 0, int(w+0.5), int(h+0.5)))
	gc := draw2d.NewGraphicContext(img)
	gc.SetDPI(dpi)
	gc.Scale(1, -1)
	gc.Translate(0, -h)
	return &ImageCanvas{
		gc:    gc,
		img:   img,
		color: color.RGBA{A: 255},
	}, nil
}

func (c *ImageCanvas) SetLineWidth(w float64) {
	c.gc.SetLineWidth(w)
}

func (c *ImageCanvas) SetLineDash(ds []float64, offs float64) {
	c.gc.SetLineDash(ds, offs)
}

func (c *ImageCanvas) SetColor(color color.Color) {
	c.gc.SetFillColor(color)
	c.gc.SetStrokeColor(color)
	c.color = color
}

func (c *ImageCanvas) Rotate(t float64) {
	c.gc.Rotate(t)
}

func (c *ImageCanvas) Translate(x, y float64) {
	c.gc.Translate(x, y)
}

func (c *ImageCanvas) Scale(x, y float64) {
	c.gc.Scale(x, y)
}

func (c *ImageCanvas) Push() {
	c.gc.Save()
}

func (c *ImageCanvas) Pop() {
	c.gc.Restore()
}

func (c *ImageCanvas) Stroke(p vecgfx.Path) {
	c.outline(p)
	c.gc.Stroke()
}

func (c *ImageCanvas) Fill(p vecgfx.Path) {
	c.outline(p)
	c.gc.Fill()
}

func (c *ImageCanvas) outline(p vecgfx.Path) {
	c.gc.BeginPath()
	for _, comp := range p {
		switch comp.Type {
		case vecgfx.MoveComp:
			c.gc.MoveTo(comp.X, comp.Y)

		case vecgfx.LineComp:
			c.gc.LineTo(comp.X, comp.Y)

		case vecgfx.ArcComp:
			c.gc.ArcTo(comp.X, comp.Y, comp.Radius,
				comp.Radius, comp.Start, comp.Finish)

		case vecgfx.CloseComp:
			c.gc.Close()

		default:
			panic(fmt.Sprintf("Unknown path component: %d", comp.Type))
		}
	}
}

func (c *ImageCanvas) DPI() float64 {
	return float64(c.gc.GetDPI())
}

func (c *ImageCanvas) FillText(font vecgfx.Font, x, y float64, str string) {
	c.gc.Save()
	c.gc.Translate(x, y+font.Extents().Ascent/vecgfx.PtInch*c.DPI())
	c.gc.Scale(1, -1)
	c.gc.DrawImage(c.textImage(font, str))
	c.gc.Restore()
}

func (c *ImageCanvas) textImage(font vecgfx.Font, str string) *image.RGBA {
	w := font.Width(str) / vecgfx.PtInch * c.DPI()
	h := font.Extents().Height / vecgfx.PtInch * c.DPI()
	img := image.NewRGBA(image.Rect(0, 0, int(w+0.5), int(h+0.5)))
	gc := draw2d.NewGraphicContext(img)

	gc.SetDPI(int(c.DPI()))
	gc.SetFillColor(c.color)
	data, ok := fontMap[font.Name()]
	if !ok {
		panic(fmt.Sprintf("Font name %s is unknown", font.Name()))
	}

	gc.SetFontData(data)
	gc.SetFontSize(font.Size)
	gc.MoveTo(0, h+font.Extents().Descent)
	gc.FillString(str)

	return img
}

var (
	fontMap = map[string]draw2d.FontData{
		"Courier": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilyMono,
			Style:  draw2d.FontStyleNormal,
		},
		"Courier-Bold": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilyMono,
			Style:  draw2d.FontStyleBold,
		},
		"Courier-Oblique": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilyMono,
			Style:  draw2d.FontStyleItalic,
		},
		"Courier-BoldOblique": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilyMono,
			Style:  draw2d.FontStyleItalic | draw2d.FontStyleBold,
		},
		"Helvetica": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleNormal,
		},
		"Helvetica-Bold": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleBold,
		},
		"Helvetica-Oblique": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleItalic,
		},
		"Helvetica-BoldOblique": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleItalic | draw2d.FontStyleBold,
		},
		"Times-Roman": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySerif,
			Style:  draw2d.FontStyleNormal,
		},
		"Times-Bold": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySerif,
			Style:  draw2d.FontStyleBold,
		},
		"Times-Italic": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySerif,
			Style:  draw2d.FontStyleItalic,
		},
		"Times-BoldItalic": draw2d.FontData{
			Name:   "Nimbus",
			Family: draw2d.FontFamilySerif,
			Style:  draw2d.FontStyleItalic | draw2d.FontStyleBold,
		},
	}
)

func (c *ImageCanvas) SavePNG(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	err = png.Encode(b, c.img)
	if err != nil {
		return err
	}
	return b.Flush()
}
