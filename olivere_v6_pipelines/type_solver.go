package olivere_v6_pipelines

import (
	"fmt"
	"github.com/dkonovenschi/aggretastic-sync/errors"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"gopkg.in/src-d/go-billy.v4"
	"os"
)

type typeSolverPipeline struct {
	BuildPath         string
	OriginPackagePath string
	OriginPackageName string
	FilesToCheck      []os.FileInfo
	typesInfo         types.Info
	typesConfig       *types.Config
	astSet            []*ast.File
	fileSet           *token.FileSet
	FS                billy.Filesystem
}

//run type solver pipeline
func (ts *typeSolverPipeline) Run() bool {
	ts.parseFiles()
	ts.newTypesInfo()
	ts.newTypesConfig()

	//run error correction pipeline on error
	err := ts.check()
	if err != nil {
		fmt.Print("|")
		ec := errorCorrectionPipeline{
			Err:               err.Error(),
			OriginPackagePath: ts.OriginPackagePath,
			OriginPackageName: ts.OriginPackageName,
			FS:                ts.FS,
		}
		ec.Run()
		return true
	}
	return false
}

//parse ast from source file
func (ts *typeSolverPipeline) parseFiles() {
	ts.astSet = []*ast.File{}
	ts.fileSet = token.NewFileSet()

	for _, file := range ts.FilesToCheck {
		if !file.IsDir() {
			name := ts.BuildPath + file.Name()
			f, err := ts.FS.Open(name)
			errors.PanicOnError(errCantOpenFile, err)

			astFile, err := parser.ParseFile(ts.fileSet, name, f, parser.ParseComments)

			errors.PanicOnError(errCantParseFile, err)
			ts.astSet = append(ts.astSet, astFile)
		}
	}
}

//run type checker
func (ts *typeSolverPipeline) check() error {
	_, err := ts.typesConfig.Check("", ts.fileSet, ts.astSet, &ts.typesInfo)
	return err
}

func (ts *typeSolverPipeline) newTypesInfo() {
	ts.typesInfo = types.Info{
		Defs: make(map[*ast.Ident]types.Object),
		Uses: make(map[*ast.Ident]types.Object),
	}
}

func (ts *typeSolverPipeline) newTypesConfig() {
	ts.typesConfig = &types.Config{
		Importer: importer.Default(),
	}
}
/*
func lookup (path string) (io.ReadCloser, error) {

}

func Import(packages map[string]*types.Package, path, srcDir string, lookup func(path string) (io.ReadCloser, error)) (pkg *types.Package, err error) {
	var rc io.ReadCloser
	var id string
	if lookup != nil {
		// With custom lookup specified, assume that caller has
		// converted path to a canonical import path for use in the map.
		if path == "unsafe" {
			return types.Unsafe, nil
		}
		id = path

		// No need to re-import if the package was imported completely before.
		if pkg = packages[id]; pkg != nil && pkg.Complete() {
			return
		}
		f, err := lookup(path)
		if err != nil {
			return nil, err
		}
		rc = f
	} else {
		var filename string
		filename, id = FindPkg(path, srcDir)
		if filename == "" {
			if path == "unsafe" {
				return types.Unsafe, nil
			}
			return nil, fmt.Errorf("can't find import: %q", id)
		}

		// no need to re-import if the package was imported completely before
		if pkg = packages[id]; pkg != nil && pkg.Complete() {
			return
		}

		// open file
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err != nil {
				// add file name to error
				err = fmt.Errorf("%s: %v", filename, err)
			}
		}()
		rc = f
	}

}
*/