package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/launchdarkly-labs/terraform-provider-launchdarkly-flag-evaluation/launchdarkly_flag_eval"
)

func main() {
	providerserver.Serve(context.Background(), launchdarkly_flag_eval.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/launchdarkly/feature-flag-eval",
	})
}
