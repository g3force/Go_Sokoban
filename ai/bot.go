package ai

import (
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	//	"unsafe"
	"github.com/g3force/Go_Sokoban/engine"
	"github.com/g3force/Go_Sokoban/log"

//	"time"
)

type HistoryTree struct {
	p    engine.Point
	sons []*HistoryTree
}

var (
	wg         sync.WaitGroup
	cDone      chan int8
	history    HistoryTree
	cHistory   chan bool
	steps      int32
	solutions  int32
	solSteps   []int32
	starttime  syscall.Timeval
	numWorkers int
	running    bool
)

func incSteps() {
	atomic.AddInt32(&steps, 1)
}

func incSolutions() {
	atomic.AddInt32(&solutions, 1)
}

// run the algo, print some output and catch if won.
// straightAhead: true: new direction are initialized with current dir, false: init with 0
func Run(e engine.Engine, single bool, outputFreq int32, printSurface bool, straightAhead bool, threads int) {
	steps, solutions, solSteps = 0, 0, []int32{}
	numWorkers = 0
	running = true
	cDone = make(chan int8, threads) // queue for threads
	cHistory = make(chan bool, 1)    // mutex on global history object

	// preprocessing
	MarkDeadFields(&e.Surface)
	e.Print()

	// init time counter
	syscall.Gettimeofday(&starttime)

	// ???
//	runtime.GOMAXPROCS(runtime.NumCPU())

	// init path
	path := Path{}
	path.Push(engine.NO_DIRECTION)

	// create history store and save initial constellation
	history = newHistoryTree(-1, -1)
	newHist := e.GetBoxesAndX()
	addHistory(&history, newHist)

	// prepare for starting workers
	wg.Add(1)
	cDone <- 1
	go runWorker(e, path, threads, single, outputFreq, printSurface, straightAhead)

	// wait for all workers to finish
	wg.Wait()
	//	time.Sleep(1 * time.Second)

	// print result
	min, sec, µsec := getTimePassed(starttime)
	log.A("Run finished with %d steps after %dm %ds %dµs.\n%d solutions found at following steps:\n%d\n", steps, min, sec, µsec, solutions, solSteps)
}

func runWorker(e engine.Engine, basePath Path, threads int, single bool, outputFreq int32, printSurface bool, straightAhead bool) {
	// make sure, process will wait for this worker
	defer wg.Done()
	gorNo := numWorkers
	numWorkers++
	log.I(gorNo, "runWorker %d created, %d running", gorNo, runtime.NumGoroutine())
	e.Id = gorNo
	//	path := Path{basePath[len(basePath)-1].Clone()}
	path := basePath[len(basePath)-1:]

	basePath = basePath[:len(basePath)-1]
	//	basePath = basePath.Clone()

	//	log.I(gorNo, "path: %d, basePath: %d", path, basePath)

	for {
		var ignoredDir = false
		ignoredDir = false
		// ### 1. check if finished
		if path.Empty() || !running {
			log.D(e.Id, "Empty path / stopped executition. Hopefully all possibilities tried ;)")
			break
		}
		// ### 2. increase the current direction to try the next possibility
		path.IncCurrentDir()
		// ### 3. check if rotation is finished. Then backtrack
		if path.CurrentCounter() > 3 {
			var dir = path.Current().PopIgnored()
			if dir == -1 {
				log.D(e.Id, "Rotation finished. Deadlock. Backtrack.")
				e.UndoStep()
				path.Pop()
				continue
			} else {
				path.SetCurrentDir(dir)
				ignoredDir = true
			}
		}
		// ### 4a. check if there is a box in direction dir and if this box is on a point
		cf := e.FigPos()                        // current figureposition
		nf := cf.Add(path.CurrentDir().Point()) // potential new figureposition
		if e.Surface[nf.Y][nf.X].Point && e.Surface[nf.Y][nf.X].Box != 0 && !ignoredDir {
			log.D(e.Id, "Do not moving a box from a point")
			path.Current().PushIgnored(path.CurrentDir())
			continue
		}
		// ### 4b. Try moving
		log.D(e.Id, "Try moving in dir=%d", path.CurrentDir())
		moved, boxMoved := e.Move(path.CurrentDir())
		if !moved {
			log.D(e.Id, "Could not move.")
			continue
		}
		// ### 5. If moved, first check if not in a loop
		newHist := e.GetBoxesAndX()
		if everBeenHere(&history, newHist) {
			log.D(e.Id, "I'v been here already. Backtrack: %d", newHist)
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
		log.D(e.Id, "Moved. Path added.")
		// ### 7. Do some statistics
		incSteps()
		if printSurface {
			e.Print()
		}
		if steps%outputFreq == 0 {
			min, sec, µsec := getTimePassed(starttime)
			log.I(gorNo, "Steps: %9d; %4dm %2ds %6dµs", steps, min, sec, µsec)
		}
		// ### 8. Do we already won? :)
		if boxMoved != 0 && e.Won() {
			incSolutions()
			min, sec, µsec := getTimePassed(starttime)
			stepsCpy := steps
			log.Lock <- 1
			log.A("%d. solution found after %d steps, %4dm %2ds %6dµs.\nPath: %d%d\n", solutions, stepsCpy, min, sec, µsec, basePath.Directions(), path.Directions())
			<-log.Lock
			e.Print()
			solSteps = append(solSteps, stepsCpy)
			if single {
				running = false
				break
			}
			e.UndoStep()
			path.Pop()
		}
		// ### 9. Decide if we just go on or if we start a new thread
		if len(cDone) < threads {
			wg.Add(1)
			cDone <- 1
			ne := e.Clone()
			log.I(gorNo, "Creating new worker")
			go runWorker(ne, path.Clone(), threads, single, outputFreq, printSurface, straightAhead)
			// go back, as we deligated current dir go worker
			//			time.Sleep(3 * time.Second)
			e.UndoStep()
			path.Pop()
		}

	}
	<-cDone
	log.I(gorNo, "runWorker %d finished", gorNo)
}

func everBeenHere(h *HistoryTree, boxes []engine.Point) bool {
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

func addHistory(h *HistoryTree, boxes []engine.Point) {
	for i := 0; i < len(boxes); i++ {
		box := boxes[i]
		son := searchSons(h, box)
		if son == -1 {
			nBoxes := make([]engine.Point, len(boxes)-i)
			for k := i; k+i < len(boxes); k++ {
				nBoxes[k] = boxes[i+k].Clone()
			}
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
	cHistory <- true

	newHis = newHistoryTree(boxList[counter].X, boxList[counter].Y)
	h.sons = append(h.sons, &newHis)
	<-cHistory // release slot
	insertNewHist(h.sons[len(h.sons)-1], boxList, counter)

	return
}

func newHistoryTree(x int8, y int8) HistoryTree {
	return HistoryTree{engine.NewPoint8(x, y), nil}
}

func searchSons(h *HistoryTree, box engine.Point) int {
	cHistory <- true

	for key, value := range h.sons {
		if value.p.X == box.X && value.p.Y == box.Y {
			<-cHistory // release slot
			return key
		}
	}

	<-cHistory // release slot
	return -1
}

func printTree(history HistoryTree) {
	log.A("%d\n", history)
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
