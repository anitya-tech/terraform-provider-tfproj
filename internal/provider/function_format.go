package provider

import (
	"context"
	"fmt"

	"github.com/anitya-tech/terraform-provider-tfproj/internal/utils/tf_project"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &FormatFunction{}

type FormatFunction struct{}

func NewFormatFunction() function.Function {
	return &FormatFunction{}
}

func (f *FormatFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "format"
}

func (f *FormatFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "format provided template",
		Description: "fill provided template with project information",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "template",
				Description: "template to fill with project information",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *FormatFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var template string
	resp.Error = req.Arguments.Get(ctx, &template)
	if resp.Error != nil {
		return
	}

	cur_project, err := tf_project.AnalyzeCurrentProject()
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("Error finding project root: %s", err.Error()))
		return
	}

	result := cur_project.Format(template)

	resp.Error = resp.Result.Set(ctx, &result)
}
