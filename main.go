package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var Red = color.New(color.FgRed).SprintFunc()
var Green = color.New(color.FgGreen).SprintFunc()

var maxLine int
var err error
var totalFileWithExceededLines = 0
var showFolderDetail = false
var excludedFormat = ""
var excludedFile = map[string]bool{}
var totalLines = 0
var totalFolderInDir = 0
var totalFilesInDir = 0

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

		if len(args) >= 4 {
			showFolderDetail = false
			if args[3] == "true" {
				showFolderDetail = true
			}

			if len(args) >= 5 {
				excludedFormat = ""
				excludedFormat = args[4]
			}
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
	if excludedFormat != "" {
		RegisterExcludedFiles(dirName)
	}

	filepath.Walk(dirName, func(name string, info os.FileInfo, _ error) error {
		if info.IsDir() && excludedFormat != "" {
			RegisterExcludedFiles(name)
		}

		if info.IsDir() {
			return nil
		}

		res := ""
		newPath := strings.ReplaceAll(name, "../", "")
		fileName := ""
		if strings.Contains(newPath, "/") {
			splitted := strings.Split(newPath, "/")[1:]
			GetSpace(&res, info.IsDir(), len(splitted))
			fileName = splitted[len(splitted)-1]
		}

		if _, shouldExclude := excludedFile[fileName]; shouldExclude {
			return nil
		}

		totalLine := GetLines(name)
		totalFiles, totalFolders := 0, 0

		if !info.IsDir() {
			fileName += fmt.Sprintf(" (%d lines)", totalLine)
			totalLines += totalLine
			totalFilesInDir++
			} else {
			totalFolderInDir++
			totalFiles, totalFolders = GetTotalFilesAndFolders(name)
		}

		fileFolderInfo := ""
		if info.IsDir() && showFolderDetail {
			fileFolderInfo = fmt.Sprintf(" [%d folder(s) and %d file(s)]", totalFolders, totalFiles)
		}

		if totalLine > maxLine {
			println(res + Red(fileName) + fileFolderInfo)
			totalFileWithExceededLines++
		} else {
			if info.IsDir() {
				println(color.New(color.FgHiWhite, color.Bold).SprintFunc()(res + fileName + fileFolderInfo))
			} else {
				println(res + fileName + fileFolderInfo)
			}
		}
		return nil
	})

	lines := "└"
	for i := 0; i < len(fmt.Sprint(totalFileWithExceededLines)+fmt.Sprint(maxLine)); i++ {
		lines += "─"
	}
	lines += "─────────────────────────────────────────"
	println(lines)
	if totalFileWithExceededLines > 0 {
		println(fmt.Sprintf(Red("%d files have exceeded lines limit [limit: %d]"), totalFileWithExceededLines, maxLine))
		println(fmt.Sprintf(Red("Total lines in all files are: %d]"), totalLines))
		println(fmt.Sprintf(Red("Total files are: %d"), totalFilesInDir))
		println(fmt.Sprintf(Red("Total folders are: %d"), totalFolderInDir))
		} else {
			println(fmt.Sprintf(Green("%d files have exceeded lines limit [limit: %d]"), totalFileWithExceededLines, maxLine))
			println(fmt.Sprintf(Green("Total lines in all files are: %d]"), totalLines))
			println(fmt.Sprintf(Green("Total files are: %d"), totalFilesInDir))
			println(fmt.Sprintf(Green("Total folders are: %d"), totalFolderInDir))
	}
}

func RegisterExcludedFiles(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), excludedFormat) {
			excludedFile[file.Name()] = true
		}
	}
}

func GetTotalFilesAndFolders(path string) (int, int) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	totalFile, totalFolder := 0, 0

	for _, file := range files {
		if file.IsDir() {
			totalFolder++
		} else {
			totalFile++
		}
	}
	return totalFile, totalFolder
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
