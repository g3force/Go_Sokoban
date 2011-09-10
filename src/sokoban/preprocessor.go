package sokoban

func MarkDeadFields() {
	for y := 0; y < len(Surface); y++ {
		for x := 0; x < len(Surface[y]); x++ {
			if !Surface[y][x].wall {
				thisPoint := Point{x,y}
				dead, dir1 := DeadEnd(thisPoint)
				if dead{
					Surface[y][x].dead=true
					dir1 = (dir1 + 2)%4 //dir1, dir2 are the directions of a possible dead wall
					dir2 := (dir1 - 1)%4 
					deadWall, p := checkForDeadWall(thisPoint, dir1, (dir2 + 2)%4)
					if deadWall {
						markDeadWall(thisPoint, p)
					}
					deadWall, p = checkForDeadWall(thisPoint, dir2, (dir1 + 2)%4)
					if deadWall {
						markDeadWall(thisPoint, p)
					}
				}
			}
		}
	}
}
//deadEdge: first dead Edge to star
//dir: direction where the wall will go on
//wallDir: direction ot the wall, left or right of the dir???
func checkForDeadWall(deadEdge Point, dir int, wallDir int) (bool, Point) {
	possDead := deadEdge
	for {
		possDead = addPoints(possDead, Direction(dir))
		possField := Surface[possDead.Y][possDead.X]
		possWall := Surface[addPoints(possDead, Direction(wallDir)).Y][addPoints(possDead, Direction(wallDir)).X]
		if possField.wall || possField.point || !possWall.wall{
			return false, possDead
		} else {
			dead, _ := DeadEnd(possDead)
			if dead {
				return true, possDead
			}
		}			
	}
	E("chekcForDeadWall: end of For loop")
	return false, possDead
}

func markDeadWall(start Point, end Point) {
	if start.X==end.X && start.Y!=end.Y   {
		if start.Y < end.Y {
			for i := start.Y; i <= end.Y; i++{
				Surface[i][start.X].dead=true
			}
		} else {
			for i := start.Y; i >= end.Y; i--{
				Surface[i][start.X].dead=true
			}
		}
	} else if start.Y==end.Y && start.X!=end.X  {
		if start.X<end.X {
			for i := start.X; i <= end.X; i++{
				Surface[start.Y][i].dead=true
			}
		} else {
			for i := start.X; i >= end.X; i--{
				Surface[start.Y][i].dead=true
			}
		}
	} else {
		I("Solo dead end")
	}
}
		
		
		
