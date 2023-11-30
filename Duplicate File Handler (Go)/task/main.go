package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	// write your code here
	if len(os.Args) < 2 {
		fmt.Println("Directory is not specified")
		os.Exit(1)
	}

	dir := os.Args[1]

	fmt.Println("Enter file format:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	fileFormat := scanner.Text()

	fmt.Println("Size sorting options:")
	fmt.Println("1. Descending")
	fmt.Println("2. Ascending")

	var sortingOption int

	for {
		fmt.Println("Enter a sorting option:")
		fmt.Scan(&sortingOption)
		if sortingOption != 2 && sortingOption != 1 {
			fmt.Println("Wrong option")
		} else {
			break
		}
	}

	fileData := make(map[int64][]string)
	sizes := make([]int64, 0, 50)

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if path != dir && !info.IsDir() {
			splitByDot := strings.Split(path, ".")
			extension := splitByDot[len(splitByDot)-1]
			if fileFormat == "" || fileFormat == extension {
				size := info.Size()
				sizes = append(sizes, size)
				fileData[size] = append(fileData[size], path)
				//fmt.Println(path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir, err)
		return
	}

	if sortingOption == 2 {
		sort.Slice(sizes, func(i, j int) bool {
			return sizes[i] < sizes[j]
		})
	} else {
		sort.Slice(sizes, func(i, j int) bool {
			return sizes[i] > sizes[j]
		})
	}

	for index, size := range sizes {
		if index > 0 && sizes[index-1] == size {
			// already did this size
			continue
		} else {
			fmt.Printf("\n%d bytes\n", size)
		}

		for _, path := range fileData[size] {
			fmt.Println(path)
		}
	}

	fmt.Println("Check for duplicates?")
	var dupCheck string
	fmt.Scan(&dupCheck)
	if dupCheck == "yes" {
		number := 0
		for index, size := range sizes {
			if index > 0 && sizes[index-1] == size {
				// already did this size
				continue
			} else {
				hashVals := make(map[string][]string)
				fmt.Printf("\n%d bytes\n", size)
				for _, path := range fileData[size] {
					// calculate hash
					f, err := os.Open(path)
					if err != nil {
						log.Fatal(err)
					}
					defer f.Close()

					h := md5.New()
					if _, err := io.Copy(h, f); err != nil {
						log.Fatal(err)
					}
					hashVal := fmt.Sprintf("%x", (h.Sum(nil)))
					hashVals[hashVal] = append(hashVals[hashVal], path)
					//fmt.Printf("Hash: %s\n", h)
				}

				for h, s := range hashVals {
					if len(s) > 1 {
						fmt.Printf("Hash: %s\n", h)
						for _, p := range s {
							number++
							fmt.Printf("%d. %s\n", number, p)
						}
					}
				}
			}
		}
	}
}
