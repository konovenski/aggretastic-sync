package cmd

import (
	"github.com/dlclark/regexp2"
	"path/filepath"
)

//extract all matched entities from dir to other dir
func ExtractFiles(reg string, from string, to string) {
	pattern := regexp2.MustCompile(reg, 0)
	MakePath(to)
	content := Ls(from)
	for _, entry := range content {
		CopyOnMatch(pattern, entry.Name(), from, to)
	}
}

//copy file if it matches pattern
func CopyOnMatch(pattern *regexp2.Regexp, filename string, from string, to string) {
	if isMatch, _ := pattern.MatchString(filename); isMatch {
		file := filepath.Join(from, filename)
		Cp(file, to)
	}
}
