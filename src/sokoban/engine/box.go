package engine

import ()

type Box struct {
	Pos   Point
	Order int8
}

func NewBox(pos Point, order int8) Box {
	return Box{pos, order}
}

func (b *Box) Clone() (box Box) {
	return Box{b.Pos.Clone(), b.Order}
}

func (b *Box) SetPos(p Point) {
	b.Pos = p
}

func (b *Box) SetOrder(order int8) {
	b.Order = order
}
