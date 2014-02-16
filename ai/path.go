package ai

import (
	"github.com/g3force/Go_Sokoban/engine"
)

type Path []Node


/**
 * list of all directions
 */
func (path Path) Directions() []engine.Direction {
	pa := []engine.Direction{}
	for i := 0; i < len(path); i++ {
		pa = append(pa, path[i].Direction())
	}
	return pa
}

func (p *Path) Clone() (path Path) {
	path = make([]Node,len(*p))
	for k,v := range *p {
		path[k] = v.Clone()
	}
	return
}

func (path Path) Current() *Node {
	if len(path) > 0 {
		return &path[len(path)-1]
	}
	panic("Current(): path is empty")
}

func (path Path) CurrentDir() engine.Direction {
	if len(path) > 0 {
		return path.Current().Direction()
	}
	return engine.NO_DIRECTION
}

func (path *Path) SetCurrentDir(dir engine.Direction ) {
	path.Current().SetDirection(dir)
}

func (path *Path) IncCurrentDir() {
	path.Current().dir++
	if path.Current().dir == 4 {
		path.Current().dir = 0
	}
	path.Current().counter++
}

func (path Path) CurrentCounter() int8 {
	if len(path) > 0 {
		return path.Current().counter
	}
	return -1
}

func (p *Path) Pop() {
	path := *p
	if len(path) > 0 {
		path = path[:len(path)-1]
		*p = path
	} else {
		panic("Pop(): Path is empty. Can not remove from path.")
	}
}

func (p *Path) Push(dir engine.Direction ) {
	path := *p
	path = append(path, Node{-1, dir, nil})
	*p = path
}

func (path Path) Empty() bool {
	if len(path) == 0 {
		return true
	}
	return false
}
