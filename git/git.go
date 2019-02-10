//Package git is a wrapper for Git tool.
package git

import (
	"bytes"
	"fmt"
	"github.com/dkonovenschi/aggretastic-sync/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"io/ioutil"
	"os"
)

var (
	errBrokenRepo   = fmt.Errorf("Broken repository! Clean cache and retry. ")
	errClone        = fmt.Errorf("Couldn't clone the repository. ")
	errAccessDenied = fmt.Errorf("Access denied! Fix rights for this directory and try again. ")
)

//wrapper for Git clone
func Clone(url, branch, path string) *git.Repository {
	fmt.Println("Git Clone:")
	repository, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName(branch),
		SingleBranch:  true,
		Progress:      os.Stdout,
	})

	errors.PanicOnError(err, errClone)
	return repository
}

//get current Head hash
func getLastCommit(repository *git.Repository) []byte {
	head, err := repository.Head()
	errors.PanicOnError(err, errBrokenRepo)
	hash := head.Hash().String()
	return []byte(hash)
}

//if repo head hash is equal to hash from lock file - return true
func IsUpToDate(repository *git.Repository, lockPath string) bool {
	hash := getLastCommit(repository)
	content, _ := ioutil.ReadFile(lockPath) //we can ignore errors here, because []byte(nil) != anyHash
	return bytes.Equal(hash, content)
}

//creates head lock file
func CreateLockFile(lockpath string, hash plumbing.Hash) {
	f, err := os.Create(lockpath)
	errors.PanicOnError(err, errAccessDenied)
	_, _ = f.Write([]byte(hash.String()))
	_ = f.Close()
}
