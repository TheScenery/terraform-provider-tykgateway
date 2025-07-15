package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKeyResource(t *testing.T) {

	t.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					resource "tykgateway_key" "test_key" {
						 	alias              = "example-key"
							org_id             = "default"
							hashed             = false
  					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tykgateway_key.test_key", "alias", "example222"),
				),
			},
		},
	})
}
