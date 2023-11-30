package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	// write your code here
	if len(os.Args) < 2 {
		fmt.Println("Directory is not specified")
		os.Exit(1)
	}

	dir := os.Args[1]
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if path != dir && !info.IsDir() {
			fmt.Println(path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir, err)
		return
	}
}
