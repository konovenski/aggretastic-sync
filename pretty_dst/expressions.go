package pretty_dst

import "github.com/dave/dst"

//UnaryExpr decorator
type UnaryExpression struct {
	*dst.UnaryExpr
}

//creates new UnaryExpr decorator
func DecorateUnaryExpression(expression *dst.UnaryExpr) *UnaryExpression {
	return &UnaryExpression{
		expression,
	}
}

//returns ComposuteLit from unaryExpr if contains any
func (ue *UnaryExpression) GetCompositeLiteral() *CompositeLiteral {
	if ue == nil {
		return nil
	}
	if compLit, ok := ue.X.(*dst.CompositeLit); ok {
		return DecorateCompositeLiteral(compLit)
	}
	return nil
}

//CompositeLit Decorator
type CompositeLiteral struct {
	*dst.CompositeLit
}

//creates new CompositeLit decorator
func DecorateCompositeLiteral(compLit *dst.CompositeLit) *CompositeLiteral {
	return &CompositeLiteral{
		compLit,
	}
}

//Removes element from compositeLit by key
func (cl *CompositeLiteral) RemoveElementByKey(key string) {
	if index := cl.GetKVElementIndex(key); index >= 0 {
		cl.RemoveElementByIndex(index)
	}
}

//Returns matched KVelement index if any
func (cl *CompositeLiteral) GetKVElementIndex(key string) int {
	for index, element := range cl.Elts {
		if element.(*dst.KeyValueExpr).Key.(*dst.Ident).Name == key {
			return index
		}
	}
	return -1
}

//removes element from composite literal by index
func (cl *CompositeLiteral) RemoveElementByIndex(index int) {
	if index >= len(cl.Elts) {
		return
	}
	cl.Elts = append(cl.Elts[:index], cl.Elts[index+1:]...)
}

//Decorator for ReturnStmt
type ReturnStatement struct {
	*dst.ReturnStmt
}

//Creates new ReturnStmt Decorator
func DecorateReturnStatement(statement *dst.ReturnStmt) *ReturnStatement {
	return &ReturnStatement{
		statement,
	}
}

//Returns unaryExpression from ReturnStmt results if contains any
func (rs *ReturnStatement) GetUnaryExpression() *UnaryExpression {
	for _, result := range rs.Results {
		if expression, ok := result.(*dst.UnaryExpr); ok {
			return DecorateUnaryExpression(expression)
		}
	}
	return nil
}

//Returns Ident from ReturnStmt results if contains any
func (rs *ReturnStatement) GetIdentifier() *dst.Ident {
	for _, result := range rs.Results {
		if identifier, ok := result.(*dst.Ident); ok {
			return identifier
		}
	}
	return nil
}
