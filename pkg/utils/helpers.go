package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// Generic helper functions on common types
func wd() {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	fmt.Println("Current Directory:", dir)

	// Walk through the directory and print the structure
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Print the file or directory
		fmt.Println(path)
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}
