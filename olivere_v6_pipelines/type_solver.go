package olivere_v6_pipelines

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
)

//"\"github.com/olivere/elastic\""

type typeSolverPipeline struct {
	BuildPath         string
	OriginPackagePath string
	OriginPackageName string
	FilesToCheck      []os.FileInfo
	typesInfo         types.Info
	typesConfig       *types.Config
	astSet            []*ast.File
	fileSet           *token.FileSet
}

//run type solver pipeline
func (ts *typeSolverPipeline) Run() bool {
	ts.parseFiles()
	ts.newTypesInfo()
	ts.newTypesConfig()

	//run error correction pipeline on error
	err := ts.check()
	if err != nil {
		fmt.Println("Fix procedure on: " + err.Error())
		ec := errorCorrectionPipeline{
			Err:err.Error(),
			OriginPackagePath:ts.OriginPackagePath,
			OriginPackageName:ts.OriginPackageName,
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
			astFile, _ := parser.ParseFile(ts.fileSet, name, nil, parser.ParseComments)
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