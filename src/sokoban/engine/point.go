package engine

import ()

// simple Point type
type Point struct {
	X int8
	Y int8
}

// add to points (their x and y)
func (p1 *Point) Add(p2 Point) Point {
	var newP Point
	newP.X = p1.X + p2.X
	newP.Y = p1.Y + p2.Y
	return newP
}

func NewPoint(x int, y int) Point {
	return Point{int8(x), int8(y)}
}

func NewPoint8(x int8, y int8) Point {
	return Point{x, y}
}

func (p *Point) Clone() (Point) {
	return NewPoint8((*p).X, (*p).Y)
}