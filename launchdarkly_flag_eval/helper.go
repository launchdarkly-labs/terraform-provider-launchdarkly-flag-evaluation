package launchdarkly_flag_eval

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

func getFlagEvaluationSchemaForType(typ attr.Type) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			value: {
				Type:     typ,
				Computed: true,
			},
			flagKey: {
				Type:     types.StringType,
				Required: true,
			},
			defaultValue: {
				Type:     typ,
				Required: true,
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
					customAttributes: {
						Type:     DynamicType{},
						Optional: true,
					},
				}),
			},
		},
	}, nil
}

func convertUserContextToLDUserContext(userKey string, userContext LDUser) (ldUserContext lduser.User, isUnknown bool) {
	// type LDUser struct {
	// 	Key       types.String `tfsdk:"key"`
	// 	Secondary types.String `tfsdk:"secondary"`
	// 	IP        types.String `tfsdk:"ip"`
	// 	Country   types.String `tfsdk:"country"`
	// 	Email     types.String `tfsdk:"email"`
	// 	FirstName types.String `tfsdk:"first_name"`
	// 	LastName  types.String `tfsdk:"last_name"`
	// 	Avatar    types.String `tfsdk:"avatar"`
	// 	Name      types.String `tfsdk:"name"`
	// 	Anonymous types.Bool   `tfsdk:"anonymous"`
	// 	Custom    Dynamic      `tfsdk:"custom"`
	// }
	for key, val := range userContext.Custom.Values {
		tflog.Info(context.Background(), fmt.Sprintf("Got %s with value %s", key, val))
	}

	builder := lduser.NewUserBuilder(userKey)

	if userContext.Key.Unknown {
		return nil, true
	}
	builder.Key(userContext.Key.Value)

	if userContext.Secondary.Unknown {
		return nil, true
	}
	if !userContext.Secondary.Null {
		builder.Secondary(userContext.Secondary.Value)
	}

	if userContext.IP.Unknown {
		return nil, true
	}
	if !userContext.IP.Null { // TODO handle unknown value
		builder.IP(userContext.IP.Value)
	}

	if userContext.Country.Unknown {
		return nil, true
	}
	if !userContext.Country.Null { // TODO handle unknown value
		builder.Country(userContext.Country.Value)
	}

	if userContext.Email.Unknown {
		return nil, true
	}
	if !userContext.Email.Null { // TODO handle unknown value
		builder.Email(userContext.Email.Value)
	}

	if userContext.FirstName.Unknown {
		return nil, true
	}
	if !userContext.FirstName.Null { // TODO handle unknown value
		builder.FirstName(userContext.FirstName.Value)
	}

	if userContext.LastName.Unknown {
		return nil, true
	}
	if !userContext.LastName.Null { // TODO handle unknown value
		builder.LastName(userContext.LastName.Value)
	}

	if userContext.Avatar.Unknown {
		return nil, true
	}
	if !userContext.Avatar.Null { // TODO handle unknown value
		builder.Avatar(userContext.Avatar.Value)
	}

	if userContext.Name.Unknown {
		return nil, true
	}
	if !userContext.Name.Null { // TODO handle unknown value
		builder.Name(userContext.Name.Value)
	}

	if userContext.Anonymous.Unknown {
		return nil, true
	}
	if !userContext.Anonymous.Null { // TODO handle unknown value
		builder.Anonymous(userContext.Anonymous.Value)
	}

	for key, val := range userContext.Custom.Values {
		var ldval ldvalue.Value

		if !val.IsFullyKnown() {
			return nil, true
		}

		switch val.Type() {
			case 
		}

		builder.Custom(key, ldval)
	}

	return builder.Build(), false
}
