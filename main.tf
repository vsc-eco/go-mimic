provider "mongodbatlas" { }

resource "mongodbatlas_project" "test" {
  name   = "go-mimic-test"
  org_id = var.mongoatlas_org_id
}

resource "mongodbatlas_cluster" "test" {
  project_id = mongodbatlas_project.test.id
  provider_name = "TENANT"
  name = "cluster-testdb"
  provider_instance_size_name = "M0"
  provider_region_name        = "US_EAST_1"  # required for M0
  backing_provider_name       = "AWS"
}