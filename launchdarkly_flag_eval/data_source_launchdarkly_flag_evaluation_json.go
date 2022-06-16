package launchdarkly_flag_eval

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

type dataSourceFlagEvaluationJSONType struct{}

func (d dataSourceFlagEvaluationJSONType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return getFlagEvaluationSchemaForType(types.Int64Type)
}

func (d dataSourceFlagEvaluationJSONType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceFlagEvaluationJSON{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceFlagEvaluationJSON struct {
	p provider
}

func (d dataSourceFlagEvaluationJSON) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var dataSourceState struct {
		FlagKey      types.String `tfsdk:"flag_key"`
		DefaultValue types.Object `tfsdk:"default_value"`
		Value        types.Object `tfsdk:"value"`
		UserContext  LDUser       `tfsdk:"context"`
	}

	diags := req.Config.Get(ctx, &dataSourceState)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var ldMap ldvalue.ValueMapBuilder
	for key, val := range dataSourceState.DefaultValue.Attrs {
		switch {
		case val.Type(ctx) == types.BoolType:
			var temp bool
			tfVal, err := val.ToTerraformValue(ctx)
			if err != nil {
				tflog.Debug(ctx, err.Error())
				return
			}
			_ = tfVal.As(&temp)
			ldMap.Set(key, ldvalue.Bool(temp))
		case val.Type(ctx) == types.StringType:
			var temp string
			tfVal, err := val.ToTerraformValue(ctx)
			if err != nil {
				tflog.Debug(ctx, err.Error())
				return
			}
			_ = tfVal.As(&temp)
			ldMap.Set(key, ldvalue.String(temp))
		case val.Type(ctx) == types.Int64Type:
			var temp int
			tfVal, err := val.ToTerraformValue(ctx)
			if err != nil {
				tflog.Debug(ctx, err.Error())
				return
			}
			_ = tfVal.As(&temp)
			ldMap.Set(key, ldvalue.Int(temp))
		//case val.Type(ctx) == types.Map:
		//	var temp map[string]
		//	tfVal, err := val.ToTerraformValue(ctx)
		//	if err != nil {
		//		tflog.Debug(ctx, err.Error())
		//		return
		//	}
		//	_ = tfVal.As(&temp)
		//	ldMap.Set(key, ldvalue.(temp))
		default:
			resp.Diagnostics.AddError(
				"Flag evaluation failed",
				"Unknown value in object type",
			)
			return
		}
	}

	userCtx, _ := convertUserContextToLDUserContext(ctx, dataSourceState.UserContext.Key.Value, dataSourceState.UserContext, resp.Diagnostics)
	evaluation, err := d.p.client.JSONVariation(dataSourceState.FlagKey.Value, userCtx, ldMap.Build().AsValue())
	if err != nil {
		resp.Diagnostics.AddError(
			"Flag evaluation failed",
			"Could not evaluate flag: "+err.Error(),
		)
		return
	}
	d.p.client.Flush()

	dataSourceState.Value = types.Object{
		Unknown: false,
		Null:    false,
		Attrs:   evaluation,
		AttrTypes:,
	}

	// set state
	diags = resp.State.Set(ctx, &dataSourceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
