package olivere_v6_pipelines

import (
	"fmt"
	"gitlab.com/dmitry.konovenschi/aggretastic-sync/cmd"
)

type packageUpdaterPipeline struct {
	BuildPath string
	RepoPath  string
	Patterns  []string
	Deps      []string
	SyncNodes []*syncNode
}

//run package updater pipeline
func (up *packageUpdaterPipeline) Run() {
	up.SyncNodes = []*syncNode{}
	up.extractRequiredFiles()
	up.enrichFiles()
	up.extractDeps()
	up.runTypeSolver()
}

//extract files from repo to build path
func (up *packageUpdaterPipeline) extractRequiredFiles() {
	for _, pattern := range up.Patterns {
		cmd.ExtractFiles(pattern, up.RepoPath, up.BuildPath)
	}
}

//extract other dependencies to build path
func (up *packageUpdaterPipeline) extractDeps() {
	for _, file := range up.Deps {
		cmd.Cp(file, up.BuildPath)
	}
}

//change every file ast and save on disk
func (up *packageUpdaterPipeline) enrichFiles() {
	files := cmd.Ls(up.BuildPath)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		fmt.Println(name)
		fu := fileUpdatePipeline{
			Filename:                   up.BuildPath + name,
			DesiredPackageName:         "aggretastic",
			TargetStructureNamePattern: "(.*)Aggregation$",
			TargetFunctionNamePattern:  "^New(.*)Aggregation$",
		}
		if node := fu.Run(); node != nil {
			up.SyncNodes = append(up.SyncNodes, node)
		}
	}
}

//run type solver
func (up *packageUpdaterPipeline) runTypeSolver() {
	ts := typeSolverPipeline{
		BuildPath:         up.BuildPath,
		OriginPackageName: "elastic",
		OriginPackagePath: "github.com/olivere/elastic",
		FilesToCheck:      cmd.Ls(up.BuildPath),
	}
	for i := true; i != false; i = ts.Run() {
	}
}
