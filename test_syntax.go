package main

import (
	"fmt"
	"strings"
)

// TestFunction demonstrates syntax highlighting
func TestFunction() {
	// This is a comment
	var message string = "Hello, world!"
	number := 42
	flag := true
	
	if flag {
		fmt.Println(message)
		for i := 0; i < number; i++ {
			fmt.Printf("Count: %d\n", i)
		}
	}
	
	result := strings.ToUpper(message)
	fmt.Println(result)
}

const MaxValue = 100

type User struct {
	Name string
	Age  int
}

