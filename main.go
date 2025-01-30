package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func CleanInput(text string) []string {
	split := strings.Fields(strings.ToLower(text))

	return split
}