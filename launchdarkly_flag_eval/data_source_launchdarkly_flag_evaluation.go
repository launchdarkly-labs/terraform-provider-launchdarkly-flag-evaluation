package launchdarkly_flag_eval

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
)

const (
	flagKey      = "flag_key"
	flagType     = "flag_type"
	userContext  = "context"
	variation    = "variation_type"
	value        = "value"
	defaultValue = "default_value"

	// KeyAttribute is the standard attribute name corresponding to User.GetKey().
	keyAttribute = "key"
	// SecondaryKeyAttribute is the standard attribute name corresponding to User.GetSecondaryKey().
	secondaryKeyAttribute = "secondary"
	// IPAttribute is the standard attribute name corresponding to User.GetIP().
	ipAttribute = "ip"
	// CountryAttribute is the standard attribute name corresponding to User.GetCountry().
	countryAttribute = "country"
	// EmailAttribute is the standard attribute name corresponding to User.GetEmail().
	emailAttribute = "email"
	// FirstNameAttribute is the standard attribute name corresponding to User.GetFirstName().
	firstNameAttribute = "first_name"
	// LastNameAttribute is the standard attribute name corresponding to User.GetLastName().
	lastNameAttribute = "last_name"
	// AvatarAttribute is the standard attribute name corresponding to User.GetAvatar().
	avatarAttribute = "avatar"
	// NameAttribute is the standard attribute name corresponding to User.GetName().
	nameAttribute = "name"
	// AnonymousAttribute is the standard attribute name corresponding to User.GetAnonymous().
	anonymousAttribute = "anonymous"
	customAttributes   = "custom"
)

type dataSourceFlagEvaluationBooleanType struct {
	p provider
}

func (r dataSourceFlagEvaluationBooleanType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return getFlagEvaluationSchemaForType(types.BoolType)
}

func (d dataSourceFlagEvaluationBooleanType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceFlagEvaluationBoolean{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceFlagEvaluationBoolean struct {
	p provider
}

type LDUser struct {
	Key       types.String `tfsdk:"key"`
	Secondary types.String `tfsdk:"secondary"`
	Ip        types.String `tfsdk:"ip"`
	Country   types.String `tfsdk:"country"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Avatar    types.String `tfsdk:"avatar"`
	Name      types.String `tfsdk:"name"`
	Anonymous types.Bool   `tfsdk:"anonymous"`
	Custom    Dynamic      `tfsdk:"custom"`
}

func (d dataSourceFlagEvaluationBoolean) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var dataSourceState struct {
		FlagKey      types.String `tfsdk:"flag_key"`
		DefaultValue types.Bool   `tfsdk:"default_value"`
		Value        types.Bool   `tfsdk:"value"`
		UserContext  LDUser       `tfsdk:"context"`
	}

	tflog.Info(ctx, "test")

	diags := req.Config.Get(ctx, &dataSourceState)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("STATE %+v", dataSourceState))

	userCtx := convertUserContextToLDUserContext(dataSourceState.UserContext.Key.Value, dataSourceState.UserContext)
	evaluation, err := d.p.client.BoolVariation(dataSourceState.FlagKey.Value, userCtx, dataSourceState.DefaultValue.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Flag evaluation failed",
			"Could not evaluate flag: "+err.Error(),
		)
		return
	}
	d.p.client.Flush()

	dataSourceState.Value = types.Bool{
		Unknown: false,
		Null:    false,
		Value:   evaluation,
	}

	// set state
	diags = resp.State.Set(ctx, &dataSourceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func convertUserContextToLDUserContext(userKey string, userContext LDUser) lduser.User {
	for key, val := range userContext.Custom.Values {
		tflog.Info(context.Background(), fmt.Sprintf("Got %s with value %s", key, val))
	}

	// builder := lduser.NewUserBuilder(userKey)
	return lduser.User{}
}
