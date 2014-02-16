package main

import (
	"fmt"
	"os"
	"github.com/g3force/Go_Sokoban/ai"
	"github.com/g3force/Go_Sokoban/engine"
	"github.com/g3force/Go_Sokoban/log"
	"strconv"
)

func main() {
	log.DebugLevel = 4
	runmode := false
	single := true
	level := "alevel"
	straightAhead := false
	outputFreq := int32(50000)
	printSurface := false
	threads := 1

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
						outputFreq = int32(of)
					}
				}
			case "-p":
				printSurface = true
			case "-t":
				if len(os.Args) > i+1 {
					t, err := strconv.Atoi(os.Args[i+1])
					if err != nil {
						panic(err)
					} else {
						threads = t
					}
				}
			}
		}
	}

	e.LoadLevel(level)
	log.I(e.Id, "Level: " + level)

	if runmode {
		ai.Run(e, single, outputFreq, printSurface, straightAhead, threads)
		return
	}
	
	// surface
	e.Print()

	var choice string
	for {
		choice = ""
		log.A("Press m for manual or r for run: ")
		fmt.Scanf("%s", &choice)
		if choice == "r" {
			ai.Run(e, single, outputFreq, printSurface, straightAhead, threads)
			break
		} else if choice == "m" {
			log.A("Manual mode\n")
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
