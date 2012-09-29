package ai

import (
	"fmt"
	"syscall"
//	"unsafe"
	"sokoban/engine"
	"sokoban/log"
)

type HistoryTree struct {
	p    engine.Point
	sons []*HistoryTree
}

//var (
	//	pointmap      = map[uint16]engine.Point{}
	//	surWidth      int8
//)

//func Init() (history *HistoryTree, path *[]DirType) {
	//	history = HistoryTree{0, nil}
	//	surWidth = int8(len(Surface[0]))
	//	for y, _ := range Surface {
	//		for x, _ := range Surface[y] {
	//			pointmap[uint16(y*int(surWidth)+x)] = NewPoint(x, y)
	//		}
	//	}
//}

// run the algo, print some output and catch if won.
// straightAhead: true: new direction are initialized with current dir, false: init with 0
func Run(e engine.Engine, single bool, outputFreq int, printSurface bool, straightAhead bool) {
	path := Path{}
	path.Push(engine.NO_DIRECTION)
	history := HistoryTree{engine.NewPoint(-1, -1), nil}
	
	MarkDeadFields(&e.Surface)
	fmt.Println()
	e.Print()
	steps, solutions, solSteps := 0, 0, []int{}
	//println(unsafe.Sizeof(HistoryTree{}))
	var ignoredDir = false
	// init time counter
	var starttime syscall.Timeval
	syscall.Gettimeofday(&starttime)



//}
//
//func runWorker(single bool, outputFreq int, printSurface bool, straightAhead bool) {

	for {
		ignoredDir = false
		// ### 1. check if finished
		if path.Empty() {
			log.D("Empty path. Hopefully all possibilities tried ;)")
			break
		}
		// ### 2. increase the last path to try the next possibility
		path.IncCurrentDir()
		// ### 3. check if rotation is finished. Then backtrack
		if path.CurrentCounter() > 3 {
			var dir = path.Current().PopIgnored()
			if dir == -1 {
				log.D("Rotation finished. Deadlock. Backtrack.")
				e.UndoStep()
				path.Pop()
				continue
			} else {
				path.SetCurrentDir(dir)
				ignoredDir = true
			}
		}
		// ### 4a. check if there is a box in direction dir and if this box is on a point
		cf := e.FigPos()                             // current figureposition
		nf := cf.Add(path.CurrentDir().Point()) // potential new figureposition
		if e.Surface[nf.Y][nf.X].Point && e.Surface[nf.Y][nf.X].Box != 0 && !ignoredDir {
			log.D("Do not moving a box from a point")
			path.Current().PushIgnored(path.CurrentDir())
			continue
		}
		// ### 4b. Try moving
		log.D("Try moving in dir=%d", path.CurrentDir())
		moved, boxMoved := e.Move(path.CurrentDir())
		if !moved {
			log.D("Could not move.")
			continue
		}
		// ### 5. If moved, first check if not in a loop
		newHist := e.GetBoxesAndX()
		if everBeenHere(&history, newHist) {
			log.D("I'v been here already. Backtrack.")
			e.UndoStep()
			continue
		}
		// ### 6. If not in a loop, append history and go on
		addHistory(&history, newHist)
		if straightAhead {
			path.Push(path.CurrentDir() - 1)
		} else {
			path.Push(-1)
		}
		log.D("Moved. Path added.")
		// ### 7. Do some statistics
		steps++
		if printSurface {
			e.Print()
		}
		if steps%outputFreq == 0 {
			min, sec, µsec := getTimePassed(starttime)
			log.I("Steps: %9d; %4dm %2ds %6dµs", steps, min, sec, µsec)
		}
		// ### 8. Do we already won? :)
		if boxMoved != 0 && e.Won() {
			solutions++
			min, sec, µsec := getTimePassed(starttime)
			fmt.Printf("%d. solution found after %d steps, %4dm %2ds %6dµs.\nPath: %d\n", solutions, steps, min, sec, µsec, path.Directions())
			e.Print()
			solSteps = append(solSteps, steps)
			if single {
				break
			}
			e.UndoStep()
			path.Pop()
		}
	}

	min, sec, µsec := getTimePassed(starttime)
	fmt.Printf("Run finished with %d steps after %dm %ds %dµs.\n%d solutions found at following steps:\n%d\n", steps, min, sec, µsec, solutions, solSteps)
}

func everBeenHere(history *HistoryTree, boxes []engine.Point) bool {
	h := history
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

func addHistory(history *HistoryTree, boxes []engine.Point) {
	h := history
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

func insertNewHist(h *HistoryTree, boxList []engine.Point, counter int8) (newHis HistoryTree) {
	counter++
	if int8(len(boxList)) == counter {
		return
	}
	newHis = newHistoryTree(boxList[counter].X, boxList[counter].Y)
	h.sons = append(h.sons, &newHis)
	insertNewHist(h.sons[len(h.sons)-1], boxList, counter)
	return
}

func newHistoryTree(x int8, y int8) HistoryTree {
	return HistoryTree{engine.NewPoint8(x, y), nil}
}

func searchSons(h *HistoryTree, box engine.Point) int {
	for key, value := range h.sons {
		if value.p.X == box.X && value.p.Y == box.Y {
			return key
		}
	}
	return -1
}

func printTree(history HistoryTree) {
	fmt.Println(history)
}

// check, if a and b are equal
func sameFields(a []engine.Point, b []engine.Point) bool {
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

