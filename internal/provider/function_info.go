package provider

import (
	"context"
	"fmt"

	"github.com/anitya-tech/terraform-provider-tfproj/internal/utils/tf_project"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var moduleInfoType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"path": types.StringType,
		"name": types.StringType,
	},
}

var projectInfoType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"path": types.StringType,
		"modules": types.MapType{
			ElemType: moduleInfoType,
		},
		"current_module": moduleInfoType,
	},
}

var _ function.Function = &InfoFunction{}

type InfoFunction struct{}

func NewInfoFunction() function.Function {
	return &InfoFunction{}
}

func (f *InfoFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "info"
}

func (f *InfoFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "get project information",
		Description: "anaylyze whole git project and find all modules",

		Parameters: []function.Parameter{},
		Return: function.ObjectReturn{
			AttributeTypes: projectInfoType.AttrTypes,
		},
	}
}

func (f *InfoFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	curr_project, err := tf_project.AnalyzeCurrentProject()
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("Error analyzing current project: %s", err.Error()))
		return
	}

	projectInfo, diags := types.ObjectValueFrom(ctx, projectInfoType.AttrTypes, curr_project)

	resp.Error = function.FuncErrorFromDiags(ctx, diags)
	if resp.Error != nil {
		return
	}

	resp.Error = resp.Result.Set(ctx, &projectInfo)
}
