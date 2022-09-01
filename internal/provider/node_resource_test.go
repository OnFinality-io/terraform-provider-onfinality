package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("onfinality_node.test", "workspace_id", "6635707676612587520"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "onfinality_node.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{},
			},
			// Update and Read testing
			{
				Config: testAccExampleResourceConfig2(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("onfinality_node.test", "node_name", "ian test2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleResourceConfig() string {
	return fmt.Sprintf(`
resource "onfinality_node" "test" {
  workspace_id         = 6635707676612587520
  network_spec_key     = "polkadot"
  node_spec = {
    key = "unit"
    multiplier = 4
  }
  node_type            = "full"
  node_name            = "ian test"
  cluster_hash         = "jm"
  storage              = "150Gi"
  image_version        = "v0.9.27"
}
`)
}

func testAccExampleResourceConfig2() string {
	return fmt.Sprintf(`
resource "onfinality_node" "test" {
  workspace_id         = 6635707676612587520
  network_spec_key     = "polkadot"
  node_spec = {
    key = "unit"
    multiplier = 4
  }
  node_type            = "full"
  node_name            = "ian test2"
  cluster_hash         = "jm"
  storage              = "150Gi"
  image_version        = "v0.9.27"
}
`)
}
