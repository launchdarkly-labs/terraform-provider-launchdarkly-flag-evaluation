locals {
  app    = "ldflags-app"
  labels = {
    app = "myapp"
  }
}

#data "ldflags_evaluation_string" "nginx_version" {
#  count         = 2
#  flag_key      = "k8s-nginx-version"
#  default_value = "1.20.0"
#  context       = {
#    key = "${local.app}-${count.index}"
#  }
#}

#data "ldflags_evaluation_int" "k8s_replicas" {
#  flag_key      = "k8s-replicas"
#  default_value = 2
#
#  context = {
#    key = "terraform-user"
#  }
#}

resource "kubernetes_deployment" "ldflags_app" {
  count = 2

  metadata {
    name   = "${local.app}-${count.index}"
    labels = local.labels
  }

  spec {
    replicas = 1

    selector {
      match_labels = local.labels
    }

    template {
      metadata {
        labels = local.labels
      }

      spec {
        container {
#          image = "nginx:${data.ldflags_evaluation_string.nginx_version[count.index].value}"
          image = "nginx:1.20.0"
          name  = "nginx-app"

          resources {
            limits = {
              cpu    = "0.5"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "50Mi"
            }
          }

          liveness_probe {
            http_get {
              path = "/"
              port = 80

              http_header {
                name  = "X-Custom-Header"
                value = "Awesome"
              }
            }

            initial_delay_seconds = 3
            period_seconds        = 3
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "ldflags_app" {
  metadata {
    name = "myapp"
  }
  spec {
    selector = local.labels

    session_affinity = "ClientIP"
    port {
      port        = 8080
      target_port = 80
    }

    type = "NodePort"
  }
}
