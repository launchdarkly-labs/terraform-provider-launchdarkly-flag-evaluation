terraform {
  required_providers {
    feature-flag-eval = {
      version = "0.2"
      source  = "launchdarkly/feature-flag-eval"
    }
  }
}

provider "feature-flag-eval" {
  sdk_key = "bla"
}

#data "launchdarkly_flag_evaluation_boolean" "mybool" {
#  flagKey = "bla"
#  default_value = false
#  context = {
#    key = "hosh@launchdarkly"
#  }
#}
#
#
#locals {
#  foo = file(yamlencode(data.launchdarkly_flag_evaluation_boolean.mybool))
#}
