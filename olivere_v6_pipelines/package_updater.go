package olivere_v6_pipelines

import (
	"fmt"
	"github.com/dkonovenschi/aggretastic-sync/cmd"
	"github.com/dkonovenschi/aggretastic-sync/errors"
	"gopkg.in/src-d/go-billy.v4"
)



type packageUpdaterPipeline struct {
	BuildPath string
	RepoPath  string
	Patterns  []string
	Deps      []string
	FS        billy.Filesystem
}

//run package updater pipeline
func (up *packageUpdaterPipeline) Run() {
	up.extractRequiredFiles()
	up.enrichFiles()
	up.extractDeps()
	up.runTypeSolver()
}

//extract files from repo to build path
func (up *packageUpdaterPipeline) extractRequiredFiles() {
	for _, pattern := range up.Patterns {
		err := cmd.ExtractFiles(up.FS, pattern, ".", up.BuildPath)
		errors.PanicOnError(errCantCopyFile, err)
	}
}

//extract other dependencies to build path
func (up *packageUpdaterPipeline) extractDeps() {
	for _, file := range up.Deps {
		err := cmd.CpFromReal(up.FS, file, up.BuildPath+file)
		errors.PanicOnError(errCantCopyFile, err)
	}
}

//change every file ast and save on disk
func (up *packageUpdaterPipeline) enrichFiles() {
	files, err := up.FS.ReadDir(up.BuildPath)
	errors.PanicOnError(errCantReadDir, err)

	fmt.Print("Update process:[")
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		fmt.Print("|")
		fu := fileUpdatePipeline{
			Filename:                   up.BuildPath + name,
			DesiredPackageName:         "aggretastic",
			TargetStructureNamePattern: "(.*)Aggregation$",
			TargetFunctionNamePattern:  "^New(.*)Aggregation$",
			FS:                         up.FS,
		}
		fu.Run()
	}
	fmt.Println("]")
}

//run type solver
func (up *packageUpdaterPipeline) runTypeSolver() {
	fileList, err := up.FS.ReadDir(up.BuildPath)
	errors.PanicOnError(errCantReadDir, err)

	ts := typeSolverPipeline{
		BuildPath:         up.BuildPath,
		OriginPackageName: "elastic",
		OriginPackagePath: "github.com/olivere/elastic",
		FilesToCheck:      fileList,
		FS:                up.FS,
	}
	fmt.Print("Fix process: [")
	for i := true; i != false; i = ts.Run() {
	}
	fmt.Println("]")
}
