package logger

import (
	"fmt"
	"time"
)

func Log(a ...interface{}) {
	prefix := time.Now().Format(time.DateTime) + "::"
	log := append([]interface{}{prefix}, a...)
	fmt.Println(log...)
}

func Error(a ...interface{}) {
	prefix := time.Now().Format(time.DateTime) + ":: ERROR:"
	log := append([]interface{}{prefix}, a...)
	fmt.Println(log...)
}
