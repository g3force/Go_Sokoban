package ai

import (
	"sokoban/engine"
	"testing"
)

func TestCloneNode(t *testing.T) {
	var node1 Node
	var node2 Node

	node1 = NewNode()
	node1.SetDirection(2)
	node1.counter = 2
	node1.ignored = []engine.Direction{1, 3}

	node2 = node1.Clone()
	if node1.counter != node2.counter {
		t.Error("counter differs")
	}
	if node1.dir != node2.dir {
		t.Error("dir differs")
	}
	if len(node1.ignored) != 2 {
		t.Error("node1.ignored != 2")
	}
	if len(node1.ignored) != len(node2.ignored) {
		t.Errorf("array length differ: %d != %d", len(node1.ignored), len(node2.ignored))
	}

	node2.dir = 3
	if node1.dir == node2.dir {
		t.Error("dir is equal, should be different")
	}

	node2.PushIgnored(3)
	if len(node2.ignored) != 3 {
		t.Error("node2.ignored != 3")
	}
	if len(node1.ignored) == len(node2.ignored) {
		t.Error("ignored arrays must be different")
	}
}

func TestClonePath(t *testing.T) {
	var path1 Path
	var path2 Path
	
	path1 = Path{}
	path1.Push(-1)
	if path1.Empty() {
		t.Error("Path should not be empty")
	}
	path1.IncCurrentDir()
	if path1.CurrentDir() != 0 {
		t.Error("IncCurrentDir failed")
	}
	path2 = path1.Clone()
	if path1.CurrentDir() != path2.CurrentDir() {
		t.Error("Clone failed")
	}
	path1.IncCurrentDir()
	if path1.CurrentDir() == path2.CurrentDir() {
		t.Error("Clone failed (2)")
	}
}
