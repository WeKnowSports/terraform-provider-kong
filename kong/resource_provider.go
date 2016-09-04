package kong

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type: schema.TypeString,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"kong_api": resourceKongApi(),
		},
	}
}
