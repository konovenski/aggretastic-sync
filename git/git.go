//Package git is a wrapper for Git tool.
package git

import (
	"bytes"
	"fmt"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"
	"os"
)

//wrapper for Git clone
func Clone(url, branch string) (*git.Repository, error) {
	fmt.Println("Git Clone:")
	fs := memfs.New()
	storage := memory.NewStorage()

	repository, err := git.Clone(storage, fs, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName(branch),
		SingleBranch:  true,
		Progress:      os.Stdout,
	})

	return repository, err
}

//get current Head hash
func getLastCommit(repository *git.Repository) ([]byte, error) {
	head, err := repository.Head()
	if err != nil {
		return nil, err
	}
	hash := head.Hash().String()
	return []byte(hash), nil
}

//if repo head hash is equal to hash from lock file - return true
func IsUpToDate(repository *git.Repository, lockPath string) (bool, error) {
	hash, err := getLastCommit(repository)
	if err != nil {
		return false, err
	}
	content, _ := ioutil.ReadFile(lockPath) //we can ignore errors here, because []byte(nil) != anyHash
	return bytes.Equal(hash, content), nil
}

//creates head lock file
func CreateLockFile(lockpath string, hash plumbing.Hash) error {
	f, err := os.Create(lockpath)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(hash.String()))
	if err != nil {
		return err
	}
	return f.Close()
}
