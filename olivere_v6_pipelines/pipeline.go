//Package olivere_v6_pipelines implements modifications pipeline for version 6 of olivere/elastic files.
package olivere_v6_pipelines

import (
	"fmt"
	"github.com/dkonovenschi/aggretastic-sync/cmd"
	"github.com/dkonovenschi/aggretastic-sync/errors"
	"gopkg.in/src-d/go-billy.v4"
	"log"
	"os"
	"strings"
)

var (
	errCantAtoi      = fmt.Errorf("Atoi conversion can't be executed: ")

	errCantOpenFile  = fmt.Errorf("Can't open file: ")
	errCantCloseFile  = fmt.Errorf("Can't close file: ")
	errCantReadDir = fmt.Errorf("Can't read from directory: ")
	errCantWriteFile  = fmt.Errorf("Can't write to file: ")
	errCantRemoveFile = fmt.Errorf("File can't be removed: ")
	errCantCopyFile = fmt.Errorf("Can't copy file: ")
	errCantParseFile = fmt.Errorf("Source file can't be parsed: ")

	errCantClone      = fmt.Errorf("Repository can't be cloned: ")
	errBrokenRepo     = fmt.Errorf("Repository is broken: ")
	errCantCreateLock = fmt.Errorf("Can't create lock file: ")
	errBrokenStorage  = fmt.Errorf("Repository storage is broken: ")

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
		repo: os.Getenv("ELASTIC_REPO"),
		repoHeadLock: os.Getenv("HEAD_LOCK_FILE"),
		repoBranch:   os.Getenv("REPO_BRANCH"),
		buildPath:             "build-tmp/",
		elasticExportPatterns: strings.Split(patterns, ", "),
		deps:                  strings.Split(files, ", "),
	}
}

//run olivere_v6 pipeline
func Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(fmt.Sprintf("Sync process can't be finished because: %v ", r))
		}
	}()

	vars := loadVars()

	//run git pipeline
	git := gitPipeline{
		Url:    vars.repo,
		Branch: vars.repoBranch,
		Lock:   vars.repoHeadLock,
	}
	isUpToDate, fs := git.Run()

	if isUpToDate {
		log.Println("Aggretastic is already up-to-date with " + vars.repoBranch)
		return
	}

	//run package updater
	updater := packageUpdaterPipeline{
		RepoPath:  vars.repoPath,
		Patterns:  vars.elasticExportPatterns,
		Deps:      vars.deps,
		FS:        fs,
		BuildPath: vars.buildPath,
	}
	updater.Run()

	buildPackage(fs, vars.buildPath)
}

//copy updater artifacts in main repository and remove deprecated
func buildPackage(fs billy.Filesystem, buildPath string) {
	originFileList, err := cmd.LsDiskByPattern("./", "^aggs_(.*).go$")
	errors.PanicOnError(errCantReadDir, err)

	buildFileList, err := fs.ReadDir(buildPath)
	errors.PanicOnError(errCantReadDir, err)

	deprecated := cmd.ListDiff(originFileList, buildFileList)

	//copy new files to project
	err = cmd.ExtractFilesFromMemory(fs, ".*", buildPath, "./")
	errors.PanicOnError(errCantCopyFile, err)

	//remove deprecated files
	err = cmd.RmListFromDisk("./", deprecated)
	errors.PanicOnError(errCantRemoveFile, err)

	if len(deprecated) > 0 {
		fmt.Println("Deprecated files has been removed:")
		for _, file := range deprecated {
			fmt.Println(file.Name())
		}
	}
}
