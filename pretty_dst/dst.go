// Package pretty_dst is a high-level wrapper for go/ast and go/dst.
//
// Abstract Syntax Trees transformations in Aggretastic-generator are used for operations on source-files.
//
// pretty_dst implement some high-level operations which are useful for aggretastic-generator.
// You can easily extend ops list in future, if future elastic releases will contain new file-patterns.
package pretty_dst

import (
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dkonovenschi/aggretastic-sync/errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"regexp"
)

//Container for Dst and it fileset
type Source struct {
	FileSet *token.FileSet
	Dst     *dst.File
}

//Creates new Source Container
func NewDst(file io.Reader) (src *Source) {
	var err error
	src = &Source{FileSet: token.NewFileSet()}
	src.Dst, err = decorator.ParseFile(src.FileSet, "", file, parser.ParseComments)
	errors.PanicOnError(err, nil)
	return src
}

//Find structure by name pattern
func (src *Source) FindStructure(pattern string) *StructureDeclaration {
	regex, _ := regexp.Compile(pattern)
	query := newStructureSearchQuery(regex)
	dst.Walk(query, src.Dst)
	if query.structure != nil {
		return query.structure
	}
	return nil
}

//Find function by name pattern
func (src *Source) FindFunction(pattern string) *Function {
	regex, _ := regexp.Compile(pattern)
	query := newFunctionSearchQuery(regex)
	dst.Walk(query, src.Dst)
	if query.function != nil {
		return query.function
	}
	return nil
}

//renames  package
func (src *Source) RenamePackage(name string) {
	src.Dst.Name.Name = name
}

//add new package import
func (src *Source) AddImport(name string, path string) {
	importSpec := NewImportSpec(name, path)

	//create import declaration if not exists
	if src.Dst.Imports == nil {
		dec := NewImportDeclaration(importSpec)
		src.Dst.Decls = append([]dst.Decl{dec}, src.Dst.Decls...)
	} else {
		AppendToImport(src.Dst.Decls[0].(*dst.GenDecl), importSpec)
	}
}

//save changes on disk
func (src *Source) Save(file io.Writer) error {
	fs, fl, err := decorator.RestoreFile(src.Dst)
	if err != nil {
		return err
	}

	if err := printer.Fprint(file, fs, fl); err != nil {
		return err
	}
	return nil
}

//restore ast and fset from dst
func (src *Source) Restore() (*token.FileSet, *ast.File) {
	fs, fl, _ := decorator.RestoreFile(src.Dst)
	return fs, fl
}
