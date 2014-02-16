package engine

type Direction int8 // 0-3, -1 for invalid
const NO_DIRECTION = Direction(-1)

// convert direction from int to Point
func (dir Direction) Point() Point {
	dir = dir % 4
	var p Point
	switch dir {
	case 0: // right
		p.X = 1
		p.Y = 0
	case 1: // down
		p.X = 0
		p.Y = 1
	case 2: // left
		p.X = -1
		p.Y = 0
	case 3: // up
		p.X = 0
		p.Y = -1
	}
	return p
}

func (dir Direction) Int() int8 {
	return (int8) (dir)
}
