#!/bin/bash

set -euo pipefail

root="$(git rev-parse --show-toplevel)"
tmp="$root/tmp"
bin="$tmp/bin"

TF_CLI_CONFIG_FILE="$tmp/dev.tfrc"

if [[ ${1:-} == "--rebuild" ]]; then
  rm -f "$TF_CLI_CONFIG_FILE"
  go build -o "$bin/terraform-provider-ldflags" "$root"
  shift
fi

export TF_CLI_CONFIG_FILE

if [[ ! -f $TF_CLI_CONFIG_FILE ]]; then
  mkdir -p "$tmp" "$bin"

  cat <<EOF >"$TF_CLI_CONFIG_FILE"
provider_installation {
  dev_overrides {
    "launchdarkly/ldflags" = "${bin}"
  }
  direct {}
}
EOF
fi

exec terraform "$@"
