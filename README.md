# terraform-provider-launchdarkly-flag-evaluation

### Test locally

Run `make install` to build the binary and install it to the relevant directory (`~/.terraform.d/plugins/registry.terraform.io/launchdarkly/feature-flag-eval/0.2/${YOUR_OS_ARCH}`)

Verify if `tfenv` is set. If not,
* check which versions are supported locally
```
tfenv list
  1.1.6
  1.1.4
* 0.14.11 (set by /usr/local/Cellar/tfenv/2.2.3/version)
  0.13.5
  0.13.2
```
* Pick the latest one to set it.
```
tfenv use 1.1.6
```

from the examples directory, run 
* `terraform init`
* `./../scripts/tfdev.sh --rebuild apply`
