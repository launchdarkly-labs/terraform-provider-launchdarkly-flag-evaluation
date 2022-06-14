package launchdarkly_flag_eval

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
)

const (
	FLAG_KEY       = "flag_key"
	FLAG_TYPE      = "flag_type"
	CONTEXT        = "context"
	VARIATION_TYPE = "variation_type"
	VALUE          = "value"
	DEFAULT_VALUE  = "default_value"

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

func dataSourceFlagEvaluation(typ schema.ValueType) *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFlagEvaluationReadWrapper(typ),

		// TODO: see what we actually need
		Schema: map[string]*schema.Schema{
			FLAG_KEY: {
				Type:     schema.TypeString,
				Required: true,
			},
			FLAG_TYPE: {
				Type:             schema.TypeString,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"boolean", "string", "float", "int"}, false)),
			},
			DEFAULT_VALUE: {
				Type:     typ,
				Required: true,
			},
			VALUE: {
				Type:     typ,
				Computed: true,
			},
			// TODO: figure out the best name for this
			CONTEXT: {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						keyAttribute: {
							Type:     schema.TypeString,
							Required: true,
							// todo 64 = CUSTOM_PROPERTY_CHAR_LIMIT
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 64)),
						},
						secondaryKeyAttribute: {
							Type:     schema.TypeString,
							Optional: true,
						},
						ipAttribute: {
							Type:     schema.TypeString,
							Optional: true,
						},
						countryAttribute: {
							Type:     schema.TypeString,
							Optional: true,
						},
						emailAttribute: {
							Type:     schema.TypeString,
							Optional: true,
						},
						firstNameAttribute: {
							Type:     schema.TypeString,
							Optional: true,
						},
						lastNameAttribute: {
							Type:     schema.TypeString,
							Optional: true,
						},
						avatarAttribute: {
							Type:     schema.TypeString,
							Optional: true,
						},
						nameAttribute: {
							Type:     schema.TypeString,
							Optional: true,
						},
						anonymousAttribute: {
							Type:     schema.TypeBool,
							Optional: true,
						},
						// customAttributes:
						// todo custom
						// todo privateAttributes
					},
				},
			},
		},
	}
}

func dataSourceFlagEvaluationReadWrapper(typ schema.ValueType) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		tflog.Info(ctx, "ENTERING FUNCTION")
		var diags diag.Diagnostics
		client := meta.(*ld.LDClient)

		flagKey := d.Get(FLAG_KEY).(string)
		rawContext := d.Get(CONTEXT).([]interface{})
		// TODO construct user object properly
		rawContextMap := rawContext[0].(map[string]interface{})
		_ = rawContextMap
		//userCtxBuilder := lduser.NewUserBuilder(rawContextMap[keyAttribute])
		userCtxBuilder := lduser.NewUserBuilder("hello-world")
		//userCtxBuilder.Name(rawContextMap[nameAttribute])
		userCtx := userCtxBuilder.Build()
		// todo rest of userCtx

		switch typ {
		case schema.TypeString:
			d.Set(FLAG_TYPE, "string")
			defaultValue := d.Get(DEFAULT_VALUE).(string)
			value, err := client.StringVariation(flagKey, userCtx, defaultValue)
			if err != nil {
				return diag.FromErr(err)
			}
			err = d.Set(VALUE, value)
			if err != nil {
				return diag.FromErr(err)
			}
			d.Set(DEFAULT_VALUE, defaultValue)
		case schema.TypeBool:
			d.Set(FLAG_TYPE, "boolean")
			defaultValue := d.Get(DEFAULT_VALUE).(bool)
			tflog.Debug(ctx, "defaultValue is "+strconv.FormatBool(defaultValue))
			tflog.Debug(ctx, "USER CTX KEY")
			tflog.Debug(ctx, userCtx.GetKey())
			value, err := client.BoolVariation(flagKey, userCtx, defaultValue)
			tflog.Debug(ctx, "VALUE VALUE")
			tflog.Debug(ctx, strconv.FormatBool(value))
			if err != nil {
				return diag.FromErr(err)
			}
			err = d.Set(VALUE, value)
			if err != nil {
				return diag.FromErr(err)
			}
			d.Set(DEFAULT_VALUE, defaultValue)
		case schema.TypeInt:
			d.Set(FLAG_TYPE, "int")
			defaultValue := d.Get(DEFAULT_VALUE).(int)
			value, err := client.IntVariation(flagKey, userCtx, defaultValue)
			if err != nil {
				return diag.FromErr(err)
			}
			err = d.Set(VALUE, value)
			if err != nil {
				return diag.FromErr(err)
			}
			d.Set(DEFAULT_VALUE, defaultValue)
		case schema.TypeFloat:
			d.Set(FLAG_TYPE, "float")
			defaultValue := d.Get(DEFAULT_VALUE).(float64)
			value, err := client.Float64Variation(flagKey, userCtx, defaultValue)
			if err != nil {
				return diag.FromErr(err)
			}
			err = d.Set(VALUE, value)
			if err != nil {
				return diag.FromErr(err)
			}
			d.Set(DEFAULT_VALUE, defaultValue)
			// case schema.TypeMap:
			// 	var jsonRaw json.RawMessage
			// 	err := jsonRaw.UnmarshalJSON([]byte(rawDefault))
			// 	if err != nil {
			// 		return diag.FromErr(err)
			// 	}

			// 	defaultValue := ldvalue.Raw(jsonRaw)
			// 	value, err := client.JSONVariation(flagKey, userCtx, defaultValue)
			// 	if err != nil {
			// 		return diag.FromErr(err)
			// 	}
		}

		d.Set(FLAG_KEY, flagKey)
		// TODO we need helper functions to convert back and forth
		// d.Set(CONTEXT, context)
		d.SetId(flagKey)

		return diags
	}
}
