package main

import (
	"sokoban/ai"
	"fmt"
	"os"
	"strconv"
	"log"
)

func main() {
	ai.DebugLevel = 3
	runmode := false
	single := true
	level := "alevel"
	ai.StraightAhead = false
	outputFreq := 50000
	printSurface := false

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
				ai.PrintInfo()
			case "-m":
				single = false
			case "-d":
				if len(os.Args) > i+1 {
					debuglevel, err := strconv.Atoi(os.Args[i+1])
					if err == nil {
						ai.DebugLevel = debuglevel
					}
				}
			case "-s":
				ai.StraightAhead = true
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

	ai.LoadLevel(level)
	ai.Init()
	log.Print("Level: " + level)
	ai.Print()
	//fmt.Printf("boxes: %d, points: %d\n", len(ai.GetBoxes()), len(ai.GetPoints()))

	if runmode {
		ai.Run(single, outputFreq, printSurface)
		return
	}

	var choice string
	for {
		choice = ""
		fmt.Print("Press m for manual or r for run: ")
		fmt.Scanf("%s", &choice)
		if choice == "r" {
			ai.Run(single, outputFreq, printSurface)
			break
		} else if choice == "m" {
			fmt.Println("Manual mode")
			var input int
			for {
				fmt.Scanf("%d", &input)
				if input >= 0 && input <= 3 {
					ai.Move(int8(input))
					ai.Print()
				} else {
					ai.UndoStep()
					ai.Print()
				}
			}
			break
		}
	}
}
