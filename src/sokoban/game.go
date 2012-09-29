package main

import (
	"fmt"
	"os"
	"sokoban/ai"
	"sokoban/engine"
	"sokoban/log"
	"strconv"
)

func main() {
	log.DebugLevel = 3
	runmode := false
	single := true
	level := "alevel"
	straightAhead := false
	outputFreq := 50000
	printSurface := false

	e := engine.NewEngine()

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
				engine.PrintInfo()
			case "-m":
				single = false
			case "-d":
				if len(os.Args) > i+1 {
					debuglevel, err := strconv.Atoi(os.Args[i+1])
					if err == nil {
						log.DebugLevel = debuglevel
					}
				}
			case "-s":
				straightAhead = true
			case "-f":
				if len(os.Args) > i+1 {
					of, err := strconv.Atoi(os.Args[i+1])
					if err == nil {
						outputFreq = of
					}
				}
			case "-p":
				printSurface = true
			}
		}
	}

	e.LoadLevel(level)
	log.D("Level: " + level)

	if runmode {
		ai.Run(e, single, outputFreq, printSurface, straightAhead)
		return
	}
	
	// surface
	e.Print()

	var choice string
	for {
		choice = ""
		fmt.Print("Press m for manual or r for run: ")
		fmt.Scanf("%s", &choice)
		if choice == "r" {
			ai.Run(e, single, outputFreq, printSurface, straightAhead)
			break
		} else if choice == "m" {
			fmt.Println("Manual mode")
			var input int
			for {
				fmt.Scanf("%d", &input)
				if input >= 0 && input <= 3 {
					e.Move(engine.Direction(input))
					e.Print()
				} else {
					e.UndoStep()
					e.Print()
				}
			}
			break
		}
	}
}
