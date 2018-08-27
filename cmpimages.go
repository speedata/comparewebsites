package main

import (
	"fmt"
	"os"
	"path/filepath"

	"cmpimage"
	"screenshot"
)

// To build: run go build cmpimages.go

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s [.json*] [dir1 dir2]\n", filepath.Base(os.Args[0]))
		os.Exit(0)
	}
	jsonfiles := []string{}
	screenshotpaths := []string{}
	for _, v := range os.Args[1:] {
		fi, err := os.Stat(v)
		if err == nil && !fi.IsDir() {
			jsonfiles = append(jsonfiles, v)
		} else {
			screenshotpaths = append(screenshotpaths, v)
		}
	}

	if len(jsonfiles) > 0 {
		err := screenshot.Dothings(jsonfiles)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	if l := len(screenshotpaths); l > 0 {
		if l != 2 {
			fmt.Println("To compare the screenshots, you have to state exactly two directories")
			os.Exit(0)
		}
		cmpimage.Dothings(screenshotpaths)
	}
}
