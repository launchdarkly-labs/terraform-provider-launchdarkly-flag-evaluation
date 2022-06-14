package launchdarkly_flag_eval

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
)

const (
	FLAG_KEY       = "flag_key"
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
	firstNameAttribute = "firstName"
	// LastNameAttribute is the standard attribute name corresponding to User.GetLastName().
	lastNameAttribute = "lastName"
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
			DEFAULT_VALUE: {
				Type:     typ,
				Required: true,
			},
			VALUE: {
				Type:     typ,
				Optional: true,
				Computed: true,
			},
			// TODO: figure out the best name for this
			CONTEXT: {
				Type:     schema.TypeList,
				MaxItems: 1,
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
		var diags diag.Diagnostics
		client := meta.(*ld.LDClient)

		flagKey := d.Get(FLAG_KEY).(string)
		rawContext := d.Get(CONTEXT).([]interface{})
		// TODO construct user object properly
		rawContextMap := rawContext[0].(map[string]*schema.Schema)
		_ = rawContextMap
		//userCtxBuilder := lduser.NewUserBuilder(rawContextMap[keyAttribute])
		userCtxBuilder := lduser.NewUserBuilder("hello-world")
		//userCtxBuilder.Name(rawContextMap[nameAttribute])
		userCtx := userCtxBuilder.Build()
		// todo rest of userCtx

		switch typ {
		case schema.TypeString:
			defaultValue := d.Get(DEFAULT_VALUE).(string)
			value, err := client.StringVariation(flagKey, userCtx, defaultValue)
			if err != nil {
				return diag.FromErr(err)
			}
			err = d.Set(VALUE, value)
			if err != nil {
				return diag.FromErr(err)
			}
		case schema.TypeBool:
			defaultValue := d.Get(DEFAULT_VALUE).(bool)
			value, err := client.BoolVariation(flagKey, userCtx, defaultValue)
			if err != nil {
				return diag.FromErr(err)
			}
			err = d.Set(VALUE, value)
			if err != nil {
				return diag.FromErr(err)
			}
		case schema.TypeInt:
			defaultValue := d.Get(DEFAULT_VALUE).(int)
			value, err := client.IntVariation(flagKey, userCtx, defaultValue)
			if err != nil {
				return diag.FromErr(err)
			}
			err = d.Set(VALUE, value)
			if err != nil {
				return diag.FromErr(err)
			}
		case schema.TypeFloat:
			defaultValue := d.Get(DEFAULT_VALUE).(float64)
			value, err := client.Float64Variation(flagKey, userCtx, defaultValue)
			if err != nil {
				return diag.FromErr(err)
			}
			err = d.Set(VALUE, value)
			if err != nil {
				return diag.FromErr(err)
			}
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

		return diags
	}
}
