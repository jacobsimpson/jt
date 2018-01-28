package debug

import (
	"fmt"
)

var level = 0

func SetLevel(l int) {
	level = l
}

func Info(format string, a ...interface{}) (n int, err error) {
	if level >= 1 {
		return fmt.Printf(format+"\n", a...)
	}
	return 0, nil
}

func Debug(format string, a ...interface{}) (n int, err error) {
	if level >= 2 {
		return fmt.Printf(format+"\n", a...)
	}
	return 0, nil
}
