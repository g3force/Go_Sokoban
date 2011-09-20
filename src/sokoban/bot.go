package sokoban

import (
	"fmt"
	"syscall"
)

// indicates the direction, that the figure moves to 
// and counts the number of rotations
type DirType struct {
	counter int8
	dir     int8
}

type HistoryTree struct {
	//	p    uint16
	p    Point
	sons []*HistoryTree
}

var (
	path          = []DirType{}   // direction path (which directions the figure has taken yet)
	history       = HistoryTree{} // history with collection of all points and figure
	StraightAhead bool            // true: new direction are initialized with current dir, false: init with 0
	//	pointmap      = map[uint16]Point{}
	//	surWidth      int8
)

func Init() {
	path = []DirType{{-1, -1}}
	history = HistoryTree{NewPoint(-1, -1), nil}
	//	history = HistoryTree{0, nil}
	//	surWidth = int8(len(Surface[0]))
	//	for y, _ := range Surface {
	//		for x, _ := range Surface[y] {
	//			pointmap[uint16(y*int(surWidth)+x)] = NewPoint(x, y)
	//		}
	//	}
}

//func getPointMapId(x, y int8) uint16 {
//	return uint16(y*surWidth + x)
//}

// return only dir from path
func GetPath() []int8 {
	pa := []int8{}
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
	steps, solutions, solSteps := 0, 0, []int{}
	// init time counter
	var starttime syscall.Timeval
	syscall.Gettimeofday(&starttime)
	for {
		// ### 1. check if finished
		if len(path) == 0 {
			I("Empty path. Hopefully all possibilities tried ;)")
			break
		}
		// ### 2. increase the last path to try the next possibility
		incLastPath()
		// ### 3. check if rotation is finished. Then backtrack
		if getLastCounter() > 3 {
			I("Rotation finished. Deadlock. Backtrack.")
			UndoStep()
			rmLastPath()
			continue
		}
		// ### 4. Try moving
		I("Try moving in dir=%d", getLastPath())
		moved, boxMoved := Move(getLastPath())
		if !moved {
			I("Could not move.")
			continue
		}
		// ### 5. If moved, first check if not in a loop
		newHist := GetBoxesAndX()
		if everBeenHere(newHist) {
			I("I'v been here already. Backtrack.")
			UndoStep()
			continue
		}
		// ### 6. If not in a loop, append history and go on
		addHistory(newHist)
		if StraightAhead {
			addToPath(getLastPath() - 1)
		} else {
			addToPath(-1)
		}
		I("Moved. Path added.")
		// ### 7. Do some statistics
		steps++
		if printSurface {
			Print()
		}
		if steps%outputFreq == 0 {
			min, sec, µsec := getTimePassed(starttime)
			D("Steps: %9d; %4dm %2ds %6dµs", steps, min, sec, µsec)
		}
		// ### 8. Do we already won? :)
		if boxMoved != EMPTY && Won() {
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

func everBeenHere(boxes []Point) bool {
	h := &history
	for i := 0; i < len(boxes); i++ {
		box := boxes[i]
		son := searchSons(h, box)
		if son == -1 {
			return false
		} else {
			h = h.sons[son]
		}
	}
	if len(history.sons) == 0 {
		return false
	}
	return true
}

func addHistory(boxes []Point) {
	h := &history
	for i := 0; i < len(boxes); i++ {
		box := boxes[i]
		son := searchSons(h, box)
		if son == -1 {
			insertNewHist(h, boxes[i:], -1)
			break
		} else {
			h = h.sons[son]
		}
	}
}

func insertNewHist(h *HistoryTree, boxList []Point, counter int8) (newHis HistoryTree) {
	counter++
	if int8(len(boxList)) == counter {
		return
	}
	newHis = newHistoryTree(boxList[counter].X, boxList[counter].Y, nil)
	h.sons = append(h.sons, &newHis)
	insertNewHist(h.sons[len(h.sons)-1], boxList, counter)
	return
}

func newHistoryTree(x int8, y int8, sons []HistoryTree) HistoryTree {
	//	return HistoryTree{getPointMapId(x, y), nil}
	return HistoryTree{NewPoint8(x, y), nil}
}

func searchSons(h *HistoryTree, box Point) int {
	for key, value := range h.sons {
		//		if pointmap[value.p].X == box.X && pointmap[value.p].Y == box.Y {
		if value.p.X == box.X && value.p.Y == box.Y {
			return key
		}
	}
	return -1
}

func printTree() {
	fmt.Println(history)
}

// check, if a and b are equal
func sameFields(a []Point, b []Point) bool {
	//	if len(a) != len(b) {
	//		return false
	//	}
	for y := 0; y < len(a); y++ {
		if a[y].X != b[y].X || a[y].Y != b[y].Y {
			return false
		}
	}
	return true
}

func getLastPath() int8 {
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

func getLastCounter() int8 {
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

func addToPath(dir int8) {
	path = append(path, DirType{-1, dir})
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

