package pretty_dst

import (
	"fmt"
	"github.com/dave/dst"
	"go/token"
)

//Function declaration decorator
type Function struct {
	name   *dst.Ident
	body   *FunctionBody
	origin *dst.FuncDecl
}

//Creates function decorator
func DecorateFunction(function *dst.FuncDecl) *Function {
	return &Function{
		name:   function.Name,
		body:   DecorateFunctionBody(function.Body),
		origin: function,
	}
}

//returns original dst function declaration
func (f *Function) Extract() *dst.FuncDecl {
	return f.origin
}

//returns function name
func (f *Function) GetName() string {
	return f.name.Name
}

//rename function
func (f *Function) Rename(newName string) {
	f.name.Name = newName
}

//returns reference to this function body
func (f *Function) GetBody() *FunctionBody {
	return f.body
}

// Function body decorator
type FunctionBody struct {
	*dst.BlockStmt
}

//creates function body decorator
func DecorateFunctionBody(body *dst.BlockStmt) *FunctionBody {
	return &FunctionBody{
		body,
	}
}

// returns first return statement in body
func (body *FunctionBody) GetFirstReturn() *ReturnStatement {
	for _, stmt := range body.List {
		if statement, ok := stmt.(*dst.ReturnStmt); ok {
			return DecorateReturnStatement(statement)
		}
	}
	return nil
}

// returns first assigment atatement in body
func (body *FunctionBody) GetFirstAssignment() *dst.AssignStmt {
	for _, stmt := range body.List {
		if statement, ok := stmt.(*dst.AssignStmt); ok {
			return statement
		}
	}
	return nil
}

//remove all elements in function body
func (body *FunctionBody) Wipe() {
	body.List = []dst.Stmt{}
}

//appends new assigment to the function body
func (body *FunctionBody) AppendNewAssigment(name string, operator token.Token, expression dst.Expr) {
	assigment := NewAssigment(name, operator, expression)
	body.List = append(body.List, assigment)
}

//appends new return statement to the function body
func (body *FunctionBody) AppendNewReturn(args ...dst.Expr) {
	assigment := NewReturnStatement(args...)
	body.List = append(body.List, assigment)
}

//Structure declaration decorator
type StructureDeclaration struct {
	name   *dst.Ident
	fields *dst.FieldList
	origin *dst.TypeSpec
}

//creates structure decorator
func DecorateStructure(dstType *dst.TypeSpec) (*StructureDeclaration, error) {
	if !IsStructure(dstType) {
		return nil, fmt.Errorf("not struct ")
	}

	return &StructureDeclaration{
		name:   dstType.Name,
		fields: dstType.Type.(*dst.StructType).Fields,
		origin: dstType,
	}, nil
}

//returns original dst structure declaration
func (s *StructureDeclaration) Extract() *dst.TypeSpec {
	return s.origin
}

//returns structure name
func (s *StructureDeclaration) GetName() string {
	return s.name.Name
}

//rename structure
func (s *StructureDeclaration) Rename(newName string) {
	s.name.Name = newName
}

//declares new field in structure
func (s *StructureDeclaration) AddField(name string, _type string) {
	field := NewDataField(name, _type)
	s.fields.List = append(s.fields.List, field)
}

//removes field from structure by name
func (s *StructureDeclaration) RemoveField(name string) error {
	index, _ := s.FindField(name)
	if index < 0 {
		return fmt.Errorf("Trying to remove unexisting field ")
	}
	s.removeField(index)
	return nil
}

//find field position in FieldList
func (s *StructureDeclaration) FindField(fieldName string) (int, *dst.Field) {
	for key, field := range s.fields.List {
		for _, name := range field.Names {
			if name.Name == fieldName {
				return key, field
			}
		}
	}
	return -1, nil
}

//removes field from FieldList by index
func (s *StructureDeclaration) removeField(id int) {
	s.fields.List = append(s.fields.List[:id], s.fields.List[id+1:]...)
}
