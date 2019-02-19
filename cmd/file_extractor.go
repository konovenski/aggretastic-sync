package cmd

import (
	"github.com/dlclark/regexp2"
	"gopkg.in/src-d/go-billy.v4"
	"os"
	"path/filepath"
)

//extract all matched entities from dir to other dir
func ExtractFiles(fs billy.Filesystem, reg, from, to string) (err error) {
	err = fs.MkdirAll(to, os.ModePerm)
	if err != nil {
		return err
	}

	pattern := regexp2.MustCompile(reg, 0)
	content, err := fs.ReadDir(from)
	if err != nil {
		return err
	}

	for _, entry := range content {
		err = CopyOnMatch(fs, pattern, entry.Name(), from, to)
	}
	return nil
}

//extract all matched entities from dir to other dir
func ExtractFilesFromMemory(fs billy.Filesystem, reg, from, to string) (err error) {
	pattern := regexp2.MustCompile(reg, 0)
	//MakePath(fs, to)
	content, err := fs.ReadDir(from)
	if err != nil {
		return err
	}

	for _, entry := range content {
		err = CopyFromMemoryOnMatch(fs, pattern, entry.Name(), from, to)
	}
	return err
}

//copy file if it matches pattern
func CopyOnMatch(fs billy.Filesystem, pattern *regexp2.Regexp, filename, from, to string) (err error) {
	isMatch, err := pattern.MatchString(filename)
	if isMatch {
		fileFrom := filepath.Join(from, filename)
		fileTo := filepath.Join(to, filename)
		err = Cp(fs, fileFrom, fileTo)
	}
	return err
}

//copy file from memory to disk if it matches pattern
func CopyFromMemoryOnMatch(fs billy.Filesystem, pattern *regexp2.Regexp, filename, from, to string) (err error) {
	isMatch, err := pattern.MatchString(filename)
	if isMatch {
		fileFrom := filepath.Join(from, filename)
		fileTo := filepath.Join(to, filename)
		err = CpFromMemory(fs, fileFrom, fileTo)
	}
	return err
}
