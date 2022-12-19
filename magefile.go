//go:build mage
// +build mage

package main

import (
	"fmt"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// Just says hello
func Build() {
	fmt.Println("Hello from MageFile build..!")
	return
}

// Just says hello
func Hello() {
	fmt.Println("Hello from MageFile..!")
	return
}

var Default = Hello
