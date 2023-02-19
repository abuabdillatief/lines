package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args
	mydir, _ := os.Getwd()

	splittedPath := strings.Split(mydir, "/")
	currentDir := splittedPath[len(splittedPath)-1]

	path := "../" + currentDir + "/"
	if len(args) > 1 {
		path = "../" + currentDir + "/" + args[1]
	}

	PrintDir(path)
}

func PrintDir(dirName string) {
	if !CheckPath(dirName) {
		println("path does not exist")
		return
	}

	filepath.Walk(dirName, func(name string, info os.FileInfo, err error) error {
		res := ""
		newPath := strings.ReplaceAll(name, "../", "")
		if strings.Contains(newPath, "/") {
			splitted := strings.Split(newPath, "/")[1:]
			GetSpace(&res, info.IsDir(), len(splitted))
			res += splitted[len(splitted)-1]
		}

		if !info.IsDir() {
			res += fmt.Sprintf(" (%d lines)", GetLines(name))
		}
		println(res)
		return nil
	})
}

// CheckPath will return false if target path does not exist
func CheckPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetSpace will add proper space for better readability in current path
func GetSpace(str *string, isDir bool, len int) {
	*str = "│"
	pointer := ""
	separator := " "

	if len <= 2 {
		*str = "├"
	}

	if isDir || (!isDir && len <= 2) {
		separator = "─"
	}

	for i := 0; i < len; i++ {
		*str += separator
	}

	if !isDir && len >= 2 {
		if len > 2 {
			pointer = "├─ "
		}

		for i := 0; i <= len; i++ {
			if i == len-1 {
				*str += pointer
			}
		}
	}
}

// GetLines will count how many lines a file has
func GetLines(path string) int {
	file, _ := os.Open(path)
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	return lineCount
}
