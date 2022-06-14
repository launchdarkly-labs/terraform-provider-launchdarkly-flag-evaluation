#!/bin/bash

set -euo pipefail

root="$(git rev-parse --show-toplevel)"
tmp="$root/tmp"
bin="$tmp/bin"

TF_CLI_CONFIG_FILE="$tmp/dev.tfrc"
export TF_CLI_CONFIG_FILE

if [[ ! -f $TF_CLI_CONFIG_FILE ]]; then
  mkdir -p "$tmp" "$bin"

  cat <<EOF >"$TF_CLI_CONFIG_FILE"
provider_installation {
  dev_overrides {
    "launchdarkly/feature-flag-eval" = "${bin}"
  }
  direct {}
}
EOF
fi

if [[ ${1:-} == "--rebuild" ]]; then
  go build -o "$bin/terraform-provider-feature-flag-eval" "$root"
  shift
fi

exec terraform "$@"
