package olivere_v6_pipelines

import (
	"github.com/dkonovenschi/aggretastic-sync/errors"
	"github.com/dkonovenschi/aggretastic-sync/git"
	"gopkg.in/src-d/go-billy.v4"
)

type gitPipeline struct {
	Url    string
	Branch string
	Lock   string
}

//run git pipeline
func (g *gitPipeline) Run() (bool, billy.Filesystem) {
	repository, err := git.Clone(g.Url, g.Branch)
	errors.PanicOnError(errCantClone, err)

	isUpToDate, err := git.IsUpToDate(repository, g.Lock)
	errors.PanicOnError(errBrokenRepo, err)

	head, err := repository.Head()
	errors.PanicOnError(errBrokenRepo, err)

	err = git.CreateLockFile(g.Lock, head.Hash())
	errors.PanicOnError(errCantCreateLock, err)

	fs, err := repository.Worktree()
	errors.PanicOnError(errBrokenStorage, err)

	return isUpToDate, fs.Filesystem
}
