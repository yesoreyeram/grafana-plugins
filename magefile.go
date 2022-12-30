//go:build mage

package main

import (
	"fmt"
)

// Build not implemented yet
func Build() {
	fmt.Println("Hello from MageFile build..!")
	return
}

// Just says hello
func Hello() {
	fmt.Println("Hello from MageFile..!")
	return
}

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build
