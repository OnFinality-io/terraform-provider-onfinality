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
  # Workspace id, can get it from url https://app.onfinality.io/workspaces/<workspace_id>/nodes
  workspace_id     = 6635707676612587520
  # Network of the node, can get from `onf network-spec list` & `onf network-spec list-backups`
  network_spec_key = "polkadot"
  # Node Spec of the node, always put key="unit", 1 * unit ~ 0.5 cpu 1.5G mem
  node_spec = {
    key        = "unit"
    multiplier = 4
  }
  # full or archive or validator, depends on network
  node_type     = "full"
  node_name     = "ian test2"
  # Cluster where the node will be deployed, check `onf info cluster` for all available clusters
  cluster_hash  = "jm"
  # Disk size of the node, <num>Gi , e.g 100Gi
  storage       = "100Gi"
  image_version = "v0.9.27"
  # <Optional> Change it to true will stop the node
  # stopped = false
}
