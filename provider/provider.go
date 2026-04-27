package provider

import (
	"context"
	"fmt"

	"terraform-provider-softcat/internal"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOFTCAT_HOSTNAME", nil),
				Description: "The API Endpoint",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOFTCAT_USERNAME", nil),
				Description: "Username used for authentication to API Endpoints",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SOFTCAT_PASSWORD", nil),
				Description: "Password used for authentication to API Endpoints",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureContextFunc: func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			hostname := d.Get("hostname").(string)
			username := d.Get("username").(string)
			password := d.Get("password").(string)

			client, err := internal.NewClient(hostname, username, password)
			if err != nil {
				return nil, diag.FromErr(fmt.Errorf("error configuring provider: %w", err))
			}

			return client, nil
		},
		ResourcesMap: map[string]*schema.Resource{
			"softcat_azure_subscription": internal.ResourceAzureSubscription(),
		},
	}
}
