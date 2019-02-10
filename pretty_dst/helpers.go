package pretty_dst

import (
	"fmt"
	"github.com/dave/dst"
	"go/token"
)

//creates new Variable object
func NewVariableObject(name string) *dst.Object {
	return dst.NewObj(dst.Var, name)
}

//creates new Left-Hand Side expression
func NewLhs(name string) []dst.Expr {
	return []dst.Expr{
		NewIdent(name, nil),
	}
}

//creates new Right-Hand Side expression from dst.Expr
func NewRhs(expression dst.Expr) []dst.Expr {
	return []dst.Expr{
		expression,
	}
}

//creates new Return Statement
func NewReturnStatement(args ...dst.Expr) *dst.ReturnStmt {
	return &dst.ReturnStmt{
		Results: args,
	}
}

//creates new Identifier
func NewIdent(name string, obj *dst.Object) *dst.Ident {
	return &dst.Ident{
		Name: name,
		Obj:  obj,
	}
}

//creates new fully-completed Assigment
func NewAssigment(name string, operator token.Token, expression dst.Expr) *dst.AssignStmt {
	return &dst.AssignStmt{
		Lhs: NewLhs(name),
		Tok: operator,
		Rhs: NewRhs(expression),
	}
}

//creates new import spec
func NewImportSpec(name string, path string) *dst.ImportSpec {
	var nameIdent *dst.Ident
	if len(name) > 0 {
		nameIdent = NewIdent(name, nil)
	}

	return &dst.ImportSpec{
		Name: nameIdent,
		Path: &dst.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s\"", path),
		},
	}
}

//creates new import declaration
func NewImportDeclaration(imps ...dst.Spec) *dst.GenDecl {
	return &dst.GenDecl{
		Tok:    token.IMPORT,
		Lparen: true,
		Specs: imps,
		Rparen: true,
	}
}

//appends importSpec to import declaration
func AppendToImport(impDec *dst.GenDecl, imp dst.Spec) {
	impDec.Specs = append(impDec.Specs, imp)
	impDec.Lparen = true
	impDec.Rparen = true
}

//creates new function call expression
func NewCallExpression(name string, arguments ...string) *dst.CallExpr {
	fun := NewIdent(name, nil)
	args := []dst.Expr{}
	for _, arg := range arguments {
		args = append(args, NewIdent(arg, nil))
	}
	return &dst.CallExpr{
		Fun:  fun,
		Args: args,
	}
}

//creates new filled Data Field (for structures and interfaces)
func NewDataField(name string, _type string) *dst.Field {
	variable := NewVariableObject(name)

	field := &dst.Field{
		Names: []*dst.Ident{{
			Name: name,
			Obj:  variable,
		}},
	}

	variable.Decl = field

	field.Type = &dst.Ident{
		Name: _type,
	}

	return field
}

//Checks if TypeSpec contain Structure object
func IsStructure(dstType *dst.TypeSpec) bool {
	_, ok := dstType.Type.(*dst.StructType)
	return ok
}

//Checks if Node contain Function object
func IsFunction(node dst.Node) (*dst.FuncDecl, bool) {
	if node == nil {
		return nil, false
	}
	if n, ok := node.(*dst.FuncDecl); ok {
		return n, true
	}
	return nil, false
}

//Checks if Node contain Type Declaration
func IsTypeDeclaration(node dst.Node) (*dst.GenDecl, bool) {
	if node == nil {
		return nil, false
	}

	n, ok := node.(*dst.GenDecl)

	if !ok || n.Tok != token.TYPE {
		return nil, false
	}
	return n, true
}

//Checks if Structure contains field with such name
func (s *StructureDeclaration) IsFieldExists(fieldName string) bool {
	_, field := findField(s.fields, fieldName)
	if field != nil {
		return true
	}
	return false
}
