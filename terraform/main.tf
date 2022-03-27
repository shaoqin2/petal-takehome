resource "google_cloud_run_service" "default" {
  name     = "cloudrun-srv"
  project  = "petal-takehome"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/petal-takehome/reverseserver:prod-2022-03-27"
        resources {
          requests = {
            cpu    = 1
            memory = "128Mi"
          }
          limits = {
            cpu    = 1
            memory = "128Mi"
          }
        }
        ports {
          container_port = 8090
        }
      }
      timeout_seconds = 3600
    }
    metadata {
      annotations = {
        # put in a low value so my google account doesn't bankrupt me if petal decides to load test my service.....
        "autoscaling.knative.dev/maxScale" = "3"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}


data "google_iam_policy" "noauth" {
  binding {
    role    = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.default.location
  project  = google_cloud_run_service.default.project
  service  = google_cloud_run_service.default.name

  policy_data = data.google_iam_policy.noauth.policy_data
}