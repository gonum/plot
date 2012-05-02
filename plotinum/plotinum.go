package plotinum

type Rect struct {
	Min, Max Point
}

type Point struct {
	X, Y float64
}

func (p Point) Dot(q Point) float64 {
	return p.X*q.X + p.Y*q.Y
}

func (p Point) Plus(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func (p Point) Minus(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

func (p Point) Scale(s float64) Point {
	return Point{p.X * s, p.Y * s}
}

type Xer interface {
	X(int) float64
	Len() int
}

type Yer interface {
	Y(int) float64
	Len() int
}

type XYer interface {
	X(int) float64
	Y(int) float64
	Len() int
}

// Pt retuns the ith point of the given XYer.
func Pt(pts XYer, i int) Point {
	return Point{ X: pts.X(i), Y: pts.Y(i) }
}

// Line implements a simple XYer using a
// slice of Points.
type Line []Point

func (l Line) X(i int) float64 {
	return l[i].X
}

func (l Line) Y(i int) float64 {
	return l[i].Y
}

func (l Line) Len() int {
	return len(l)
}
