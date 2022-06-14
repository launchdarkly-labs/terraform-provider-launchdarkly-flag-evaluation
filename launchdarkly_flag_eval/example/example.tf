provider "launchdarkly_flag-evals" {
  sdk_key = "bla"
}

data "launchdarkly_flag_evaluation_boolean" "mybool" {
  flagKey = "bla"
  default_value = false
  context = {
    key = "hosh@launchdarkly"
  }
}

resource "aws_ec2_instance" "myuinstance" {
  count = data.launchdarkly_flag_evaluation_boolean.mybool.value ? 1 : 0

}
