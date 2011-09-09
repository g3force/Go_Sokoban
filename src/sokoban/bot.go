package sokoban

import (
	"fmt"
	"syscall"
)

var (
	path    = []int{-1}
	history = [][][]int8{}
)

func Init() {
	path = []int{-1}
	history = [][][]int8{GetSimpleSurface()}
}

func GetPath() []int {
	return path
}

func Run(single bool) {
	Init()
	j := 0
	steps := 0
	solutions := 0
	solSteps := []int{}
	var starttime syscall.Timeval
	syscall.Gettimeofday(&starttime)

	for {
		moved, finished := Step()
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
		if j == 0 && steps%5000 == 0 {
			min, sec, µsec := getTimePassed(starttime)
			D("Steps: %d; %dmin, %dsec, %dµsec", steps, min, sec, µsec)
		}
		if Won() {
			solutions++
			min, sec, µsec := getTimePassed(starttime)
			fmt.Printf("%d. solution found after %d steps, %4dm %2ds %6dµs.\nPath: %d\n", solutions, steps, min, sec, µsec, path)
			Print()
			UndoStep()
			rmLastPath()
			solSteps = append(solSteps, steps)
			if single {
				break
			}
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
		incLastPath()
		if getLastPath() > 3 {
			I("Rotation finished. Deadlock. Backtrack.")
			UndoStep()
			rmLastPath()
		} else {
			I("Try moving in dir=%d", getLastPath())
			moved, boxMoved := Move(getLastPath())
			if moved {
				if (boxMoved && !deadEnd(addPoints(GetFigPos(), Direction(getLastPath())))) || !boxMoved {
					newHist := GetSimpleSurface()
					hit := false
					for i := 0; i < len(history); i++ {
						if sameFields(history[i], newHist) {
							I("I'v been here already. Backtrack.")
							UndoStep()
							hit = true
							break
						}
					}
					if !hit {
						history = append(history, newHist)
						addToPath(-1)
						I("Moved. Path added.")
						hasMoved = true
					}
				} else {
					I("Dead End")
					UndoStep()
				}
			} else {
				I("Could not move.")
			}
		}
	}
	I("End Step. Path: %d", path)
	return
}

func deadEnd(box Point) bool {
	var p Point
	hit := false
	x := 0

	if Surface[box.Y][box.X].point {
		return false
	}

	for i := 0; i < 5; i++ {
		x = i % 4
		p = addPoints(box, Direction(x))
		//		D("%t, p=%d, box=%d", !IsInSurface(p), p, box)
		if !IsInSurface(p) || Surface[p.Y][p.X].wall {
			if hit {
				return true
			} else {
				hit = true
			}
		} else {
			hit = false
		}
	}

	return false
}

func sameFields(a [][]int8, b [][]int8) bool {
	if len(a) != len(b) {
		return false
	}
	for y := 0; y < len(a); y++ {
		if len(a[y]) != len(b[y]) {
			return false
		}
		for x := 0; x < len(a[y]); x++ {
			if a[y][x] != b[y][x] {
				return false
			}
		}
	}
	return true
}

func getLastPath() int {
	if len(path) > 0 {
		return path[len(path)-1]
	}
	return -1
}

func incLastPath() {
	path[len(path)-1]++
}

func rmLastPath() {
	if len(path) > 0 {
		path = path[:len(path)-1]
	} else {
		E("Path is empty. Can not remove from path.")
	}

}

func addToPath(dir int) {
	path = append(path, dir)
}

func abs(a int) int {
	if a < 0 {
		a = a * -1
	}
	return a
}

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

