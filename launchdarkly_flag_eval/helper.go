package launchdarkly_flag_eval

import (
	"context"
	"fmt"
	"math/big"

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

func convertUserContextToLDUserContext(ctx context.Context, userKey string, userContext LDUser, diags diag.Diagnostics) (ldUserContext lduser.User, isUnknown bool) {
	// for key, val := range userContext.Custom.Values {
	// 	tflog.Info(ctx, fmt.Sprintf("Got %s with value %s", key, val))
	// }

	builder := lduser.NewUserBuilder(userKey)
	if userContext.Key.Unknown {
		tflog.Info(ctx, "Key is unknown\n")
		return lduser.User{}, true
	}
	builder.Key(userContext.Key.Value)

	if userContext.Secondary.Unknown {
		tflog.Info(ctx, "Secondary is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.Secondary.Null {
		builder.Secondary(userContext.Secondary.Value)
	}

	if userContext.IP.Unknown {
		tflog.Info(ctx, "IP is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.IP.Null {
		builder.IP(userContext.IP.Value)
	}

	if userContext.Country.Unknown {
		tflog.Info(ctx, "Country is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.Country.Null {
		builder.Country(userContext.Country.Value)
	}

	if userContext.Email.Unknown {
		tflog.Info(ctx, "Email is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.Email.Null {
		builder.Email(userContext.Email.Value)
	}

	if userContext.FirstName.Unknown {
		tflog.Info(ctx, "FirstName is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.FirstName.Null {
		builder.FirstName(userContext.FirstName.Value)
	}

	if userContext.LastName.Unknown {
		tflog.Info(ctx, "LastName is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.LastName.Null {
		builder.LastName(userContext.LastName.Value)
	}

	if userContext.Avatar.Unknown {
		tflog.Info(ctx, "Avatar is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.Avatar.Null {
		builder.Avatar(userContext.Avatar.Value)
	}

	if userContext.Name.Unknown {
		tflog.Info(ctx, "Name is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.Name.Null {
		builder.Name(userContext.Name.Value)
	}

	if userContext.Anonymous.Unknown {
		tflog.Info(ctx, "Anonymous is unknown\n")
		return lduser.User{}, true
	}
	if !userContext.Anonymous.Null {
		builder.Anonymous(userContext.Anonymous.Value)
	}

	for key, val := range userContext.Custom.Values {
		ldval, isUnknown := convert(ctx, key, val, diags)
		if isUnknown {
			return lduser.User{}, true
		}
		builder.Custom(key, ldval)
	}

	lduser := builder.Build()
	tflog.Info(ctx, fmt.Sprintf("%+v", lduser))

	return lduser, false
}

func convert(ctx context.Context, key string, val tftypes.Value, diags diag.Diagnostics) (value ldvalue.Value, isUnknown bool) {
	tflog.Info(ctx, fmt.Sprintf("THESE ARE SOME VALUES: %s = %v", key, val.IsFullyKnown()))
	if !val.IsFullyKnown() {
		return ldvalue.Value{}, true
	}

	switch {
	case val.Type().Is(tftypes.Bool):
		tflog.Info(ctx, "THIS IS A BOOL TYPE")
		var v bool
		err := val.As(&v)
		if err != nil {
			diags.AddAttributeError(nil, "Invalid type", "Can not convert value to boolean")
			return ldvalue.Value{}, true
		}
		return ldvalue.Bool(v), false
	case val.Type().Is(tftypes.String):
		tflog.Info(ctx, "THIS IS A STRING TYPE")
		var v string
		err := val.As(&v)
		if err != nil {
			diags.AddAttributeError(nil, "Invalid type", "Can not convert value to string")
			return ldvalue.Value{}, true
		}
		return ldvalue.String(v), false
	case val.Type().Is(tftypes.Number):
		tflog.Info(ctx, "THIS IS A NUMBER TYPE")
		// test := val.ToTerraformValue()
		var vf64 *big.Float
		err := val.As(&vf64)
		if err != nil {
			tflog.Info(ctx, fmt.Sprintf("failed to convert %v to int", val))
			diags.AddAttributeError(nil, "Invalid type", "Can not convert value to big.float")
			return ldvalue.Value{}, true
		}
		tflog.Info(ctx, fmt.Sprintf("big.float is: %v", vf64))

		if vf64.IsInt() {
			tflog.Info(ctx, "THIS IS AN INT WITHIN A NUMBER TYPE")
			f, accuracy := vf64.Int64()
			_ = accuracy

			// builder.Custom(key, ldvalue.Int(int(f)))
			return ldvalue.Int(int(f)), false
		}

		f, accuracy := vf64.Float64()
		_ = accuracy

		// builder.Custom(key, ldvalue.Float64(f))
		return ldvalue.Float64(f), false

	// case val.Type().Is(tftypes.Object{}): // tftypes.Object.Is(val.Type()):
	// 	var obj map[string]tftypes.Value

	// 	err := val.As(&obj)
	// 	if err != nil {
	// 		diags.AddAttributeError(nil, "Invalid type", "Can not convert value to map")
	// 		return ldvalue.Value{}, true
	// 	}

	// 	ldvalBuilder := ldvalue.ObjectBuildWithCapacity(len(obj))
	// 	for k, v := range obj {
	// 		newldval, isUnknown := convert(ctx, k, v, diags)
	// 		if isUnknown {
	// 			return ldvalue.Value{}, true
	// 		}

	// 		ldvalBuilder.Set(k, newldval)
	// 	}

	// 	return ldvalBuilder.Build(), false

	default:
		// todo object/array
		tflog.Info(ctx, fmt.Sprintf("THIS IS A VALUE: %+v", val))
	}
	return ldvalue.Value{}, false
}
