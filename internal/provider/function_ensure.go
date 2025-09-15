package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/anitya-tech/terraform-provider-tfproj/internal/utils/tf_project"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &EnsureFunction{}

type EnsureFunction struct{}

func NewEnsureFunction() function.Function {
	return &EnsureFunction{}
}

func (f *EnsureFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "ensure"
}

func (f *EnsureFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "ensure provided file exists",
		Description: "template will be formatted like 'format' function",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "file",
				Description: "file to ensure exists",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *EnsureFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
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

	_, err = os.Stat(result)
	if os.IsNotExist(err) {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("Resource %s not found", result))
		return
	} else if err != nil {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("Error checking resource %s: %s", result, err.Error()))
		return
	}

	resp.Error = resp.Result.Set(ctx, &result)
}
