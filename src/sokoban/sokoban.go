package sokoban

import (
	"fmt"
	"os"
	"strings"
)

// constants for displaying and saving the surface
const GROUND = 8
const WALL = 9
const POINT = 10

// constants for indicating what a field "contains"
const EMPTY = 0
const BOX = 1
const FIGURE = 2

// simple Point type
type Point struct {
	X int
	Y int
}

// History Type for saving a change after a move
type HistoryType struct {
	OldPos   Point
	NewPos   Point
	BoxMoved bool
}

// single field within the surface
type Field struct {
	wall    bool
	point   bool
	dead    bool
	contain int
}

var (
	Surface [][]Field     // the current Surface
	History []HistoryType // history, indicating the past way
	figPos  Point         // current position of figure
	points  []Point       // Array of all points
)

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
func Move(dir int) (success bool, boxMoved bool) {
	success = false
	boxMoved = false
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
	case f.contain == BOX:
		nnf := addPoints(nf, Direction(dir)) // potential new boxposition
		if !IsInSurface(nnf) {
			I("Can not move: blocked box (surface border)")
			return
		}
		if Surface[nnf.Y][nnf.X].dead{
			I("Can not move: Dead field")
			return
		}
		if Surface[nnf.Y][nnf.X].wall || Surface[nnf.Y][nnf.X].contain != EMPTY {
			I("Can not move: blocked box")
			return
		}
		I("Move box")
		boxMoved = true
		Surface[nnf.Y][nnf.X].contain = Surface[nf.Y][nf.X].contain
		fallthrough // go to next case statment to also move the figure
	case f.contain == EMPTY || f.contain == BOX:
		var hist HistoryType
		hist.NewPos = nf
		hist.OldPos = cf
		hist.BoxMoved = boxMoved
		History = append(History, hist)
		Surface[nf.Y][nf.X].contain = Surface[cf.Y][cf.X].contain
		Surface[cf.Y][cf.X].contain = EMPTY
		figPos = nf // refresh figureposition
		success = true
	case f.contain == FIGURE:
		E("Duplicate figures or bad direction")
	default:
		E("Unknown field")
	}
	return
}

// undo the last step (move figure and box to their old positions)
func UndoStep() {
	if len(History) > 0 {
		history := History[len(History)-1] // get last history
		Surface[history.OldPos.Y][history.OldPos.X].contain = Surface[history.NewPos.Y][history.NewPos.X].contain
		Surface[history.NewPos.Y][history.NewPos.X].contain = EMPTY
		figPos = history.OldPos
		// also move box back, if neccessary
		if history.BoxMoved {
			var boxPoint Point
			boxPoint.X = history.NewPos.X + (history.NewPos.X - history.OldPos.X)
			boxPoint.Y = history.NewPos.Y + (history.NewPos.Y - history.OldPos.Y)
			Surface[history.NewPos.Y][history.NewPos.X].contain = Surface[boxPoint.Y][boxPoint.X].contain
			Surface[boxPoint.Y][boxPoint.X].contain = EMPTY
		}
		History = History[:len(History)-1] // remove from history
	}
}

func GetFigPos() Point {
	return figPos
}

//func Find(object int) Point {
//	return FindNext(object, 0, 0)
//}
//
//func FindNext(object int, startx int, starty int) (p Point) {
//	for y := starty; y < len(Surface); y++ {
//		for x := startx; x < len(Surface[y]); x++ {
//			if Surface[y][x].contain == object {
//				p.X = x
//				p.Y = y
//				return
//			}
//		}
//	}
//	return
//}

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
				field = Field{false, false, false, BOX}
			case '@':
				field = Field{false, false, false, FIGURE}
				figPos = Point{x, y}
			case '.':
				field = Field{false, true, false, EMPTY}
			case '*':
				field = Field{false, true, false, BOX}
			case '+':
				field = Field{false, true, false, FIGURE}
			default:
				E("Unknown character in level file: '%c'", char)
			}
			Surface[y] = append(Surface[y], field)
			if field.point {
				points = append(points, Point{x, y})
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
		if Surface[p.Y][p.X].contain != BOX {
			return false
		}
	}
	return true
}

// print the current Surface
func Print() {
	for y := 0; y < len(Surface); y++ {
		for x := 0; x < len(Surface[y]); x++ {
			if Surface[y][x].wall {
				fmt.Print("#")
			} else if Surface[y][x].dead {
				fmt.Print("☠")
			} else if Surface[y][x].point {
				switch Surface[y][x].contain {
				case EMPTY:
					fmt.Print("*")
				case BOX:
					fmt.Print("%")
				case FIGURE:
					fmt.Print("+")
				}
			} else {
				switch Surface[y][x].contain {
				case EMPTY:
					fmt.Print(" ")
				case BOX:
					fmt.Print("$")
				case FIGURE:
					fmt.Print("x")
				}
			}
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

// return Point array of all boxes and the figure
func GetBoxesAndX() (field []Point) {
	field = append(field, figPos)
	//	temp := 0
	for y := 0; y < len(Surface); y++ {
		for x := 0; x < len(Surface[y]); x++ {
			if Surface[y][x].contain == BOX {
				field = append(field, Point{x, y})
				//				temp++
			}
		}
	}
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
			if Surface[y][x].contain == BOX {
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

