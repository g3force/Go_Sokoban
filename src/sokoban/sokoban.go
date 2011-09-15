package sokoban

import (
	"fmt"
	"os"
	"strings"
)

// constants for indicating what a field "contains"
const EMPTY = 0

// simple Point type
type Point struct {
	X int
	Y int
}

// History Type for saving a change after a move
type HistoryType struct {
	OldPos   Point
	NewPos   Point
	BoxMoved int
}

// single field within the surface
type Field struct {
	wall  bool
	point bool
	dead  bool
	box   int
}

type MovingResult struct {
	moved    bool
	boxMoved bool
	box      int
}

var (
	Surface [][]Field         // the current Surface
	History []HistoryType     // history, indicating the past way
	figPos  Point             // current position of figure
	points  []Point           // Array of all points
	boxes   = map[int]Point{} // Array of all boxes
)

func GetBoxes() map[int]Point {
	return boxes
}

// convert direction from int to Point
func Direction(dir int) Point {
	dir = dir % 4
	var p Point
	switch dir {
	case 0:
		p.X = 1
		p.Y = 0
	case 1:
		p.X = 0
		p.Y = 1
	case 2:
		p.X = -1
		p.Y = 0
	case 3:
		p.X = 0
		p.Y = -1
	}
	return p
}

/* try moving figure in specified direction.
 * Returns, if figure was moved and if figure moved a box.
 */
func Move(dir int) (success bool, boxMoved int) {
	success = false
	boxMoved = EMPTY
	cf := GetFigPos()                   // current figureposition
	nf := addPoints(cf, Direction(dir)) // potential new figureposition
	if !IsInSurface(nf) {
		I("Can not move: surface border")
		return
	}
	// check type of field of new figureposition
	switch f := Surface[nf.Y][nf.X]; {
	case f.wall == true:
		I("Can not move: wall")
		return
	case f.box != EMPTY: // if box
		nnf := addPoints(nf, Direction(dir)) // potential new boxposition
		if !IsInSurface(nnf) {
			I("Can not move: blocked box (surface border)")
			return
		}
		if Surface[nnf.Y][nnf.X].dead {
			I("Can not move: Dead field")
			return
		}
		if Surface[nnf.Y][nnf.X].wall || Surface[nnf.Y][nnf.X].box != EMPTY {
			I("Can not move: blocked box")
			return
		}
		I("Move box")
		boxMoved = f.box
		boxes[f.box] = nnf // nnf is position that box should move to, f.box is current box
		Surface[nnf.Y][nnf.X].box = Surface[nf.Y][nf.X].box
		fallthrough // go to next case statment to also move the figure
	case f.box >= EMPTY: // actually always...
		var hist HistoryType
		hist.NewPos = nf
		hist.OldPos = cf
		hist.BoxMoved = boxMoved
		History = append(History, hist)
		Surface[nf.Y][nf.X].box = Surface[cf.Y][cf.X].box
		Surface[cf.Y][cf.X].box = EMPTY
		figPos = nf // refresh figureposition
		success = true
	default:
		E("Unknown field")
	}
	return
}

// undo the last step (move figure and box to their old positions)
func UndoStep() {
	if len(History) > 0 {
		history := History[len(History)-1] // get last history
		Surface[history.OldPos.Y][history.OldPos.X].box = Surface[history.NewPos.Y][history.NewPos.X].box
		Surface[history.NewPos.Y][history.NewPos.X].box = EMPTY
		figPos = history.OldPos
		// also move box back, if neccessary
		if history.BoxMoved != EMPTY {
			var boxPoint Point
			boxPoint.X = history.NewPos.X + (history.NewPos.X - history.OldPos.X)
			boxPoint.Y = history.NewPos.Y + (history.NewPos.Y - history.OldPos.Y)
			Surface[history.NewPos.Y][history.NewPos.X].box = Surface[boxPoint.Y][boxPoint.X].box
			Surface[boxPoint.Y][boxPoint.X].box = EMPTY
			boxes[Surface[history.NewPos.Y][history.NewPos.X].box] = history.NewPos
		}
		History = History[:len(History)-1] // remove from history
	}
}

func GetFigPos() Point {
	return figPos
}

// load level from specified file (relative to binary file)
func LoadLevel(filename string) {
	raw, _ := contents(filename)
	// remove the "\r" from stupid windows files...
	raw = strings.Replace(raw, "\r", "", -1)
	// get single lines in an array
	lines := strings.Split(raw, "\n")

	Surface = [][]Field{{}}
	var field Field
	y := 0
	boxId := 0
	for _, line := range lines {
		// filter empty lines and lines that do not start with '#'
		if len(line) == 0 || line[0] != '#' {
			continue
		}
		for x, char := range line {
			switch char {
			case '#':
				field = Field{true, false, false, EMPTY}
			case ' ':
				field = Field{false, false, false, EMPTY}
			case '$':
				boxId++
				field = Field{false, false, false, boxId}
			case '@':
				field = Field{false, false, false, EMPTY}
				figPos = Point{x, y}
			case '.':
				field = Field{false, true, false, EMPTY}
			case '*':
				boxId++
				field = Field{false, true, false, boxId}
			case '+':
				field = Field{false, true, false, EMPTY}
				figPos = Point{x, y}
			default:
				E("Unknown character in level file: '%c'", char)
			}
			Surface[y] = append(Surface[y], field)
			if field.point {
				points = append(points, Point{x, y})
			}
			if field.box != EMPTY {
				boxes[boxId] = Point{x, y}
			}
		}
		y++
		Surface = append(Surface, []Field{})
	}
	// the last sub-array of Surface is always empty, so remove it...
	if len(Surface[len(Surface)-1]) == 0 {
		Surface = Surface[:len(Surface)-1]
	}
	return
}

// loop over all points and check, if there is a box. Else return false
func Won() bool {
	for _, p := range points {
		if Surface[p.Y][p.X].box == EMPTY {
			return false
		}
	}
	return true
}

// print the current Surface
func Print() {
	for y := 0; y < len(Surface); y++ {
		for x := 0; x < len(Surface[y]); x++ {
			switch field := Surface[y][x]; {
			case field.wall:
				fmt.Print("#")
			case figPos.X == x && figPos.Y == y:
				if field.point {
					fmt.Print("+")
				} else {
					fmt.Print("x")
				}
			case field.box == EMPTY:
				if field.point {
					fmt.Print("*")
				} else if field.dead {
					fmt.Print("☠")
				} else {
					fmt.Print(" ")
				}
			default: // field has box
				if field.point {
					fmt.Print("%")
				} else {
					fmt.Print("$")
				}
			}
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

// return Point array of all boxes and the figure
func GetBoxesAndX() (field []Point) {
	//	field = boxes
	//	field[0] = figPos
	field = append(field, figPos)

	var sort = map[int]map[int]int{}

	for _, box := range boxes {
		if sort[box.Y] == nil {
			sort[box.Y] = map[int]int{}
		}
		sort[box.Y][box.X] = 0
	}

	for y, _ := range sort {
		for x, _ := range sort[y] {
			field = append(field, Point{x, y})
		}
	}
	//	D("field: %d", field)

	//	for i := 1; i < len(boxes)+1; i++ {
	//		field = append(field, boxes[i])
	//	}
	//	temp := 0
	//	for y := 0; y < len(Surface); y++ {
	//		for x := 0; x < len(Surface[y]); x++ {
	//			if Surface[y][x].box != EMPTY {
	//				field = append(field, Point{x, y})
	//				//				temp++
	//			}
	//		}
	//	}
	return
}

// print a legend of the Surface output
func PrintInfo() {
	fmt.Println("Surface Field association:")
	fmt.Printf("EMPTY\t\t' '\n")
	fmt.Printf("BOX\t\t'$'\n")
	fmt.Printf("FIGURE\t\t'x'\n")
	fmt.Printf("EMPTY POINT\t'*'\n")
	fmt.Printf("BOX POINT\t'%'\n")
	fmt.Printf("FIGURE POINT\t'+'\n")
	fmt.Printf("WALL\t\t'#'\n")
	fmt.Printf("DEAD FIELD\t'☠'\n")
}

// return number of boxes on the surface
func CountBoxes() int {
	count := 0
	for y := 0; y < len(Surface); y++ {
		for x := 0; x < len(Surface[y]); x++ {
			if Surface[y][x].box != EMPTY {
				count++
			}
		}
	}
	return count
}

// check if the surface border was reached
func IsInSurface(p Point) bool {
	if p.Y < 0 || p.X < 0 || p.Y >= len(Surface) || p.X >= len(Surface[0]) {
		D("not in surface: %d", p)
		return false
	}
	return true
}

// add to points (their x and y)
func addPoints(p1 Point, p2 Point) Point {
	var np Point
	np.X = p1.X + p2.X
	np.Y = p1.Y + p2.Y
	return np
}

// read from the specified file and return whole content in a string
func contents(filename string) (string, os.Error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var result []byte
	buf := make([]byte, 100)
	for {
		n, err := f.Read(buf[0:])
		result = append(result, buf[0:n]...)
		if err != nil {
			if err == os.EOF {
				break
			}
			return "", err
		}
	}
	return string(result), nil
}

