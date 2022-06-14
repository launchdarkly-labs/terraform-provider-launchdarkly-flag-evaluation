package launchdarkly_flag_eval

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
	"gopkg.in/launchdarkly/go-server-sdk.v5/ldcomponents"
)

const LAUNCHDARKLY_SDK_KEY = "LAUNCHDARKLY_SDK_KEY"

const sdk_key = "sdk_key"

var stderr = os.Stderr

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	client     *ld.LDClient
}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			sdk_key: {
				Type:        types.StringType,
				Required:    true,
				Sensitive:   true,
				Description: "The LaunchDarkly SDK key associated with the project and environment you would like to evaluate flags on",
			},
			// TODO add attributes:
			// host
			// polling/streaming
			// other things?
		},
	}, nil
}

type providerData struct {
	sdkKey types.String `tfsdk:"sdk_key"`
}

// func Provider() *schema.Provider {
// 	return &schema.Provider{
// 		Schema: map[string]*schema.Schema{
// sdk_key: {
// 	Type:        schema.TypeString,
// 	Required:    true,
// 	ForceNew:    true,
// 	DefaultFunc: schema.EnvDefaultFunc(LAUNCHDARKLY_SDK_KEY, nil),
// 	Description: "The LaunchDarkly SDK key associated with the project and environment you would like to evaluate flags on",
// },
// 		},
// 		DataSourcesMap: map[string]*schema.Resource{
// "feature-flag-eval_boolean": dataSourceFlagEvaluation(schema.TypeBool),
// "feature-flag-eval_string":  dataSourceFlagEvaluation(schema.TypeString),
// "feature-flag-eval_int":     dataSourceFlagEvaluation(schema.TypeInt),
// "feature-flag-eval_float":   dataSourceFlagEvaluation(schema.TypeFloat),
// // "launchdarkly_flag_evaluation_json":    dataSourceFlagEvaluation(schema.TypeMap),
// 		},
// 		ConfigureContextFunc: configureSDK,
// 	}
// }

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var sdkKey string
	if config.sdkKey.Unknown || config.sdkKey.Null {
		sdkKey = os.Getenv(LAUNCHDARKLY_SDK_KEY)
	} else {
		sdkKey = config.sdkKey.Value
	}

	if sdkKey == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"SDK key not found",
			"SDK Key must be provided",
		)
		return
	}

	var ldConfig ld.Config
	// TODO is there a shutdown call?
	ldConfig.Events = ldcomponents.SendEvents().FlushInterval(10 * time.Second)
	// default poll interval is 30 seconds
	ldConfig.DataSource = ldcomponents.PollingDataSource()
	c, err := ld.MakeCustomClient(sdkKey, ldConfig, 5*time.Second)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to initialize LaunchDarkly client:\n\n"+err.Error(),
		)
		return
	}

	p.client = c
	p.configured = true
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{}, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"feature-flag-eval_boolean": dataSourceFlagEvaluationBooleanType{},
		// "feature-flag-eval_string":  dataSourceFlagEvaluationStringType{},
		// "feature-flag-eval_int":     dataSourceFlagEvaluationIntType{},
		// "feature-flag-eval_float":   dataSourceFlagEvaluationFloatType{},
		// "feature-flag-eval_json":   dataSourceFlagEvaluationJSONType{},
	}, nil
}
