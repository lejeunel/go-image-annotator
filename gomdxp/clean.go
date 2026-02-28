package main

import (
	"fmt"
	"os"
)

func Clean() error {
	fmt.Println("Cleaning...")

	err := os.RemoveAll("path/to/dir")
	if err != nil {
		panic(err)
	}
	return nil
}
