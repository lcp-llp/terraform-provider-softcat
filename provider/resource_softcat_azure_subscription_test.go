package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSoftcatAzureSubscription_basic(t *testing.T) {
	if os.Getenv("SOFTCAT_ENABLE_AZURE_SUBSCRIPTION_ACC") == "" {
		t.Skip("set SOFTCAT_ENABLE_AZURE_SUBSCRIPTION_ACC=1 to run this acceptance test skeleton")
	}

	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSoftcatAzureSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSoftcatAzureSubscriptionConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("softcat_azure_subscription.test", "id"),
					resource.TestCheckResourceAttrSet("softcat_azure_subscription.test", "order_id"),
					resource.TestCheckResourceAttrSet("softcat_azure_subscription.test", "status"),
				),
			},
		},
	})
}

func testAccCheckSoftcatAzureSubscriptionDestroy(_ *terraform.State) error {
	return nil
}

func testAccSoftcatAzureSubscriptionConfig() string {
	return fmt.Sprintf(`
provider "softcat" {
  hostname = %q
  username = %q
  password = %q
}

resource "softcat_azure_subscription" "test" {
  basket_name    = "acctest-softcat-order"
  msid           = %q
  azure_budget   = "1"
  azure_nickname = "acctest-softcat-sub"
  azure_contact  = %q
  quantity       = 1

  checkout_data {
    purchase_order_number  = "acctest-softcat-order"
    additional_information = %q
    csp_terms              = true
  }
}
`,
		os.Getenv("SOFTCAT_HOSTNAME"),
		os.Getenv("SOFTCAT_USERNAME"),
		os.Getenv("SOFTCAT_PASSWORD"),
		os.Getenv("SOFTCAT_MSID"),
		os.Getenv("SOFTCAT_AZURE_CONTACT"),
		os.Getenv("SOFTCAT_AZURE_CONTACT"),
	)
}
