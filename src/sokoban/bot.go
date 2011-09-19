package sokoban

import (
	"fmt"
	"syscall"
)

// indicates the direction, that the figure moves to 
// and counts the number of rotations
type DirType struct {
	counter int
	dir     int
}

type HistoryTree struct {
	x int
	y int
	sons []HistoryTree
}

var (
	path          = []DirType{} // direction path (which directions the figure has taken yet)
	history       = HistoryTree{} // history with collection of all points and figure
	StraightAhead bool          // true: new direction are initialized with current dir, false: init with 0
)

func Init() {
	path = []DirType{{-1, -1}}
	history = HistoryTree{-1,-1, nil}
}

// return only dir from path
func GetPath() []int {
	pa := []int{}
	for i := 0; i < len(path); i++ {
		pa = append(pa, path[i].dir)
	}
	return pa
}

// run the algo by calling Step(), print some output and catch if won.
func Run(single bool, outputFreq int, printSurface bool) {
	Init()
	MarkDeadFields()
	fmt.Println()
	Print()
	j := 0
	steps := 0
	solutions := 0
	solSteps := []int{}
	var starttime syscall.Timeval
	syscall.Gettimeofday(&starttime)

	for {
		moved, finished := Step()
		// finished indicates, that all possibilities were tried
		if finished {
			break
		}
		if moved {
			steps++
			j = 0
		} else {
			j++
			if j > 1000 {
				E("Bot seems not to move")
				break
			}
		}
		if printSurface {
			Print()
		}
		if j == 0 && steps%outputFreq == 0 {
			min, sec, µsec := getTimePassed(starttime)
			D("Steps: %6d; %4dm %2ds %6dµs", steps, min, sec, µsec)
		}
		if Won() {
			solutions++
			min, sec, µsec := getTimePassed(starttime)
			fmt.Printf("%d. solution found after %d steps, %4dm %2ds %6dµs.\nPath: %d\n", solutions, steps, min, sec, µsec, GetPath())
			Print()
			solSteps = append(solSteps, steps)
			if single {
				break
			}
			UndoStep()
			rmLastPath()
		}
	}
	min, sec, µsec := getTimePassed(starttime)
	fmt.Printf("Run finished with %d steps after %dm %ds %dµs.\n%d solutions found at following steps:\n%d\n", steps, min, sec, µsec, solutions, solSteps)
}

func Step() (hasMoved bool, finished bool) {
	hasMoved = false
	finished = false
	if len(path) == 0 {
		I("Empty path. Hopefully all possibilities tried ;)")
		finished = true
		return
	} else {
		incLastPath() // first of all: increase the last path to try the next possibility
		if getLastCounter() > 3 {
			I("Rotation finished. Deadlock. Backtrack.")
			UndoStep()
			rmLastPath()
		} else {
			I("Try moving in dir=%d", getLastPath())
			moved, _ := Move(getLastPath())
			if moved {
				newHist := GetBoxesAndX() // TODO improve this?
				//				D("hist: %d", newHist)
				hit := false
				if everBeenHere(newHist) {
					I("I'v been here already. Backtrack.")
					UndoStep()
					hit = true
				}
				if !hit {
					addHistory(newHist)
					printTree()
					if StraightAhead {
						addToPath(getLastPath() - 1)
					} else {
						addToPath(-1)
					}
					I("Moved. Path added.")
					hasMoved = true
				}
			} else {
				I("Could not move.")
			}
		}
	}
	I("End Step. Path: %d", path)
	return
}

func everBeenHere(boxes []Point) bool {
	h := history
	for i := 0; i <  len(boxes); i++ {
		box := boxes[i]
		son := searchSons(h, box)
		if son != -1 {
			return false
		}
	}
	if len(history.sons)==0 {
		return false
	}
	return true
}

func addHistory(boxes []Point) {
	h := history
	for i := 0; i <  len(boxes); i++ {
		box := boxes[i]
		son := searchSons(h, box)
		if son == -1 {
			insertNewHist(h, boxes[i:])
			break
		} else {
			h=h.sons[son]
		}
	}
}

func insertNewHist(h HistoryTree, boxList []Point) {
	E("historytree: %d, boxes: %d", h, boxList)	
	for _, box := range boxList {
		newHis := newHistoryTree(box.X, box.Y, nil)
		h.sons = append(h.sons, newHis)
		h = newHis
	}
}

func newHistoryTree(x int, y int, sons []HistoryTree) HistoryTree {
	return HistoryTree{x, y, nil}
}

func searchSons(h HistoryTree, box Point) int {
	for key, value:=range h.sons {
		if value.x == box.X && value.y == box.Y {
			return key 	
		}
	}
	return -1
}

func printTree() {
	fmt.Print(history)
}

// check, if a and b are equal
func sameFields(a []Point, b []Point) bool {
	if len(a) != len(b) {
		return false
	}
	for y := 0; y < len(a); y++ {
		if a[y].X != b[y].X || a[y].Y != b[y].Y {
			return false
		}
	}
	return true
}


func getLastPath() int {
	if len(path) > 0 {
		return path[len(path)-1].dir
	}
	return -1
}

func incLastPath() {
	path[len(path)-1].dir++
	if path[len(path)-1].dir == 4 {
		path[len(path)-1].dir = 0
	}
	path[len(path)-1].counter++
}

func getLastCounter() int {
	if len(path) > 0 {
		return path[len(path)-1].counter
	}
	return -1
}

func rmLastPath() {
	if len(path) > 0 {
		path = path[:len(path)-1]
	} else {
		E("Path is empty. Can not remove from path.")
	}

}

func addToPath(dir int) {
	path = append(path, DirType{-1, dir})
}

func abs(a int) int {
	if a < 0 {
		a = a * -1
	}
	return a
}

// return min, sec and µsec since specified starttime
func getTimePassed(starttime syscall.Timeval) (min, sec, µsec int) {
	var time syscall.Timeval
	syscall.Gettimeofday(&time)
	sec = int(time.Sec - starttime.Sec)
	µsec = int(time.Usec - starttime.Usec)
	if µsec < 0 {
		sec--
		µsec = 1000000 + µsec
	}
	min = sec / 60
	sec = sec % 60
	return
}

