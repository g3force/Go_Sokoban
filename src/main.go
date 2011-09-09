package main

import (
	"sokoban"
	"fmt"
	"os"
	"strconv"
)

func main() {
	sokoban.DebugLevel = 3
	runmode := false
	single := true
	level := "level1"
	sokoban.StraightAhead = false
	outputFreq := 5000

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
			case "-d":
				if len(os.Args) > i+1 {
					debuglevel, err := strconv.Atoi(os.Args[i+1])
					if err == nil {
						sokoban.DebugLevel = debuglevel
					}
				}
			case "-s":
				sokoban.StraightAhead = true
			case "-f":
				if len(os.Args) > i+1 {
					of, err := strconv.Atoi(os.Args[i+1])
					if err == nil {
						outputFreq = of
					}
				}
			}
		}
	}

	sokoban.LoadLevel(level)
	sokoban.Init()
	sokoban.Print()

	if runmode {
		sokoban.Run(single, outputFreq)
		return
	}

	var choice string
	for {
		choice = ""
		fmt.Print("Press m for manual or r for run: ")
		fmt.Scanf("%s", &choice)
		if choice == "r" {
			sokoban.Run(single, outputFreq)
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

