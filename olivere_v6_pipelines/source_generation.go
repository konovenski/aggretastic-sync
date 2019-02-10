package olivere_v6_pipelines

import "github.com/dave/jennifer/jen"

type sourceGenerationPipeline struct {
	SyncNodes []*syncNode
}

func (sg *sourceGenerationPipeline) run () {
	f := jen.NewFile("aggretastic")

	imps := jen.Dict{}
	notImps := jen.Dict{}
	for _, node := range sg.SyncNodes{
		switch node.aggregationType{
		case "Injectable":
			imps[jen.Lit(node.functionName)] = jen.Id(node.functionName)
		case "NotInjectable":
			notImps[jen.Lit(node.functionName)] = jen.Id(node.functionName)
		}
	}

	f.Var().Id("aggMap").Op("=").Map(jen.String()).Map(jen.String()).Interface().Values(jen.Dict{
		jen.Lit("Injectable"): jen.Map(jen.String()).Interface().Values(imps),
		jen.Lit("NotInjectable"): jen.Map(jen.String()).Interface().Values(notImps),
	})
	_ = f.Save("generated-aggregations-mapping.go")

}
