//Package cmd is a high-level wrapper on bash functions.
package cmd

import (
	"fmt"
	"gitlab.com/dmitry.konovenschi/aggretastic-sync/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

var (
	errAccessDenied = fmt.Errorf("Access denied! ")
	errFileDoesNotExists = fmt.Errorf("You try to copy nonexistent file! ")
	errPathDoesNotExists = fmt.Errorf("You try to ls nonexistent path! ")
)

//Wrapper for mkdir -p
func MakePath(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	errors.PanicOnError(err, errAccessDenied)
}

//wrapper for cp
func Cp(args ...string) {
	err := exec.Command("cp", args...).Run()
	errors.PanicOnError(err, errFileDoesNotExists)
}

//wrapper for ls
func Ls(dir string) []os.FileInfo {
	content, err := ioutil.ReadDir(dir)
	errors.PanicOnError(err, errPathDoesNotExists)
	return content
}

//ls and extract entities by pattern
func LsByPattern(dir string, pattern string) []os.FileInfo {
	regex := regexp.MustCompile(pattern)
	content := Ls(dir)
	filtered := []os.FileInfo{}
	for _, file := range content {
		if regex.MatchString(file.Name()) {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

//return entities from first list which has not been found in second list
func ListDiff(first []os.FileInfo, second []os.FileInfo) []os.FileInfo {
	diff := []os.FileInfo{}
	for _, file := range first {
		if ListContains(file, second) {
			continue
		}
		diff = append(diff, file)
	}
	return diff
}

// check if entity contains in list
func ListContains(entry os.FileInfo, list []os.FileInfo) bool {
	for _, file := range list {
		if file.Name() == entry.Name() {
			return true
		}
	}
	return false
}

//wrapper for rm
func Rm(filepath string) {
	err := os.RemoveAll(filepath)
	if err != nil {
		errors.PanicOnError(err, errAccessDenied)
	}
}

//remove all entities from list in path
func RmList(path string, list []os.FileInfo) {
	for _, file := range list {
		Rm(path + file.Name())
	}
}
