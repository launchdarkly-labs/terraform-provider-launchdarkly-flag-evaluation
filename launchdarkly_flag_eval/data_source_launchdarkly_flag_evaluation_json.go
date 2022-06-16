package launchdarkly_flag_eval

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

type dataSourceFlagEvaluationJSONType struct{}

func (d dataSourceFlagEvaluationJSONType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return getFlagEvaluationSchemaForType(types.ObjectType{})
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

	fmt.Println("REQ: ", req)
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

	fmt.Println("DEFAULT ATTRS: ", dataSourceState.DefaultValue.Attrs)

	var ldMap ldvalue.ValueMapBuilder
	for key, val := range dataSourceState.DefaultValue.Attrs {
		fmt.Println("KV: ", key, val)
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

	attrTypes, attrs, err := convertJSONEvaluationToTFAttrs(evaluation)
	if err != nil {
		resp.Diagnostics.AddError(
			"Flag evaluation failed",
			"Could not evaluate flag: "+err.Error(),
		)
		return
	}
	dataSourceState.Value = types.Object{
		Unknown:   false,
		Null:      false,
		Attrs:     attrs,
		AttrTypes: attrTypes,
	}

	// set state
	diags = resp.State.Set(ctx, &dataSourceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func convertJSONEvaluationToTFAttrs(evaluation ldvalue.Value) (map[string]attr.Type, map[string]attr.Value, error) {
	var attrType map[string]attr.Type
	var attrValue map[string]attr.Value
	for k, v := range evaluation.AsValueMap().AsMap() {
		switch v.Type() {
		case ldvalue.BoolType:
			attrType[k] = types.BoolType
			attrValue[k] = types.Bool{Value: v.BoolValue()}
		case ldvalue.StringType:
			attrType[k] = types.StringType
			attrValue[k] = types.String{Value: v.StringValue()}
		case ldvalue.NumberType:
			attrType[k] = types.Int64Type
			attrValue[k] = types.Int64{Value: int64(v.IntValue())}
		default:
			return attrType, attrValue, errors.New("unknown value in evaluation map")
		}
	}
	return attrType, attrValue, nil
}
