# terraform-provider-launchdarkly-flag-evaluation

This repo contains an `ldflags` provider that wraps the [LaunchDarkly Go SDK](https://docs.launchdarkly.com/sdk/server-side/go). The provider directory `ldflags/` contains five flag evaluation data sources:

- `ldflags_evaluation_boolean`
- `ldflags_evaluation_string`
- `ldflags_evaluation_int`
- `ldflags_evaluation_float`
- `ldflags_evaluation_json`

Examples of each can be found in `ldflags/example/example.tf`.

### Test locally

To get set up to contribute to the provider, first run `make install` to build the binary and install it to the relevant directory (`~/.terraform.d/plugins/registry.terraform.io/launchdarkly/ldflags/0.2/${YOUR_OS_ARCH}`).

Verify if `tfenv` is set. If not, check which versions are supported locally using `tfenv list`.

```
tfenv list
  1.1.6
  1.1.4
* 0.14.11 (set by /usr/local/Cellar/tfenv/2.2.3/version)
  0.13.5
  0.13.2
```

Use `tfenv use <version>` to set it to the latest version.

```
tfenv use 1.1.6
```

From the examples directory, run the following:

- `terraform init`
- `./../scripts/tfdev.sh --rebuild apply`

The tfdev.sh script will rebuild and reinstall the binary and also write a dev_overrides configuration to ensure that your `ldflags` provider configurations will point to your local binary.
