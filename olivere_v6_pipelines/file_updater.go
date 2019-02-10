package olivere_v6_pipelines

import (
	"github.com/dave/dst"
	"github.com/dkonovenschi/aggretastic-sync/cmd"
	"github.com/dkonovenschi/aggretastic-sync/pretty_dst"
	"go/token"
	"strings"
)

type syncNode struct {
	aggregationType string
	functionName    string
}

func newSyncLeaf(aggrType string, functionName string) *syncNode {
	return &syncNode{
		aggregationType: aggrType,
		functionName:    functionName,
	}
}

type fileUpdatePipeline struct {
	Filename string
	src      *pretty_dst.Source

	DesiredPackageName string

	TargetStructureNamePattern string
	structure                  *pretty_dst.StructureDeclaration
	structureInitExpression    *dst.UnaryExpr

	TargetFunctionNamePattern string
	function                  *pretty_dst.Function
	functionBody              *pretty_dst.FunctionBody

	strategy modificationStrategy
}

//Run file updater pipeline
func (fu *fileUpdatePipeline) Run() *syncNode {
	fu.parseFile()
	defer fu.saveFile()

	fu.renamePackage()
	fu.findTargetStructure()
	if fu.structure == nil {
		return nil
	}
	fu.pickStrategy()
	fu.enrichStructure()

	fu.findTargetFunction()
	if fu.function == nil {
		return nil
	}
	fu.enrichFunction()
	return fu.fillSyncNode()
}

func (fu *fileUpdatePipeline) fillSyncNode() *syncNode {
	return newSyncLeaf(fu.strategy.name(), fu.function.GetName())
}

//parse ast from file
func (fu *fileUpdatePipeline) parseFile() {
	var err error
	fu.src, err = pretty_dst.NewDst(fu.Filename)
	if err != nil {
		panic(err)
	}
}

//save ast changes on disk without search_ prefix
func (fu *fileUpdatePipeline) saveFile() {
	filename := strings.Replace(fu.Filename, "search_aggs_", "aggs_", -1)
	fu.src.Save(filename)
	if filename != fu.Filename {
		cmd.Rm(fu.Filename)
	}
}

func (fu *fileUpdatePipeline) renamePackage() {
	fu.src.RenamePackage(fu.DesiredPackageName)
}

func (fu *fileUpdatePipeline) findTargetStructure() {
	fu.structure = fu.src.FindStructure(fu.TargetStructureNamePattern)
}

func (fu *fileUpdatePipeline) enrichStructure() {
	fu.strategy.enrichStructure(fu.structure)
}

func (fu *fileUpdatePipeline) findTargetFunction() {
	fu.function = fu.src.FindFunction(fu.TargetFunctionNamePattern)
	fu.functionBody = fu.function.GetBody()
}

func (fu *fileUpdatePipeline) enrichFunction() {
	fu.findStructureInitExpression()
	fu.functionBody.Wipe()
	fu.functionBody.AppendNewAssigment("a", token.DEFINE, fu.structureInitExpression)
	fu.strategy.initCustomField(fu.functionBody)
	newReturn := pretty_dst.NewIdent("a", nil)
	fu.functionBody.AppendNewReturn(newReturn)
}

func (fu *fileUpdatePipeline) findStructureInitExpression() {
	ret := fu.functionBody.GetFirstReturn()
	if id := ret.GetIdentifier(); id != nil {
		assigment := fu.functionBody.GetFirstAssignment()
		fu.structureInitExpression = assigment.Rhs[0].(*dst.UnaryExpr)
	} else {
		ret.GetUnaryExpression().GetCompositeLiteral().RemoveElementByKey("subAggregations")
		fu.structureInitExpression = ret.GetUnaryExpression().UnaryExpr
	}
}

//pick modification stratedy. Depends on target structure fieldset
func (fu *fileUpdatePipeline) pickStrategy() {
	if fu.structure.IsFieldExists("subAggregations") {
		fu.strategy = injectableStrategy{}
	} else {
		fu.strategy = notInjectableStrategy{}
	}
}

type modificationStrategy interface {
	enrichStructure(*pretty_dst.StructureDeclaration)
	initCustomField(*pretty_dst.FunctionBody)
	name() string
}

type injectableStrategy struct {
	modificationStrategy
}

func (injectableStrategy) enrichStructure(structure *pretty_dst.StructureDeclaration) {
	structure.AddField("*Injectable", "")
	_ = structure.RemoveField("subAggregations")
}

func (injectableStrategy) initCustomField(body *pretty_dst.FunctionBody) {
	body.AppendNewAssigment("a.Injectable",
		token.ASSIGN,
		pretty_dst.NewCallExpression("newInjectable", "a"),
	)
}

func (injectableStrategy) name() string{
	return "Injectable"
}

type notInjectableStrategy struct {
	modificationStrategy
}

func (notInjectableStrategy) enrichStructure(structure *pretty_dst.StructureDeclaration) {
	structure.AddField("*NotInjectable", "")
}

func (notInjectableStrategy) initCustomField(body *pretty_dst.FunctionBody) {
	body.AppendNewAssigment("a.NotInjectable",
		token.ASSIGN,
		pretty_dst.NewCallExpression("newNotInjectable", "a"),
	)
}

func (notInjectableStrategy) name() string{
	return "NotInjectable"
}