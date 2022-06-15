package launchdarkly_flag_eval

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

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

func convertUserContextToLDUserContext(userKey string, userContext LDUser, diags diag.Diagnostics) (ldUserContext lduser.User, isUnknown bool) {
	for key, val := range userContext.Custom.Values {
		tflog.Info(context.Background(), fmt.Sprintf("Got %s with value %s", key, val))
	}

	tflog.Info(context.Background(), "what")
	builder := lduser.NewUserBuilder(userKey)
	if userContext.Key.Unknown {
		tflog.Info(context.Background(), "Key is unknown\n")
		return lduser.User{}, true
	}
	builder.Key(userContext.Key.Value)

	if userContext.Secondary.Unknown {
		tflog.Info(context.Background(), "Secondary is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.Secondary.Null {
		builder.Secondary(userContext.Secondary.Value)
	}

	if userContext.IP.Unknown {
		tflog.Info(context.Background(), "IP is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.IP.Null {
		builder.IP(userContext.IP.Value)
	}

	if userContext.Country.Unknown {
		tflog.Info(context.Background(), "Country is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.Country.Null {
		builder.Country(userContext.Country.Value)
	}

	if userContext.Email.Unknown {
		tflog.Info(context.Background(), "Email is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.Email.Null {
		builder.Email(userContext.Email.Value)
	}

	if userContext.FirstName.Unknown {
		tflog.Info(context.Background(), "FirstName is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.FirstName.Null {
		builder.FirstName(userContext.FirstName.Value)
	}

	if userContext.LastName.Unknown {
		tflog.Info(context.Background(), "LastName is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.LastName.Null {
		builder.LastName(userContext.LastName.Value)
	}

	if userContext.Avatar.Unknown {
		tflog.Info(context.Background(), "Avatar is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.Avatar.Null {
		builder.Avatar(userContext.Avatar.Value)
	}

	if userContext.Name.Unknown {
		tflog.Info(context.Background(), "Name is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.Name.Null {
		builder.Name(userContext.Name.Value)
	}

	if userContext.Anonymous.Unknown {
		tflog.Info(context.Background(), "Anonymous is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.Anonymous.Null {
		builder.Anonymous(userContext.Anonymous.Value)
	}

	for key, val := range userContext.Custom.Values {
		tflog.Info(context.Background(), fmt.Sprintf("hello %s: %t", key, val.IsFullyKnown()))
		if !val.IsFullyKnown() {
			return lduser.User{}, true
		}

		switch {
		case tftypes.Bool.Is(val.Type()):
			var v bool
			err := val.As(&v)
			if err != nil {
				diags.AddAttributeError(nil, "Invalid type", "Can not convert value to boolean")
				return lduser.User{}, true
			}
			builder.Custom(key, ldvalue.Bool(v))
		case tftypes.String.Is(val.Type()):
			var v string
			err := val.As(&v)
			if err != nil {
				diags.AddAttributeError(nil, "Invalid type", "Can not convert value to string")
				return lduser.User{}, true
			}
			builder.Custom(key, ldvalue.String(v))
		case tftypes.Number.Is(val.Type()):
			var vf64 float64
			err := val.As(&vf64)
			if err == nil {
				builder.Custom(key, ldvalue.Float64(vf64))
				continue
			}

			var vint int
			err = val.As(&vint)
			if err != nil {
				diags.AddAttributeError(nil, "Invalid type", "Can not convert value to int or float64")
				return lduser.User{}, true
			}
			builder.Custom(key, ldvalue.Int(vint))
		default:
			// todo object/array
			tflog.Info(context.Background(), fmt.Sprintf("val: %+v", val))
		}
	}

	return builder.Build(), false
}
