// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate ./list

package plotter

import (
	"math"
	"sort"

	"github.com/gonum/plot"
	"github.com/gonum/plot/palette"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// Contour implements the Plotter interface, drawing
// a contour plot of the values in the GridFunc field. The
// order of rows is from high index to low index down.
type Contour struct {
	GridFunc GridFunc

	// Levels describes the contour heights to plot.
	Levels []float64

	// LineStyle is the style of the contour lines.
	draw.LineStyle

	// Palette is the color palette used to render
	// the heat map.
	Palette palette.Palette

	// Min and Max define the dynamic range of the
	// heat map.
	Min, Max float64
}

// NewContour creates as new contour plotter for the given data,
// using the provided palette.
func NewContour(g GridFunc, levels []float64, p palette.Palette) *Contour {
	var min, max float64
	type minMaxer interface {
		Min() float64
		Max() float64
	}
	switch g := g.(type) {
	case minMaxer:
		min, max = g.Min(), g.Max()
	default:
		min, max = math.Inf(1), math.Inf(-1)
		c, r := g.Dims()
		for i := 0; i < c; i++ {
			for j := 0; j < r; j++ {
				v := g.Z(i, j)
				if math.IsNaN(v) {
					continue
				}
				min = math.Min(min, v)
				max = math.Max(max, v)
			}
		}
	}

	return &Contour{
		GridFunc:  g,
		Levels:    levels,
		LineStyle: DefaultLineStyle,
		Palette:   p,
		Min:       min,
		Max:       max,
	}
}

const naive = false

// Plot implements the Plot method of the plot.Plotter interface.
func (h *Contour) Plot(c draw.Canvas, plt *plot.Plot) {
	pal := h.Palette.Colors()
	if len(pal) == 0 {
		panic("contour: empty palette")
	}

	c.SetLineStyle(h.LineStyle)

	trX, trY := plt.Transforms(&c)

	// Collate contour paths and draw them.
	//
	// The alternative naive approach is to draw each line segment as
	// conrec returns it. The integrated path approach allows graphical
	// optimisations and is necessary for contour fill shading.
	if !naive {
		cp := contourPaths(h.GridFunc, h.Levels, trX, trY)
		ps := float64(len(pal)-1) / (h.Levels[len(h.Levels)-1] - h.Levels[0])
		if len(h.Levels) == 1 {
			ps = 0
		}

		for _, z := range h.Levels {
			if math.IsNaN(z) {
				continue
			}
			for _, pa := range cp[z] {
				if isLoop(pa) {
					pa.Close()
				}
				col := pal[int((z-h.Levels[0])*ps+0.5)]
				if col != nil {
					c.SetColor(col)
					c.Stroke(pa)
				}
			}
		}
	} else {
		var pa vg.Path
		sort.Float64s(h.Levels)
		ps := float64(len(pal)-1) / (h.Levels[len(h.Levels)-1] - h.Levels[0])
		if len(h.Levels) == 1 {
			ps = 0
		}

		conrec(h.GridFunc, h.Levels, func(_, _ int, l line, z float64) {
			if math.IsNaN(z) {
				return
			}

			pa = pa[:0]

			x1, y1 := trX(l.p1.X), trY(l.p1.Y)
			x2, y2 := trX(l.p2.X), trY(l.p2.Y)

			if !c.Contains(draw.Point{x1, y1}) || !c.Contains(draw.Point{x2, y2}) {
				return
			}

			pa.Move(x1, y1)
			pa.Line(x2, y2)
			pa.Close()

			col := pal[int((z-h.Levels[0])*ps+0.5)]
			if col != nil {
				c.SetColor(col)
				c.Stroke(pa)
			}
		})
	}
}

// DataRange implements the DataRange method
// of the plot.DataRanger interface.
func (h *Contour) DataRange() (xmin, xmax, ymin, ymax float64) {
	c, r := h.GridFunc.Dims()
	return h.GridFunc.X(0), h.GridFunc.X(c - 1), h.GridFunc.Y(0), h.GridFunc.Y(r - 1)
}

// GlyphBoxes implements the GlyphBoxes method
// of the plot.GlyphBoxer interface.
func (h *Contour) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	c, r := h.GridFunc.Dims()
	b := make([]plot.GlyphBox, 0, r*c)
	for i := 0; i < c; i++ {
		for j := 0; j < r; j++ {
			b = append(b, plot.GlyphBox{
				X: plt.X.Norm(h.GridFunc.X(i)),
				Y: plt.Y.Norm(h.GridFunc.Y(j)),
				Rect: draw.Rect{
					Min:  draw.Point{-2.5, -2.5},
					Size: draw.Point{5, 5},
				},
			})
		}
	}
	return b
}

// isLoop returns true iff a vg.Path is a closed loop.
func isLoop(p vg.Path) bool {
	s := p[0]
	e := p[len(p)-1]
	return s.X == e.X && s.Y == e.Y
}

// contourPaths returns a collection of vg.Paths describing contour lines based
// on the input data in m cut at the given levels. The trX and trY function
// are coordinate transforms. The returned map contains slices of paths keyed
// on the value of the contour level.
func contourPaths(m GridFunc, levels []float64, trX, trY func(float64) vg.Length) map[float64][]vg.Path {
	sort.Float64s(levels)

	ends := make(map[float64]endMap)
	conts := make(contourSet)
	conrec(m, levels, func(_, _ int, l line, z float64) {
		paths(l, z, ends, conts)
	})
	ends = nil

	// TODO(kortschak): Check that all non-loop paths have
	// both ends at boundary. If any end is not at a boundary
	// it must have a partner near by. Find this partner and join
	// the two conts by merging the near by ends at the mean
	// location. This operation is done level by level to ensure
	// close contours of different heights are not joined.
	// A partner should be a float error different end, but I
	// suspect that is is possible for a bi- or higher order
	// furcation so it may be that the path ends at middle node
	// of another path. This needs to be investigated.
	//
	// If no partner is found and not close to boundary, panic.

	// Excise loops from crossed paths.
	for c := range conts {
		c.exciseLoops(conts)
	}

	// Build vg.Paths.
	paths := make(map[float64][]vg.Path)
	for c := range conts {
		paths[c.z] = append(paths[c.z], c.path(trX, trY))
	}

	return paths
}

// contourSet hold a working collection of contours.
type contourSet map[*contour]struct{}

// endMap holds a working collection of available ends.
type endMap map[point]*contour

// paths extends a conrecLine function to build a set of conts that represent paths along
// contour lines. It is used as the engine for a closure where ends and conts are closed
// around in a conrecLine function, and l and z are the line and height values provided
// by conrec. At the end of a conrec call, conts will contain a map keyed on the set of
// paths with each value being the height of the contour.
func paths(l line, z float64, ends map[float64]endMap, conts contourSet) {
	zEnds, ok := ends[z]
	if !ok {
		zEnds = make(endMap)
		ends[z] = zEnds
		c := newContour(l, z)
		zEnds[l.p1] = c
		zEnds[l.p2] = c
		conts[c] = struct{}{}
		return
	}

	c1, ok1 := zEnds[l.p1]
	c2, ok2 := zEnds[l.p2]

	// New segment.
	if !ok1 && !ok2 {
		c := newContour(l, z)
		zEnds[l.p1] = c
		zEnds[l.p2] = c
		conts[c] = struct{}{}
		return
	}

	if ok1 {
		// Add l.p2 to end of l.p1's contour.
		if !c1.extend(l, zEnds) {
			panic("internal link")
		}
	} else if ok2 {
		// Add l.p1 to end of l.p2's contour.
		if !c2.extend(l, zEnds) {
			panic("internal link")
		}
	}

	if c1 == c2 {
		return
	}

	// Join conts.
	if ok1 && ok2 {
		if !c1.connect(c2, zEnds) {
			panic("internal link")
		}
		delete(conts, c2)
	}
}

// excise removes all the elements from e1 to e2 in the list
// and creates a new list holding these elements. e1 must
// precede e2 in the list, otherwise excise will panic.
// If e1 or e2 are not in the list, excise panics.
func (l *list) excise(e1, e2 *element) *list {
	if e1.list != l || e2.list != l {
		panic("contour: mismatched list element")
	}

	var l2 list

	p := e1.prev
	n := e2.next
	l2.root.next = e1
	l2.root.prev = e2
	e1.prev = &l2.root
	e2.next = &l2.root
	p.next = n
	n.prev = p

	for e := e1; e != &l2.root; e = e.next {
		if e.list == nil {
			panic("contour: e2 before e1")
		}
		e.list = &l2
		l2.len++
		l.len--
	}
	return &l2
}

type contour struct {
	l *list
	z float64
}

func newContour(l line, z float64) *contour {
	li := newList()
	li.PushFront(l.p1)
	li.PushBack(l.p2)
	return &contour{l: li, z: z}
}

func (c *contour) path(trX, trY func(float64) vg.Length) vg.Path {
	var pa vg.Path
	e := c.l.Front()
	p := e.Value
	pa.Move(trX(p.X), trY(p.Y))
	for e = e.Next(); e != nil; e = e.Next() {
		p = e.Value
		pa.Line(trX(p.X), trY(p.Y))
	}

	return pa
}

func (c *contour) front() point { return c.l.Front().Value }
func (c *contour) back() point  { return c.l.Back().Value }

func (c *contour) extend(l line, ends endMap) (ok bool) {
	switch c.front() {
	case l.p1:
		c.l.PushFront(l.p2)
		delete(ends, l.p1)
		ends[l.p2] = c
		return true
	case l.p2:
		c.l.PushFront(l.p1)
		delete(ends, l.p2)
		ends[l.p1] = c
		return true
	}

	switch c.back() {
	case l.p1:
		c.l.PushBack(l.p2)
		delete(ends, l.p1)
		ends[l.p2] = c
		return true
	case l.p2:
		c.l.PushBack(l.p1)
		delete(ends, l.p2)
		ends[l.p1] = c
		return true
	}

	return false
}

func (c *contour) dropFront() { c.l.Remove(c.l.Front()) }
func (c *contour) dropBack()  { c.l.Remove(c.l.Back()) }
func (c *contour) connect(b *contour, ends endMap) bool {
	f1 := c.front()
	f2 := b.front()
	b1 := c.back()
	b2 := b.back()
	switch {
	case f1 == f2:
		delete(ends, f2)
		ends[b2] = c
		b.dropFront()
		for i, e := b.l.Len(), b.l.Front(); i > 0; i, e = i-1, e.Next() {
			c.l.PushFront(e.Value)
		}
		return true
	case f1 == b2:
		delete(ends, b2)
		ends[f2] = c
		b.dropBack()
		c.l.PushFrontList(b.l)
		ends[c.l.Front().Value] = c
		return true
	case b1 == f2:
		delete(ends, f2)
		ends[b2] = c
		b.dropFront()
		c.l.PushBackList(b.l)
		return true
	case b1 == b2:
		delete(ends, b2)
		ends[f2] = c
		b.dropBack()
		for i, e := b.l.Len(), b.l.Back(); i > 0; i, e = i-1, e.Prev() {
			c.l.PushBack(e.Value)
		}
		return true
	default:
		return false
	}
}

func (c *contour) exciseLoops(conts contourSet) {
	f := c.l.Front()
	seen := make(map[point]*element)
	for e := f; e != nil; e = e.Next() {
		if p, ok := seen[e.Value]; ok && e.Value != f.Value {
			nl := c.l.excise(p.Next(), e)
			nl.PushFront(e.Value)
			conts[&contour{l: nl, z: c.z}] = struct{}{}
		}
		seen[e.Value] = e
	}
}
