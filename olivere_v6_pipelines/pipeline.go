//Package olivere_v6_pipelines implements modifications pipeline for version 6 of olivere/elastic files.
package olivere_v6_pipelines

import (
	"fmt"
	"gitlab.com/dmitry.konovenschi/aggretastic-sync/cmd"
	"os"
	"strings"
)

type olivere_v6_vars struct {
	repo         string
	repoPath     string
	repoHeadLock string
	repoBranch   string

	buildPath             string
	elasticExportPatterns []string
	deps                  []string
}

//load required variables from env
func loadVars() olivere_v6_vars {
	patterns := os.Getenv("ELASTIC_EXPORT_PATTERNS")
	files := os.Getenv("AGGRETASTIC_PACKAGE_FILES")
	if files == "" {
		panic("This pipeline requires aaha aggregation types. Please, specify them in Conf.env file. ")
	}
	return olivere_v6_vars{
		repo:                  os.Getenv("ELASTIC_REPO"),
		repoPath:              os.Getenv("REPO_PATH"),
		repoHeadLock:          os.Getenv("HEAD_LOCK_FILE"),
		repoBranch:            os.Getenv("REPO_BRANCH"),
		buildPath:             os.Getenv("BUILD_PATH"),
		elasticExportPatterns: strings.Split(patterns, ", "),
		deps:                  strings.Split(files, ", "),
	}
}

//run olivere_v6 pipeline
func Run() {
	vars := loadVars()

	//clean tmp files after finish
	defer cmd.Rm(vars.repoPath)
	defer cmd.Rm(vars.buildPath)

	//run git pipeline
	git := gitPipeline{
		Url:    vars.repo,
		Branch: vars.repoBranch,
		Path:   vars.repoPath,
		Lock:   vars.repoHeadLock,
	}
	if git.Run() {
		fmt.Println("Aggretastic is already up-to-date with " + vars.repoBranch)
		return
	}

	//run package updater
	updater := packageUpdaterPipeline{
		BuildPath: vars.buildPath,
		RepoPath:  vars.repoPath,
		Patterns:  vars.elasticExportPatterns,
		Deps:      vars.deps,
	}
	updater.Run()

	generator := sourceGenerationPipeline{
		SyncNodes:updater.SyncNodes,
	}
	generator.run()

	buildPackage(vars.buildPath)
}

//copy updater artifacts in main repository and remove deprecated
func buildPackage(buildPath string) {
	originFileList := cmd.LsByPattern("./", "^aggs_(.*).go$")
	buildFileList := cmd.Ls(buildPath)
	deprecated := cmd.ListDiff(originFileList, buildFileList)

	//copy new files to project
	cmd.ExtractFiles(".*", buildPath, "./")

	//remove deprecated files
	cmd.RmList("./", deprecated)
	if len(deprecated) > 0 {
		fmt.Println("Deprecated files has been removed:")
		for _, file := range deprecated {
			fmt.Println(file.Name())
		}
	}
}
