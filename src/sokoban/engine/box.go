package engine

import ()

type Box struct {
	Pos   Point
	Order int8
}

func NewBox(pos Point, order int8) Box {
	return Box{pos, order}
}

func (b *Box) SetPos(p Point) {
	b.Pos = p
}

func (b *Box) SetOrder(order int8) {
	b.Order = order
}
