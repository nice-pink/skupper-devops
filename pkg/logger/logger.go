package logger

import (
	"fmt"
	"time"
)

func Log(a ...interface{}) {
	log := append([]interface{}{time.Now().Format(time.RFC850)}, a...)
	fmt.Println(log...)
}

func Error(a ...interface{}) {
	prefix := time.Now().Format(time.RFC850) + ":: ERROR:"
	log := append([]interface{}{prefix}, a...)
	fmt.Println(log...)
}
