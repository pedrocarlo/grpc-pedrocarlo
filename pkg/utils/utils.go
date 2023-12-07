package utils

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func Log_trace(msg ...any) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	function_names := strings.Split(frame.Function, ".")
	var function_name string = frame.Function
	if len(function_names) > 0 {
		function_name = function_names[len(function_names)-1]
	}
	log.Printf("%s %s\n", function_name, fmt.Sprint(msg...))
}

func Log_fatal_trace(err error) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Printf("%s %s\n", frame.Function, err)
}
