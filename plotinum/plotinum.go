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
