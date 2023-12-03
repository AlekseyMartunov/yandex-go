// Package main here, just test file uses in testing ExitAnalyzer
package main

import (
	"fmt"
	"os"
)

func main() {
	os.Exit(1) // want "calling os.Exit in main is prohibited"
	myFunc()
}

func myFunc() {
	fmt.Println("do something")
	os.Exit(1)
}
