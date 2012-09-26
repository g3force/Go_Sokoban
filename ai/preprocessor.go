package ai

// check, if given point is a dead corner
func DeadCorner(point Point) (found bool, x int8) {
	var p Point
	hit := false
	x = 0
	found = false
	if Surface[point.Y][point.X].point {
		return
	}

	// check clockwise, if there is a wall or not.
	// If there is a wall two times together, corner is dead
	for i := 0; i < 5; i++ {
		x = int8(i % 4)
		p = AddPoints(point, Direction(x))
		if !IsInSurface(p) || Surface[p.Y][p.X].wall {
			if hit {
				found = true
				return
			} else {
				hit = true
			}
		} else {
			hit = false
		}
	}
	return
}

func MarkDeadFields() {
	for y := 0; y < len(Surface); y++ {
		for x := 0; x < len(Surface[y]); x++ {
			thisPoint := NewPoint(x, y)
			// walls can't be dead fields
			if !IsInSurface(thisPoint) || Surface[y][x].wall {
				continue
			}
			dead, dir1 := DeadCorner(thisPoint)
			if dead {
				Surface[y][x].dead = true
				dir1 = (dir1 + 2) % 4 //dir1, dir2 are the directions of a possible dead wall
				dir2 := (dir1 - 1) % 4
				deadWall, p := checkForDeadWall(thisPoint, dir1, (dir2+2)%4)
				if deadWall {
					markDeadWall(thisPoint, p)
				}
				deadWall, p = checkForDeadWall(thisPoint, dir2, (dir1+2)%4)
				if deadWall {
					markDeadWall(thisPoint, p)
				}
			}
		}
	}
}
//deadEdge: first dead Edge to star
//dir: direction where the wall will go on
//wallDir: direction of the wall, left or right of the dir???
func checkForDeadWall(deadEdge Point, dir int8, wallDir int8) (bool, Point) {
	possDead := deadEdge
	for {
		possDead = AddPoints(possDead, Direction(dir))
		if !IsInSurface(possDead) {
			return false, possDead
		}
		possField := Surface[possDead.Y][possDead.X]
		possWallPos := AddPoints(possDead, Direction(wallDir))
		if !IsInSurface(possWallPos) {
			return false, possDead
		}
		possWall := Surface[possWallPos.Y][possWallPos.X]
		if possField.wall || possField.point || !possWall.wall {
			return false, possDead
		} else {
			dead, _ := DeadCorner(possDead)
			if dead {
				return true, possDead
			}
		}
	}
	E("checkForDeadWall: end of For loop")
	return false, possDead
}

func markDeadWall(start Point, end Point) {
	if start.X == end.X && start.Y != end.Y {
		if start.Y < end.Y {
			for i := start.Y; i <= end.Y; i++ {
				Surface[i][start.X].dead = true
			}
		} else {
			for i := start.Y; i >= end.Y; i-- {
				Surface[i][start.X].dead = true
			}
		}
	} else if start.Y == end.Y && start.X != end.X {
		if start.X < end.X {
			for i := start.X; i <= end.X; i++ {
				Surface[start.Y][i].dead = true
			}
		} else {
			for i := start.X; i >= end.X; i-- {
				Surface[start.Y][i].dead = true
			}
		}
	} else {
		I("Solo dead end")
	}
}

