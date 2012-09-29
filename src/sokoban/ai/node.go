package ai

import (
	"sokoban/engine"
)

// indicates the direction, that the figure moves to 
// and counts the number of rotations
type Node struct {
	counter int8        // number of rotations
	dir     engine.Direction   // direction (0-3)
	ignored []engine.Direction // ignored directions during rotation
}

func NewNode() (node Node) {
	node.counter = -1
	node.dir = engine.NO_DIRECTION
	node.ignored = nil
	return 
}

func (node Node) Direction() engine.Direction {
	return node.dir
}

func (node *Node) SetDirection(dir engine.Direction) {
	node.dir = dir
}

func (node *Node) PushIgnored(dir engine.Direction) {
	node.ignored = append(node.ignored, dir)
}

//returns and deletes the ignored direction if there was a direction ignored in this Node, -1 if not
func (node *Node) PopIgnored() engine.Direction {
	if len(node.ignored) == 0 {
		return engine.NO_DIRECTION
	}
	// get
	dir := node.ignored[len(node.ignored)-1]
	// delete
	node.ignored = node.ignored[:len(node.ignored)-1]
	return dir
}
