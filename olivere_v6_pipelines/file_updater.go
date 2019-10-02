package olivere_v6_pipelines

import (
	"github.com/dave/dst"
	"github.com/konovenschi/aggretastic-sync/errors"
	"github.com/konovenschi/aggretastic-sync/pretty_dst"
	"go/token"
	"gopkg.in/src-d/go-billy.v4"
	"strings"
)

type fileUpdatePipeline struct {
	Filename string
	FS       billy.Filesystem
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
func (fu *fileUpdatePipeline) Run() {
	fu.parseFile()
	defer fu.saveFile()

	fu.renamePackage()
	fu.findTargetStructure()
	if fu.structure == nil {
		return
	}
	fu.pickStrategy()
	fu.enrichStructure()

	fu.findTargetFunction()
	if fu.function == nil {
		return
	}
	fu.enrichFunction()
}

//parse ast from file
func (fu *fileUpdatePipeline) parseFile() {
	file, err := fu.FS.Open(fu.Filename)
	errors.PanicOnError(errCantOpenFile, err)

	fu.src = pretty_dst.NewDst(file)
}

//save ast changes on disk without search_ prefix
func (fu *fileUpdatePipeline) saveFile() {
	filename := strings.Replace(fu.Filename, "search_aggs_", "aggs_", -1)
	file, err := fu.FS.Create(filename)
	errors.PanicOnError(errCantOpenFile, err)
	defer func() {
		err := file.Close()
		errors.PanicOnError(errCantCloseFile, err)
	}()

	err = fu.src.Save(file)
	errors.PanicOnError(errCantWriteFile, err)

	if filename != fu.Filename {
		err = fu.FS.Remove(fu.Filename)
		errors.PanicOnError(errCantRemoveFile, err)
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

func (injectableStrategy) name() string {
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

func (notInjectableStrategy) name() string {
	return "NotInjectable"
}
