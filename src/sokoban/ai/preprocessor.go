package ai

import (
	"sokoban/engine"
	"sokoban/log"
)

// check, if given point is a dead corner
func DeadCorner(surface engine.Surface, point engine.Point) (found bool, x int8) {
	var p engine.Point
	hit := false
	x = 0
	found = false
	if surface[point.Y][point.X].Point {
		return
	}

	// check clockwise, if there is a wall or not.
	// If there is a wall two times together, corner is dead
	for i := 0; i < 5; i++ {
		x = int8(i % 4)
		p = point.Add(engine.Direction(x).Point())
		if !surface.In(p) || surface[p.Y][p.X].Wall {
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

func MarkDeadFields(surface *engine.Surface) {
	for y := 0; y < len(*surface); y++ {
		for x := 0; x < len((*surface)[y]); x++ {
			thisPoint := engine.NewPoint(x, y)
			// walls can't be dead fields
			if !surface.In(thisPoint) || (*surface)[y][x].Wall {
				continue
			}
			dead, dir1 := DeadCorner(*surface, thisPoint)
			if dead {
				(*surface)[y][x].Dead = true
				dir1 = (dir1 + 2) % 4 //dir1, dir2 are the directions of a possible dead wall
				dir2 := (dir1 - 1) % 4
				deadWall, p := checkForDeadWall(*surface, thisPoint, dir1, (dir2+2)%4)
				if deadWall {
					markDeadWall(surface, thisPoint, p)
				}
				deadWall, p = checkForDeadWall(*surface, thisPoint, dir2, (dir1+2)%4)
				if deadWall {
					markDeadWall(surface, thisPoint, p)
				}
			}
		}
	}
}
//deadEdge: first dead Edge to star
//dir: direction where the wall will go on
//wallDir: direction of the wall, left or right of the dir???
func checkForDeadWall(surface engine.Surface, deadEdge engine.Point, dir int8, wallDir int8) (bool, engine.Point) {
	possDead := deadEdge
	for {
		possDead = possDead.Add(engine.Direction(dir).Point())
		if !surface.In(possDead) {
			return false, possDead
		}
		possField := surface[possDead.Y][possDead.X]
		possWallPos := possDead.Add(engine.Direction(wallDir).Point())
		if !surface.In(possWallPos) {
			return false, possDead
		}
		possWall := surface[possWallPos.Y][possWallPos.X]
		if possField.Wall || possField.Point || !possWall.Wall {
			return false, possDead
		} else {
			dead, _ := DeadCorner(surface, possDead)
			if dead {
				return true, possDead
			}
		}
	}
	log.E(-1, "checkForDeadWall: end of For loop")
	return false, possDead
}

func markDeadWall(surface *engine.Surface, start engine.Point, end engine.Point) {
	if start.X == end.X && start.Y != end.Y {
		if start.Y < end.Y {
			for i := start.Y; i <= end.Y; i++ {
				(*surface)[i][start.X].Dead = true
			}
		} else {
			for i := start.Y; i >= end.Y; i-- {
				(*surface)[i][start.X].Dead = true
			}
		}
	} else if start.Y == end.Y && start.X != end.X {
		if start.X < end.X {
			for i := start.X; i <= end.X; i++ {
				(*surface)[start.Y][i].Dead = true
			}
		} else {
			for i := start.X; i >= end.X; i-- {
				(*surface)[start.Y][i].Dead = true
			}
		}
	} else {
		log.D(-1, "Solo dead end")
	}
	return
}

