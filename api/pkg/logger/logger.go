package logger

import "fmt"

// Log logs by tag
func Log(tag string, data ...interface{}) {
	for _, element := range data {
		fmt.Printf("[ %s ]: %v\n", tag, element)
	}
}

// Assert .
func Assert(i interface{}) {
	if i != nil {
		fmt.Println(i)
		panic(i)
	}
}
