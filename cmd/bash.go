//Package cmd is a high-level wrapper on bash functions.
package cmd

import (
	"bytes"
	"gopkg.in/src-d/go-billy.v4"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

//wrapper for cp
func Cp(fs billy.Filesystem, from string, to string) error {
	src, err := fs.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := fs.Create(to)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)

	return err
}

//wrapper for cp
func CpFromMemory(fs billy.Filesystem, from string, to string) error {
	src, err := fs.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()

	data, err := ioutil.ReadAll(src)
	err = ioutil.WriteFile(to, data, os.ModePerm)

	return err
}

func CpFromReal(fs billy.Filesystem, from string, to string) error {
	file, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	dstFile, err := fs.Create(to)
	if err != nil {
		return err
	}
	_, err = io.Copy(dstFile, bytes.NewReader(file))
	return err
}

//ls and extract entities by pattern
func LsDiskByPattern(dir string, pattern string) ([]os.FileInfo, error) {
	regex := regexp.MustCompile(pattern)
	content, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filtered := []os.FileInfo{}
	for _, file := range content {
		if regex.MatchString(file.Name()) {
			filtered = append(filtered, file)
		}
	}
	return filtered, nil
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


//remove all entities from list in path
func RmListFromDisk(path string, list []os.FileInfo) error {
	for _, file := range list {
		err := os.Remove(path + file.Name())
		if err != nil {
			return err
		}
	}
	return nil
}
