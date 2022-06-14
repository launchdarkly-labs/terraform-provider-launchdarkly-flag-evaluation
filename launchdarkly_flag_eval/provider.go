package launchdarkly_flag_eval

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
	"gopkg.in/launchdarkly/go-server-sdk.v5/ldcomponents"
)

const LAUNCHDARKLY_SDK_KEY = "LAUNCHDARKLY_SDK_KEY"

const sdk_key = "sdk_key"

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			sdk_key: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc(LAUNCHDARKLY_SDK_KEY, nil),
				Description: "The LaunchDarkly SDK key associated with the project and environment you would like to evaluate flags on",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"feature-flag-eval_boolean": dataSourceFlagEvaluation(schema.TypeBool),
			"feature-flag-eval_string":  dataSourceFlagEvaluation(schema.TypeString),
			"feature-flag-eval_int":     dataSourceFlagEvaluation(schema.TypeInt),
			"feature-flag-eval_float":   dataSourceFlagEvaluation(schema.TypeFloat),
			// "launchdarkly_flag_evaluation_json":    dataSourceFlagEvaluation(schema.TypeMap),
		},
		ConfigureContextFunc: configureSDK,
	}
}

func configureSDK(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	var config ld.Config

	sdkKey := d.Get(sdk_key).(string)

	// TODO is there a shutdown call?
	config.Events = ldcomponents.SendEvents().FlushInterval(10 * time.Second)
	// default poll interval is 30 seconds
	config.DataSource = ldcomponents.PollingDataSource()
	client, err := ld.MakeCustomClient(sdkKey, config, 5*time.Second)
	if err != nil {
		return client, diag.FromErr(err)
	}
	return client, diag.Diagnostics{}
}
