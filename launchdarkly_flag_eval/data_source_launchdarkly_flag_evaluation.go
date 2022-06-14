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
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			flagKey: {
				Type:     types.StringType,
				Required: true,
			},
			flagType: {
				Type:     types.StringType,
				Computed: true,
			},
			defaultValue: {
				Type:     types.BoolType, // TODO refactor to pass type via wrapper function
				Required: true,
			},
			value: {
				Type:     types.BoolType, // TODO refactor to pass type via wrapper function
				Computed: true,
			},
			userContext: {
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					keyAttribute: {
						Type:     types.StringType,
						Required: true,
					},
					secondaryKeyAttribute: {
						Type:     types.StringType,
						Optional: true,
					},
					ipAttribute: {
						Type:     types.StringType,
						Optional: true,
					},
					countryAttribute: {
						Type:     types.StringType,
						Optional: true,
					},
					emailAttribute: {
						Type:     types.StringType,
						Optional: true,
					},
					firstNameAttribute: {
						Type:     types.StringType,
						Optional: true,
					},
					lastNameAttribute: {
						Type:     types.StringType,
						Optional: true,
					},
					avatarAttribute: {
						Type:     types.StringType,
						Optional: true,
					},
					nameAttribute: {
						Type:     types.StringType,
						Optional: true,
					},
					anonymousAttribute: {
						Type:     types.BoolType,
						Optional: true,
					},
					// customAttributes: {
					// 	Optional:   true,
					// 	// TODO: should this also be a map (can we pass json in custom attributes?)
					// 	Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{}, tfsdk.MapNestedAttributesOptions{}),
					// },
					// TODO private attributes??
				}),
			},
		},
	}, nil
}

func (d dataSourceFlagEvaluationBooleanType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceFlagEvaluationBoolean{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceFlagEvaluationBoolean struct {
	p provider
}

func (d dataSourceFlagEvaluationBoolean) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
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
		//Custom    ldvalue.ValueMap
	}

	var dataSourceState struct {
		FlagKey      types.String `tfsdk:"flag_key"`
		FlagType     types.String `tfsdk:"flag_type"`
		DefaultValue types.Bool   `tfsdk:"default_value"`
		Value        types.Bool   `tfsdk:"value"`
		UserContext  lduser.User  `tfsdk:"context"`
	}

	tflog.Info(ctx, "test")

	diags := req.Config.Get(ctx, &dataSourceState)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("%+v", dataSourceState))

	// set state
	diags = resp.State.Set(ctx, &dataSourceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// evaluation, err := r.p.client.BoolVariation(flagKey, userCtx, defaultValue)
}

// func dataSourceFlagEvaluationReadWrapper(typ schema.ValueType) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 		tflog.Info(ctx, "ENTERING FUNCTION")
// 		var diags diag.Diagnostics
// 		client := meta.(*ld.LDClient)

// 		flagKey := d.Get(FLAG_KEY).(string)
// 		rawContext := d.Get(CONTEXT).([]interface{})
// 		// TODO construct user object properly
// 		rawContextMap := rawContext[0].(map[string]interface{})
// 		_ = rawContextMap
// 		//userCtxBuilder := lduser.NewUserBuilder(rawContextMap[keyAttribute])
// 		userCtxBuilder := lduser.NewUserBuilder("hello-world")
// 		//userCtxBuilder.Name(rawContextMap[nameAttribute])
// 		userCtx := userCtxBuilder.Build()
// 		// todo rest of userCtx

// 		switch typ {
// 		case types.StringType:
// 			d.Set(FLAG_TYPE, "string")
// 			defaultValue := d.Get(DEFAULT_VALUE).(string)
// 			value, err := client.StringVariation(flagKey, userCtx, defaultValue)
// 			if err != nil {
// 				return diag.FromErr(err)
// 			}
// 			err = d.Set(VALUE, value)
// 			if err != nil {
// 				return diag.FromErr(err)
// 			}
// 			d.Set(DEFAULT_VALUE, defaultValue)
// 		case schema.TypeBool:
// 			d.Set(FLAG_TYPE, "boolean")
// 			defaultValue := d.Get(DEFAULT_VALUE).(bool)
// 			value, err := client.BoolVariation(flagKey, userCtx, defaultValue)
// 			if err != nil {
// 				return diag.FromErr(err)
// 			}
// 			err = d.Set(VALUE, value)
// 			if err != nil {
// 				return diag.FromErr(err)
// 			}
// 			d.Set(DEFAULT_VALUE, defaultValue)
// 		case schema.TypeInt:
// 			d.Set(FLAG_TYPE, "int")
// 			defaultValue := d.Get(DEFAULT_VALUE).(int)
// 			value, err := client.IntVariation(flagKey, userCtx, defaultValue)
// 			if err != nil {
// 				return diag.FromErr(err)
// 			}
// 			err = d.Set(VALUE, value)
// 			if err != nil {
// 				return diag.FromErr(err)
// 			}
// 			d.Set(DEFAULT_VALUE, defaultValue)
// 		case schema.TypeFloat:
// 			d.Set(FLAG_TYPE, "float")
// 			defaultValue := d.Get(DEFAULT_VALUE).(float64)
// 			value, err := client.Float64Variation(flagKey, userCtx, defaultValue)
// 			if err != nil {
// 				return diag.FromErr(err)
// 			}
// 			err = d.Set(VALUE, value)
// 			if err != nil {
// 				return diag.FromErr(err)
// 			}
// 			d.Set(DEFAULT_VALUE, defaultValue)
// 			// case schema.TypeMap:
// 			// 	var jsonRaw json.RawMessage
// 			// 	err := jsonRaw.UnmarshalJSON([]byte(rawDefault))
// 			// 	if err != nil {
// 			// 		return diag.FromErr(err)
// 			// 	}

// 			// 	defaultValue := ldvalue.Raw(jsonRaw)
// 			// 	value, err := client.JSONVariation(flagKey, userCtx, defaultValue)
// 			// 	if err != nil {
// 			// 		return diag.FromErr(err)
// 			// 	}
// 		}

// 		d.Set(FLAG_KEY, flagKey)
// 		// TODO we need helper functions to convert back and forth
// 		d.Set(CONTEXT, rawContextMap)
// 		d.SetId(flagKey)

// 		return diags
// 	}
// }
