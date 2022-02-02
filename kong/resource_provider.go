package kong

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"kong_api":                            resourceKongAPI(),
			"kong_service":                        resourceKongService(),
			"kong_route":                          resourceKongRoute(),
			"kong_consumer":                       resourceKongConsumer(),
			"kong_api_plugin":                     resourceKongPlugin(),
			"kong_plugin":                         resourceKongPlugin(),
			"kong_consumer_basic_auth_credential": resourceKongBasicAuthCredential(),
			"kong_consumer_key_auth_credential":   resourceKongKeyAuthCredential(),
			"kong_consumer_jwt_credential":        resourceKongJWTCredential(),
			"kong_api_plugin_key_auth":            resourceKongKeyAuthPlugin(),
			"kong_consumer_acl_group":             resourceKongConsumerACLGroup(),
			"kong_certificate":                    resourceKongCertificate(),
			"kong_sni":                            resourceKongSNI(),
			"kong_upstream":                       resourceKongUpstream(),
			"kong_target":                         resourceKongTarget(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Address:  d.Get("address").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	return config.Client()
}
