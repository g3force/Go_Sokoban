package sokoban

import (
	"fmt"
	"syscall"
)

type DirType struct {
	counter int
	dir int
}

var (
	path    = []DirType{}
	history = [][]Point{}
	StraightAhead bool
)

func Init() {
	path = []DirType{{-1,-1}}
	history = [][]Point{GetBoxesAndX()}
}

func GetPath() []int {
	pa := []int{}
	for i := 0; i < len(path); i++ {
		pa = append(pa, path[i].dir)
	}
	return pa
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
		if getLastCounter() > 3{
			I("Rotation finished. Deadlock. Backtrack.")
			UndoStep()
			rmLastPath()
		} else {
			I("Try moving in dir=%d", getLastPath())
			moved, boxMoved := Move(getLastPath())
			if moved {
				if (boxMoved && !deadEnd(addPoints(GetFigPos(), Direction(getLastPath())))) || !boxMoved {
					newHist := GetBoxesAndX()
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
						if StraightAhead {
							addToPath(getLastPath()-1)
						} else {
							addToPath(-1)
						}
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
	if path[len(path)-1].dir==4 {
		path[len(path)-1].dir = 0
	}
	path[len(path)-1].counter++
}

func getLastCounter() int{
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

