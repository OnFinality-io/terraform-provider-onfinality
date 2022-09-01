terraform {
  required_providers {
    onfinality = {
      source  = "OnFinality-io/onfinality"
      version = "0.1.0"
    }
  }
}

variable "onf_access_key" {}
variable "onf_secret_key" {}

provider "onfinality" {
  access_key = var.onf_access_key
  secret_key = var.onf_secret_key
}

resource "onfinality_node" "n1" {
  workspace_id     = 6635707676612587520
  network_spec_key = "polkadot"
  node_spec = {
    key        = "unit"
    multiplier = 4
  }
  node_type     = "full"
  node_name     = "ian test2"
  cluster_hash  = "jm"
  storage       = "100Gi"
  image_version = "v0.9.27"
  stopped       = false
}
