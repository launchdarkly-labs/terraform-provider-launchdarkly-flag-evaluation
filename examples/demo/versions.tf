terraform {
  required_providers {
    ldflags = {
      version = "0.2"
      source  = "launchdarkly/ldflags"
    }

    kubernetes = {
      version = "2.11.0"
      source  = "hashicorp/kubernetes"
    }
  }
}

provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "minikube"
}

provider "ldflags" {
  sdk_key = "sdk-2aaff62c-d031-47d1-9594-add194fce944"
}
