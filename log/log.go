package log

import (
	"fmt"
//	"os"
)

// 0 nothing
// 1 Errors
// 2 Warnings 	+ 1
// 3 Info 		+ 2
// 4 Debug 		+ 3
var DebugLevel = 4
var Lock = make(chan int, 1)

func debug(tag string, worker int, newLine bool, message string, args ...interface{}) {
	Lock <- 1
	line := fmt.Sprintf(tag+"\t%d\t", (worker)) + fmt.Sprintf(message+"\n", args...)
	fmt.Print(line)
	<-Lock
//	fileName := fmt.Sprintf("worker_%d.log", worker)
//	file, err := os.Create(fileName) // For read access.
//    if err == nil {
////	defer file.Close()
//    	file.WriteString(line)
//    }
}

func A(message string, args ...interface{}) {
	fmt.Printf(message, args...)
}

func E(worker int, message string, args ...interface{}) {
	if DebugLevel > 0 {
		debug("Error  ", worker, true, message, args...)
	}
}

func W(worker int, message string, args ...interface{}) {
	if DebugLevel > 1 {
		debug("Warning", worker, true, message, args...)
	}
}

func D(worker int, message string, args ...interface{}) {
	if DebugLevel > 3 {
		debug("Debug  ", worker, true, message, args...)
	}
}

func I(worker int, message string, args ...interface{}) {
	if DebugLevel > 2 {
		debug("Info   ", worker, true, message, args...)
	}
}

