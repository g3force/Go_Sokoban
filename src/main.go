package main

import (
	"sokoban"
	"fmt"
	"os"
)

func main() {
	sokoban.DebugLevel = 3
	runmode := false
	single := true
	level := "level1"
	sokoban.StraightAhead = false 

	if len(os.Args) > 1 {
		for i, _ := range os.Args {
			switch os.Args[i] {
			case "-r":
				runmode = true
			case "-l":
				if len(os.Args) > i+1 {
					level = os.Args[i+1]

				}
			case "-i":
				sokoban.PrintInfo()
			case "-m":
				single = false
			case "-s":
				sokoban.StraightAhead = true
			}
		}
	}

	sokoban.LoadLevel(level)
	sokoban.Init()
	sokoban.Print()

	if runmode {
		sokoban.Run(single)
		return
	}

	var choice string
	for {
		choice = ""
		fmt.Print("Press m for manual or r for run: ")
		fmt.Scanf("%s", &choice)
		if choice == "r" {
			sokoban.Run(single)
			break
		} else if choice == "m" {
			fmt.Println("Manual mode")
			var input int
			for {
				fmt.Scanf("%d", &input)
				if input >= 0 && input <= 3 {
					sokoban.Move(input)
					sokoban.Print()
				} else {
					sokoban.UndoStep()
					sokoban.Print()
				}
			}
			break
		}
	}
}

