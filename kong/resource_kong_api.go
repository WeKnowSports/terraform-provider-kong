package kong

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKongApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongApiCreate,
		Read:   resourceKongApiRead,
		Update: resourceKongApiUpdate,
		Delete: resourceKongApiDelete,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The API name. If none is specified, will default to the request_host or request_path.",
			},

			"request_host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The public DNS address that points to your API. For example, mockbin.com. At least request_host or request_path or both should be specified.",
			},

			"request_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The public path that points to your API. For example, /someservice. At least request_host or request_path or both should be specified.",
			},

			"strip_request_path": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Strip the request_path value before proxying the request to the final API. For example a request made to /someservice/hello will be resolved to upstream_url/hello. By default is false.",
			},

			"preserve_host": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Preserves the original Host header sent by the client, instead of replacing it with the hostname of the upstream_url. By default is false.",
			},

			"upstream_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The base target URL that points to your API server, this URL will be used for proxying requests. For example, https://mockbin.com.",
			},
		},
	}
}

func resourceKongApiCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKongApiRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKongApiUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKongApiDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
