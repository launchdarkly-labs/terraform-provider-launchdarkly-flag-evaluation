package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/launchdarkly-labs/terraform-provider-launchdarkly-flag-evaluation/launchdarkly_flag_eval"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: launchdarkly_flag_eval.Provider})
}
