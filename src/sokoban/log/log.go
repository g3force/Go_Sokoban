package log

import "fmt"

// 0 nothing
// 1 Errors
// 2 Warnings 	+ 1
// 3 Info 		+ 2
// 4 Debug 		+ 3
var DebugLevel = 4

func debug(tag string, message string, args ...interface{}) {
	fmt.Printf(tag+"\t"+message+"\n", args...)
}

func E(message string, args ...interface{}) {
	if DebugLevel > 0 {
		debug("Error  ", message, args...)
	}
}

func W(message string, args ...interface{}) {
	if DebugLevel > 1 {
		debug("Warning", message, args...)
	}
}

func D(message string, args ...interface{}) {
	if DebugLevel > 3 {
		debug("Debug  ", message, args...)
	}
}

func I(message string, args ...interface{}) {
	if DebugLevel > 2 {
		debug("Info   ", message, args...)
	}
}

