package olivere_v6_pipelines

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"gitlab.com/dmitry.konovenschi/aggretastic-sync/pretty_dst"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strconv"
)

type errorCorrectionPipeline struct {
	Err               string
	OriginPackagePath string
	OriginPackageName string
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
	//compile regex for col/line
	findLineColRegex := regexp2.MustCompile(":\\d*", 0)
	//get line and filename
	group, _ := findLineColRegex.FindStringMatch(ec.Err)
	ec.filename = ec.Err[:group.Index]
	ec.line, _ = strconv.Atoi(group.String()[1:])
	//get column
	group, _ = findLineColRegex.FindNextMatch(group)
	ec.column, _ = strconv.Atoi(group.String()[1:])
}

//extract AST from file
func (ec *errorCorrectionPipeline) parseFile() {
	ec.fileSet = token.NewFileSet()
	ec.ast, _ = parser.ParseFile(ec.fileSet, ec.filename, nil, parser.ParseComments)
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
	src, _ := pretty_dst.NewDst(ec.filename)
	src.AddImport("", ec.OriginPackagePath)
	src.Save(ec.filename)
}

//find error pointer in ast and fix error
func (ec *errorCorrectionPipeline) fixError() {
	vis := &errorSleuth{
		column: ec.column,
		line:   ec.line,
		fset:   ec.fileSet,
		packageName: ec.OriginPackageName,
	}
	ast.Walk(vis, ec.ast)
}

//save changes on disk
func (ec *errorCorrectionPipeline) saveFile() {
	f, _ := os.Create(ec.filename)
	defer f.Close()
	if err := printer.Fprint(f, ec.fileSet, ec.ast); err != nil {
		log.Fatal(err)
	}
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
