package olivere_v6_pipelines

import "github.com/dkonovenschi/aggretastic-sync/git"

type gitPipeline struct {
	Url string
	Branch string
	Path string
	Lock string
}

//run git pipeline
func (g *gitPipeline) Run() bool {
	repository := git.Clone(g.Url, g.Branch, g.Path)
	utd := git.IsUpToDate(repository, g.Lock)
	head, _ := repository.Head()
	git.CreateLockFile(g.Lock, head.Hash())
	return utd
}
