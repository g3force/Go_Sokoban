package engine

import (
	"io"
	"os"
	"sokoban/log"
	"strings"
)

// constants for indicating what a field "contains"
const EMPTY = 0

// History Type for saving a change after a move
type HistoryType struct {
	OldPos   Point
	NewPos   Point
	BoxMoved int8
}

// single field within the surface
type Field struct {
	Wall  bool
	Point bool
	Dead  bool
	Box   int8
}

type MovingResult struct {
	moved    bool
	boxMoved bool
	box      int8
}

type Surface [][]Field

type Engine struct {
	Surface      Surface       // the current Surface
	History      []HistoryType // history, indicating the past way
	figPos       Point         // current position of figure
	points       []Point       // Array of all points
	boxes        map[int8]*Box // Array of all boxes
	boxesOrdered map[int8]*Box
	Id 				int
}

func NewEngine() (e Engine) {
	e.boxes = map[int8]*Box{}
	e.boxesOrdered = map[int8]*Box{}
	return
}

func (e *Engine) Clone() (ne Engine) {
	ne = NewEngine()
	ne.Surface = e.Surface.Clone()
	ne.History = []HistoryType{} // empty
	ne.figPos = Point{e.FigPos().X, e.FigPos().Y}
	ne.points = e.points // won't change
	for k, v := range e.boxes {
		box := v.Clone()
    	ne.boxes[k] = &box
    	ne.boxesOrdered[box.Order] = &box
	}
	return
}

func (e Engine) Points() []Point {
	return e.points
}

func (e Engine) Boxes() map[int8]*Box {
	return e.boxes
}

func (e Engine) FigPos() Point {
	return e.figPos
}

func (surface Surface) AmountOfFields() (fields int, dead int8) {
	for y := 0; y < len(surface); y++ {
		for x := 0; x < len(surface[y]); x++ {
			if !surface[y][x].Wall {
				fields++
			}
			if surface[y][x].Dead {
				dead++
			}
		}
	}
	return
}

func (surface *Surface) Clone() (ns Surface) {
	// create new array (rows)
	ns = make([][]Field, len(*surface))
	// add rows
	for y:=0 ; y < len(*surface) ; y++ {
		row := make([]Field, len((*surface)[y]))
		copy(row, (*surface)[y])
		ns[y] = row
	}
	return
}

/* try moving figure in specified direction.
 * Returns, if figure was moved and if figure moved a box.
 */
func (e *Engine) Move(dir Direction) (success bool, boxMoved int8) {
	success = false
	boxMoved = EMPTY
	cf := e.FigPos()          // current figureposition
	nf := cf.Add(dir.Point()) // potential new figureposition
	if !e.Surface.In(nf) {
		log.D(e.Id,"Can not move: surface border")
		return
	}
	// check type of field of new figureposition
	switch f := e.Surface[nf.Y][nf.X]; {
	case f.Wall == true:
		log.D(e.Id,"Can not move: wall")
		return
	case f.Box != EMPTY: // if box
		nnf := nf.Add(dir.Point()) // potential new boxposition
		if !e.Surface.In(nnf) {
			log.D(e.Id,"Can not move: blocked box (surface border)")
			return
		}
		if e.Surface[nnf.Y][nnf.X].Dead {
			log.D(e.Id,"Can not move: Dead field")
			return
		}
		if e.Surface[nnf.Y][nnf.X].Wall || e.Surface[nnf.Y][nnf.X].Box != EMPTY {
			log.D(e.Id,"Can not move: blocked box")
			return
		}
		log.D(e.Id,"Move box")
		boxMoved = f.Box
		e.boxes[f.Box].SetPos(nnf) // nnf is position that box should move to, f.Box is current box
		e.reOrderBoxes(f.Box, dir)
		e.Surface[nnf.Y][nnf.X].Box = e.Surface[nf.Y][nf.X].Box
		fallthrough // go to next case statment to also move the figure
	case f.Box >= EMPTY: // actually always...
		var hist HistoryType
		hist.NewPos = nf
		hist.OldPos = cf
		hist.BoxMoved = boxMoved
		e.History = append(e.History, hist)
		e.Surface[nf.Y][nf.X].Box = e.Surface[cf.Y][cf.X].Box
		e.Surface[cf.Y][cf.X].Box = EMPTY
		e.figPos = nf // refresh figureposition
		success = true
	default:
		log.E(e.Id,"Unknown field")
	}
	return
}

func (e *Engine) reOrderBoxes(curBoxId int8, dir Direction) {
	switch dir {
	case 1: // down
	case -1: // up
	case 3: // up
		dir = -1 // convert
	default:
		// nothing todo, as there is no up/down movement
		return
	}
	curBox := e.boxes[curBoxId]
	curOrder := curBox.Order
	if curOrder+dir.Int() <= 0 || curOrder+dir.Int() > int8(len(e.boxes)) {
		return
	}
	nextBox := e.boxesOrdered[curOrder+dir.Int()]

	//	D("1. curBox:%d nextBox:%d dir=%d", *curBox, *nextBox, dir)

	if dir == -1 {
		if nextBox.Pos.Y == curBox.Pos.Y {
			if nextBox.Pos.X <= curBox.Pos.X {
				return
			}
		} else if nextBox.Pos.Y <= curBox.Pos.Y {
			return
		}
	} else {
		if nextBox.Pos.Y == curBox.Pos.Y {
			if nextBox.Pos.X >= curBox.Pos.X {
				return
			}
		} else if nextBox.Pos.Y >= curBox.Pos.Y {
			return
		}
	}
	nextBox.SetOrder(curOrder)
	curBox.SetOrder(curOrder + dir.Int())
	e.boxesOrdered[curOrder+dir.Int()] = curBox
	e.boxesOrdered[curOrder] = nextBox
	e.reOrderBoxes(curBoxId, dir)
}

// undo the last step (move figure and box to their old positions)
func (e *Engine) UndoStep() {
	if len(e.History) > 0 {
		history := e.History[len(e.History)-1] // get last history
		e.Surface[history.OldPos.Y][history.OldPos.X].Box = e.Surface[history.NewPos.Y][history.NewPos.X].Box
		e.Surface[history.NewPos.Y][history.NewPos.X].Box = EMPTY
		e.figPos = history.OldPos
		// also move box back, if neccessary
		if history.BoxMoved != EMPTY {
			var boxPoint Point
			boxPoint.X = history.NewPos.X + (history.NewPos.X - history.OldPos.X)
			boxPoint.Y = history.NewPos.Y + (history.NewPos.Y - history.OldPos.Y)
			e.Surface[history.NewPos.Y][history.NewPos.X].Box = e.Surface[boxPoint.Y][boxPoint.X].Box
			e.Surface[boxPoint.Y][boxPoint.X].Box = EMPTY
			e.boxes[e.Surface[history.NewPos.Y][history.NewPos.X].Box].SetPos(history.NewPos)
			// if movement was up or down
			if history.NewPos.X-history.OldPos.X == 0 {
				e.reOrderBoxes(history.BoxMoved, (Direction)((int8)(history.OldPos.Y-history.NewPos.Y)))
			}
		}
		e.History = e.History[:len(e.History)-1] // remove from history
	}
}

// load level from specified file (relative to binary file)
func (e *Engine) LoadLevel(filename string) {
	raw, err := readLevelAsString(filename)
	if err != nil {
		panic(err)
	}
	// remove the "\r" from stupid windows files...
	raw = strings.Replace(raw, "\r", "", -1)
	// get single lines in an array
	lines := strings.Split(raw, "\n")

	e.Surface = Surface{{}}
	var field Field
	y := 0
	boxId := int8(0)
	maxlen := 0
	var char uint8

	for _, line := range lines {
		if len(line) > 0 && line[0] == '#' && len(line) > maxlen {
			maxlen = len(line)
		}
	}

	for _, line := range lines {
		// filter empty lines and lines that do not start with '#'
		if len(line) == 0 || line[0] != '#' {
			continue
		}
		for x := 0; x < maxlen; x++ {
			char = '#'
			if x < len(line) {
				char = line[x]
			}
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
				e.figPos = NewPoint(x, y)
			case '.':
				field = Field{false, true, false, EMPTY}
			case '*':
				boxId++
				field = Field{false, true, false, boxId}
			case '+':
				field = Field{false, true, false, EMPTY}
				e.figPos = NewPoint(x, y)
			default:
				log.E(e.Id,"Unknown character in level file: '%c'", char)
			}
			e.Surface[y] = append(e.Surface[y], field)
			if field.Point {
				e.points = append(e.points, NewPoint(x, y))
			}
			if field.Box != EMPTY {
				box := NewBox(NewPoint(x, y), boxId)
				e.boxes[boxId] = &box
				e.boxesOrdered[boxId] = &box
			}
		}
		y++
		e.Surface = append(e.Surface, []Field{})
	}
	// the last sub-array of Surface is always empty, so remove it...
	if len(e.Surface[len(e.Surface)-1]) == 0 {
		e.Surface = e.Surface[:len(e.Surface)-1]
	}
	return
}

// loop over all points and check, if there is a box. Else return false
func (e *Engine) Won() bool {
	for _, p := range e.points {
		if e.Surface[p.Y][p.X].Box == EMPTY {
			return false
		}
	}
	return true
}

// print the current Surface
func (e *Engine) Print() {
	log.Lock <- 1
	var x, y int8
	for y = 0; y < int8(len(e.Surface)); y++ {
		log.A("%3d ", e.Id)
		for x = 0; x < int8(len(e.Surface[y])); x++ {
			switch field := e.Surface[y][x]; {
			case field.Wall:
				log.A("#")
			case e.figPos.X == x && e.figPos.Y == y:
				if field.Point {
					log.A("+")
				} else {
					log.A("x")
				}
			case field.Box == EMPTY:
				if field.Point {
					log.A("*")
				} else if field.Dead {
					log.A("☠")
				} else {
					log.A(" ")
				}
			default: // field has box
				if field.Point {
					log.A("%%")
				} else {
					log.A("$")
				}
			}
			log.A(" ")
		}
		log.A("\n")
	}
	fieldnr, deadnr := e.Surface.AmountOfFields()
	log.A("Boxes: %d\n", len(e.Boxes()))
	log.A("Points: %d\n", len(e.Points()))
	log.A("Fields: %d\n", fieldnr)
	log.A("DeadFields: %d\n", deadnr)
	<-log.Lock
}

// return Point array of all boxes and the figure
func (e *Engine) GetBoxesAndX() (field []Point) {
	field = append(field, e.figPos)

	for i := int8(1); i <= int8(len(e.boxesOrdered)); i++ {
		field = append(field, e.boxesOrdered[i].Pos)
	}
	return
}

// print a legend of the Surface output
func PrintInfo() {
	log.Lock <- 1
	log.A("Surface Field association:\n")
	log.A("EMPTY\t\t' '\n")
	log.A("BOX\t\t'$'\n")
	log.A("FIGURE\t\t'x'\n")
	log.A("EMPTY POINT\t'*'\n")
	log.A("BOX POINT\t'%'\n")
	log.A("FIGURE POINT\t'+'\n")
	log.A("WALL\t\t'#'\n")
	log.A("DEAD FIELD\t'☠'\n")
	<-log.Lock
}

// return number of boxes on the surface
func (surface Surface) CountBoxes() int8 {
	count := int8(0)
	for y := 0; y < len(surface); y++ {
		for x := 0; x < len(surface[y]); x++ {
			if surface[y][x].Box != EMPTY {
				count++
			}
		}
	}
	return count
}

// check if the surface border was reached
func (surface Surface) In(p Point) bool {
	if p.Y < 0 || p.X < 0 || p.Y >= int8(len(surface)) || p.X >= int8(len(surface[p.Y])) {
		return false
	}
	return true
}

// read from the specified file and return whole content in a string
func readLevelAsString(filename string) (string, error) {
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
			if err == io.EOF {
				break
			}
			return "", err
		}
	}
	return string(result), nil
}
