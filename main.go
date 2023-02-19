package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var maxLine int
var Red = color.New(color.FgRed).SprintFunc()
var Green = color.New(color.FgGreen).SprintFunc()
var err error
var totalFileWithExceededLines = 0

func main() {
	maxLine = 1000

	args := os.Args
	mydir, _ := os.Getwd()

	splittedPath := strings.Split(mydir, "/")
	currentDir := splittedPath[len(splittedPath)-1]

	path := "../" + currentDir + "/"
	if len(args) > 1 {
		path = "../" + currentDir + "/" + args[1]
	}

	if len(args) >= 3 {
		maxLine, err = strconv.Atoi(args[2])
		if err != nil {
			println("incorrect input on second param, default to 1000")
		}
	}

	PrintDir(path)
}

func PrintDir(dirName string) {
	if !CheckPath(dirName) {
		println("path does not exist")
		return
	}

	paths := strings.Split(dirName, "/")
	dirName = strings.Join(paths, "/")
	dirName = strings.TrimSuffix(dirName, "/")

	filepath.Walk(dirName, func(name string, info os.FileInfo, _ error) error {
		res := ""
		newPath := strings.ReplaceAll(name, "../", "")
		fileName := ""
		if strings.Contains(newPath, "/") {
			splitted := strings.Split(newPath, "/")[1:]
			GetSpace(&res, info.IsDir(), len(splitted))
			fileName = splitted[len(splitted)-1]
		}

		totalLine := GetLines(name)
		if !info.IsDir() {
			fileName += fmt.Sprintf(" (%d lines)", totalLine)
		}

		if totalLine > maxLine {
			println(res + Red(fileName))
			totalFileWithExceededLines++
		} else {
			println(res + fileName)
		}
		return nil
	})

	lines := "└"
	for i := 0; i < len(fmt.Sprint(totalFileWithExceededLines)); i++ {
		lines += "─"
	}
	lines += "───────────────────────────────"
	println(lines)
	if  totalFileWithExceededLines > 0{
		println(fmt.Sprintf(Red("%d files have exceeded lines limit"), totalFileWithExceededLines))
	} else {
		println(fmt.Sprintf(Green("%d files have exceeded lines limit"), totalFileWithExceededLines))
	}

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
