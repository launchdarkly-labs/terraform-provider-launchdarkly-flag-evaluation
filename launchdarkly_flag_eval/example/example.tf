terraform {
  required_providers {
    feature-flag-eval = {
      version = "0.2"
      source  = "launchdarkly.com/launchdarkly/feature-flag-eval"
    }
  }
}

provider "feature-flag-eval" {
  sdk_key = "sdk-2aaff62c-d031-47d1-9594-add194fce944"
}

data "launchdarkly_flag_evaluation_boolean" "mybool" {
  flagKey       = "boolean-flag"
  default_value = false
  context = {
    key = "hosh@launchdarkly.com"
  }
}

output "variation_value" {
  value = data.launcharkly_flag_evaluation_boolean.mybool
}


locals {
  foo = file(yamlencode(data.launchdarkly_flag_evaluation_boolean.mybool))
}
