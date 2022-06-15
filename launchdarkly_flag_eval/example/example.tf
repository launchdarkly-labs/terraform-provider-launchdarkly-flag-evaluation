terraform {
  required_providers {
    feature-flag-eval = {
      version = "0.2"
      source  = "launchdarkly/feature-flag-eval"
    }
  }
}

provider "feature-flag-eval" {
  sdk_key = "sdk-5b155426-b39a-491d-bc2a-cdfb42f8f3c6"
}

data "feature-flag-eval_boolean" "mybool" {
  flag_key      = "terraform-eval"
  default_value = false
  context {
    key = "mchheda@launchdarkly.com"
    custom = tomap({
      "host" = "mchheda-local"
    })
  }
}

output "variation_value" {
  value = data.feature-flag-eval_boolean.mybool
}

# locals {
#   foo = file(data.feature-flag-eval_boolean.mybool.flag_key)
# }
