// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgl

import (
	"fmt"
	"image/color"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/gonum/plot/vg"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
)

// dpi is the number of dots per inch.
const dpi = 96

// Canvas implements the vg.Canvas interface, using OpenGL as a backend.
type Canvas struct {
	gc    draw2d.GraphicContext
	win   *glfw.Window
	w, h  vg.Length // width and height of the canvas
	width vg.Length // width of the current line
	color []color.Color
}

// make sure we implement vg.Canvas
var _ vg.Canvas = (*Canvas)(nil)

// channel where all OpenGL calls are made
var glchan chan func()
var done chan struct{}

func init() {
	// OpenGL calls need to be run on the main OS-thread.
	// calling LockOSThread during init() ensures that is the case.
	runtime.LockOSThread()
}

// Run runs OpenGL commands.
// Run needs to be run on the main OS-thread.
func Run(f func() error) error {
	err := gl.Init()
	if err != nil {
		return fmt.Errorf("vgl: could not initialize OpenGL context (%v)", err)
	}

	quit := make(chan struct{})

	done = make(chan struct{})
	glchan = make(chan func())
	go func() {
		err = f()
		quit <- struct{}{}
	}()

	for {
		select {
		case fct := <-glchan:
			fct()
			done <- struct{}{}
		case <-quit:
			return err
		}
	}

	return err
}

// call runs an OpenGL command on the main OS-thread
func call(fct func()) {
	glchan <- fct
	<-done
}

// New returns a new canvas with the size specified rounded up to the nearest
// pixel.
func New(w, h vg.Length, name string) (*Canvas, error) {
	ww := int(w/vg.Inch*dpi + 0.5)
	hh := int(h/vg.Inch*dpi + 0.5)
	var window *glfw.Window
	var err error
	call(func() {
		window, err = newWindow(ww, hh, name)
	})
	if err != nil {
		return nil, err
	}

	c := &Canvas{
		gc:    draw2dgl.NewGraphicContext(ww, hh),
		win:   window,
		w:     vg.Length(ww) / dpi * vg.Inch,
		h:     vg.Length(hh) / dpi * vg.Inch,
		color: []color.Color{color.White},
	}

	runtime.SetFinalizer(c, func(c *Canvas) {
		err := c.Close()
		if err != nil {
			panic(err)
		}
	})

	call(func() {
		c.win.MakeContextCurrent()

		// a blank canvas
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.Viewport(0, 0, int32(ww), int32(hh))
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		// map abstract coords directly to window coords.
		gl.Ortho(0, float64(ww), 0, float64(hh), -1, 1)
		// invert Y axis so increasing Y goes down
		c.gc.Scale(1, -1)
		// shift up origin to upper-left corner
		c.gc.Translate(0, -float64(hh))

	})

	//go c.run()
	vg.Initialize(c)
	return c, nil
}

func (c *Canvas) run() {
	for !c.win.ShouldClose() {
		call(c.paint)
		time.Sleep(10 * time.Millisecond)
	}
}

func (c *Canvas) paint() {
	//gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	c.gc.BeginPath()
	draw2d.RoundRect(c.gc, 200, 200, 600, 600, 100, 100)
	c.gc.SetFillColor(color.RGBA{0, 0, 0, 0xFF})
	c.gc.Fill()
	gl.Flush() // single buffered, so needs a flush
	c.win.SwapBuffers()
	glfw.PollEvents()
	fmt.Printf("done\n")
}

func (c *Canvas) Paint() {
	call(glfw.WaitEvents)
	call(c.paint)
	call(glfw.WaitEvents)
	//c.run()
}

func (c *Canvas) DPI() float64 {
	return float64(c.gc.GetDPI())
}

func (c *Canvas) Size() (w, h vg.Length) {
	return c.w, c.h
}

func (c *Canvas) SetLineWidth(w vg.Length) {
	c.width = w
	c.gc.SetLineWidth(w.Dots(c))
}

func (c *Canvas) SetLineDash(ds []vg.Length, offs vg.Length) {
	dashes := make([]float64, 0, len(ds))
	for _, d := range ds {
		dashes = append(dashes, d.Dots(c))
	}
	c.gc.SetLineDash(dashes, offs.Dots(c))
}

func (c *Canvas) SetColor(clr color.Color) {
	if clr == nil {
		clr = color.Black
	}
	c.gc.SetFillColor(clr)
	c.gc.SetStrokeColor(clr)
	c.color[len(c.color)-1] = clr
}

func (c *Canvas) Rotate(r float64) {
	c.gc.Rotate(r)
}

func (c *Canvas) Translate(x, y vg.Length) {
	c.gc.Translate(x.Dots(c), y.Dots(c))
}

func (c *Canvas) Scale(x, y float64) {
	c.gc.Scale(x, y)
}

func (c *Canvas) Push() {
	c.color = append(c.color, c.color[len(c.color)-1])
	c.gc.Save()
}

func (c *Canvas) Pop() {
	c.color = c.color[:len(c.color)-1]
	c.gc.Restore()
}

func (c *Canvas) Stroke(p vg.Path) {
	if c.width <= 0 {
		return
	}
	c.outline(p)
	c.gc.Stroke()
}

func (c *Canvas) Fill(p vg.Path) {
	c.outline(p)
	c.gc.Fill()
}

func (c *Canvas) outline(p vg.Path) {
	c.gc.BeginPath()
	for _, comp := range p {
		switch comp.Type {
		case vg.MoveComp:
			c.gc.MoveTo(comp.X.Dots(c), comp.Y.Dots(c))

		case vg.LineComp:
			c.gc.LineTo(comp.X.Dots(c), comp.Y.Dots(c))

		case vg.ArcComp:
			c.gc.ArcTo(comp.X.Dots(c), comp.Y.Dots(c),
				comp.Radius.Dots(c), comp.Radius.Dots(c),
				comp.Start, comp.Angle)

		case vg.CloseComp:
			c.gc.Close()

		default:
			panic(fmt.Sprintf("Unknown path component: %d", comp.Type))
		}
	}
}

func (c *Canvas) FillString(font vg.Font, x, y vg.Length, str string) {
	c.gc.Save()
	defer c.gc.Restore()

	data, ok := fontMap[font.Name()]
	if !ok {
		panic(fmt.Sprintf("Font name %s is unknown", font.Name()))
	}
	if !registeredFont[font.Name()] {
		draw2d.RegisterFont(data, font.Font())
		registeredFont[font.Name()] = true
	}
	c.gc.SetFontData(data)
	c.gc.SetFontSize(font.Size.Points())
	c.gc.Translate(x.Dots(c), y.Dots(c))
	c.gc.Scale(1, -1)
	c.gc.FillString(str)
}

// Close releases all opengl contexts attached to this canvas.
func (c *Canvas) Close() error {
	if c.win != nil {
		c.win.Destroy()
		c.win = nil
	}
	return nil
}

var (
	// RegisteredFont contains the set of font names
	// that have already been registered with draw2d.
	registeredFont = map[string]bool{}

	// FontMap contains a mapping from vg's font
	// names to draw2d.FontData for the corresponding
	// font.  This is needed to register the  fonts with
	// draw2d.
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
