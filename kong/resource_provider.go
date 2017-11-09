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
				Type:     schema.TypeString,
				Required: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional:    true,
				Default: "",
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional:    true,
				Default: "",
			},
			"headers": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Default:  nil,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"kong_api":                            resourceKongAPI(),
			"kong_consumer":                       resourceKongConsumer(),
			"kong_consumer_acl":                   resourceKongConsumerACL(),
			"kong_api_plugin":                     resourceKongPlugin(),
			"kong_plugin":                         resourceKongPlugin(),
			"kong_consumer_basic_auth_credential": resourceKongBasicAuthCredential(),
			"kong_consumer_key_auth_credential":   resourceKongKeyAuthCredential(),
			"kong_consumer_jwt_credential":        resourceKongJWTCredential(),
			"kong_api_plugin_key_auth":            resourceKongKeyAuthPlugin(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Address: d.Get("address").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Headers:  d.Get("headers").(map[string]interface{}),
	}

	return config.Client()
}
