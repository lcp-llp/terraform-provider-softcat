package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"softcat": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	for _, envVar := range []string{"SOFTCAT_HOSTNAME", "SOFTCAT_USERNAME", "SOFTCAT_PASSWORD", "SOFTCAT_MSID", "SOFTCAT_AZURE_CONTACT"} {
		if os.Getenv(envVar) == "" {
			t.Skipf("acceptance test skipped: %s must be set", envVar)
		}
	}
}
