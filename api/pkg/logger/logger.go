package logger

import "fmt"

// Log logs by tag
func Log(tag string, data []interface{}) {
	for _, element := range data {
		fmt.Println("[", tag, "]", element)
	}
}
