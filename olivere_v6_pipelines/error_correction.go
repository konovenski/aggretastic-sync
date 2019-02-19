package olivere_v6_pipelines

import (
	"fmt"
	"github.com/dkonovenschi/aggretastic-sync/errors"
	"github.com/dkonovenschi/aggretastic-sync/pretty_dst"
	"github.com/dlclark/regexp2"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"gopkg.in/src-d/go-billy.v4"
	"strconv"
)



type errorCorrectionPipeline struct {
	Err               string
	OriginPackagePath string
	OriginPackageName string
	FS                billy.Filesystem
	filename          string
	line              int
	column            int
	fileSet           *token.FileSet
	ast               *ast.File
}

//runs error correction pipeline
func (ec *errorCorrectionPipeline) Run() {
	ec.parseError()
	ec.parseFile()
	ec.ensureImports()
	ec.fixError()
	ec.saveFile()
}

//extract filename, col and line from error message
func (ec *errorCorrectionPipeline) parseError() {
	var err error
	//compile regex for col/line
	findLineColRegex := regexp2.MustCompile(":\\d*", 0)

	//get line and filename
	group, _ := findLineColRegex.FindStringMatch(ec.Err)
	ec.filename = ec.Err[:group.Index]
	ec.line, err = strconv.Atoi(group.String()[1:])
	errors.PanicOnError(errCantAtoi, err)

	//get column
	group, _ = findLineColRegex.FindNextMatch(group)
	ec.column, err = strconv.Atoi(group.String()[1:])
	errors.PanicOnError(errCantAtoi, err)
}

func (ec *errorCorrectionPipeline) getSource() billy.File {
	reader, err := ec.FS.Open(ec.filename)
	errors.PanicOnError(errCantOpenFile, err)
	return reader
}

//extract AST from file
func (ec *errorCorrectionPipeline) parseFile() {
	ec.fileSet = token.NewFileSet()

	var err error
	ec.ast, err = parser.ParseFile(ec.fileSet, ec.filename, ec.getSource(), parser.ParseComments)
	errors.PanicOnError(errCantParseFile, err)
}

//ensure that we have only one given import declaration
func (ec *errorCorrectionPipeline) ensureImports() {
	if ec.ast.Imports != nil {
		for _, imp := range ec.ast.Imports {
			//do not change anything if package already imported
			if imp.Path.Value == fmt.Sprintf("\"%s\"", ec.OriginPackagePath) {
				return
			}
		}
	}
	//add import to file and regenerate ast tree
	ec.addImport()
	ec.parseFile()
}

//add import to ast and save to disk
func (ec *errorCorrectionPipeline) addImport() {
	src := pretty_dst.NewDst(ec.getSource())

	src.AddImport("", ec.OriginPackagePath)
	file, err := ec.FS.Create(ec.filename)
	errors.PanicOnError(errCantOpenFile, err)

	err = src.Save(file)
	errors.PanicOnError(errCantWriteFile, err)

	err = file.Close()
	errors.PanicOnError(errCantCloseFile, err)
}

//find error pointer in ast and fix error
func (ec *errorCorrectionPipeline) fixError() {
	vis := &errorSleuth{
		column:      ec.column,
		line:        ec.line,
		fset:        ec.fileSet,
		packageName: ec.OriginPackageName,
	}
	ast.Walk(vis, ec.ast)
}

//save changes on disk
func (ec *errorCorrectionPipeline) saveFile() {
	f, err := ec.FS.Create(ec.filename)
	errors.PanicOnError(errCantOpenFile, err)
	defer func() {
		err := f.Close()
		errors.PanicOnError(errCantCloseFile, err)
	}()

	err = printer.Fprint(f, ec.fileSet, ec.ast)
	errors.PanicOnError(errCantWriteFile, err)
}

//implementation of ast.Visitor interface
type errorSleuth struct {
	column      int
	line        int
	fset        *token.FileSet
	packageName string
}

//Find error pointer and fix error
func (v *errorSleuth) Visit(node ast.Node) ast.Visitor {
	n, ok := node.(*ast.Ident)
	if !ok {
		return v
	}

	//fix error if position matches
	position := v.fset.Position(n.Pos())
	if position.Line == v.line && position.Column == v.column {
		n.Name = fmt.Sprintf("%s.%s", v.packageName, n.Name)
	}
	return v
}
