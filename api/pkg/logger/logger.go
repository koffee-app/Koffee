package logger

import (
	"log"
)

// Log logs by tag
func Log(tag string, data ...interface{}) {
	for _, element := range data {
		log.Printf("[ %s ]: %v\n", tag, element)
	}
}

// Assert .
func Assert(i interface{}) {
	if i != nil {
		log.Println(i)
		panic(i)
	}
}
