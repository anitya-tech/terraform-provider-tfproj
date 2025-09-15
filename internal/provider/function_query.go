package provider

import (
	"context"
	"fmt"

	filequery "github.com/anitya-tech/terraform-provider-tfproj/internal/utils/file_query"
	"github.com/anitya-tech/terraform-provider-tfproj/internal/utils/tf_project"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &QueryFunction{}

type QueryFunction struct{}

func NewQueryFunction() function.Function {
	return &QueryFunction{}
}

func (f *QueryFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "query"
}

func (f *QueryFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "query data from file",
		Description: "query data from file",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "file",
				Description: "file to query data from",
			},
			function.StringParameter{
				AllowNullValue: true,
				Name:           "pattern",
				Description:    "pattern to query data, like .foo.bar",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *QueryFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var template string
	var pattern *string
	resp.Error = req.Arguments.Get(ctx, &template, &pattern)
	if resp.Error != nil {
		return
	}

	cur_project, err := tf_project.AnalyzeCurrentProject()
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("Error finding project root: %s", err.Error()))
		return
	}

	filename := cur_project.Format(template)

	qsrc, err := filequery.Parse(filename)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("Error parsing file: %s", err.Error()))
		return
	}

	result := qsrc.Query(pattern)
	if _, ok := result.(string); !ok {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("query result is not a string: %T", result))
		return
	}

	resp.Error = resp.Result.Set(ctx, &result)
}
