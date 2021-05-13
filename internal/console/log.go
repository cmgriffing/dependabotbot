package console

import (
	"fmt"
	"log"
)

func Log(v ...interface{}) {
	fmt.Println(v...)
}

func Error(v ...interface{}) {
	log.Fatal(v...)
}
