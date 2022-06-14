terraform {
  required_providers {
    feature-flag-eval = {
      version = "0.2"
      source  = "launchdarkly/feature-flag-eval"
    }
  }
}

provider "feature-flag-eval" {
  sdk_key = "sdk-2aaff62c-d031-47d1-9594-add194fce944"
}

data "feature-flag-eval_boolean" "mybool" {
  flag_key      = "boolean-flag"
  default_value = false
  context = {
    key = "hosh@launchdarkly.com"
  }
}

output "variation_value" {
  value = data.feature-flag-eval_boolean.mybool
}

# locals {
#   foo = file(data.feature-flag-eval_boolean.mybool.flag_key)
# }
