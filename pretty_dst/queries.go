package pretty_dst

import (
	"github.com/dave/dst"
	"regexp"
)

//Container for Structure search query
type structureSearchQuery struct {
	namePattern *regexp.Regexp
	structure   *StructureDeclaration
}

//Creates new Structure search query
func newStructureSearchQuery(regexp *regexp.Regexp) *structureSearchQuery {
	return &structureSearchQuery{
		namePattern: regexp,
	}
}

//Search structure in Node
func (q *structureSearchQuery) Visit(node dst.Node) dst.Visitor {
	structure, found := findStructureByPattern(node, q.namePattern)
	if found {
		q.structure = structure
	}
	return q
}

func findStructureByPattern(node dst.Node, pattern *regexp.Regexp) (*StructureDeclaration, bool) {
	n, isType := IsTypeDeclaration(node)
	if !isType {
		return nil, false
	}

	for _, spec := range n.Specs {
		structure, err := DecorateStructure(spec.(*dst.TypeSpec))
		if err == nil && pattern.MatchString(structure.GetName()) {
			return structure, true
		}
	}
	return nil, false
}

//Container for Function search query
type functionSearchQuery struct {
	namePattern *regexp.Regexp
	function    *Function
}

//Creates new Function search query
func newFunctionSearchQuery(regexp *regexp.Regexp) *functionSearchQuery {
	return &functionSearchQuery{
		namePattern: regexp,
	}
}

//Search function in Node
func (q *functionSearchQuery) Visit(node dst.Node) dst.Visitor {
	function, found := findFunctionByPattern(node, q.namePattern)
	if found {
		q.function = function
	}
	return q
}

func findFunctionByPattern(node dst.Node, pattern *regexp.Regexp) (*Function, bool) {
	n, isFunction := IsFunction(node)
	if !isFunction {
		return nil, false
	}

	function := DecorateFunction(n)

	if pattern.MatchString(function.GetName()) {
		return function, true
	}
	return nil, false
}
