package engine

import (
	"testing"
)

func TestCloneEngine(t *testing.T) {
	var e1 Engine
	var e2 Engine
	box1 := NewBox(Point{0,0}, 0)
	box2 := NewBox(Point{1,1}, 1)
	e1 = NewEngine()
	
	e1.figPos = Point{2,2}
	e1.boxes[0] = &box1
	e2 = e1.Clone()
	
	if e1.figPos.X != e2.figPos.X ||
		e1.figPos.Y != e2.figPos.Y {
		t.Error("figPos not cloned")
	}
	e1.figPos = Point{1,1}
	if e2.figPos.X == 1 {
		t.Error("figbox is a reference!")
	}
	
	if len(e1.boxes) != 1 || len(e2.boxes) != 1 {
		t.Error("boxes should contain one box")
	}
	e1.boxes[0] = &box2
	if e2.boxes[0].Pos.X != 0 {
		t.Error("boxes is a reference")
	}
}

func TestClonePoint(t *testing.T) {
	p1 := NewPoint(0,0)
	p2 := p1.Clone()
	p1.X = 1
	if p2.X == 1 {
		t.Error("ClonePoint: reference!")
	}
}