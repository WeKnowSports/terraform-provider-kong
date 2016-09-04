package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

type API struct {
	ID               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	RequestHost      string `json:"request_host,omitempty"`
	RequestPath      string `json:"request_path,omitempty"`
	StripRequestPath bool   `json:"strip_request_path,omitempty"`
	PreserveHost     bool   `json:"preserve_host,omitempty"`
	UpstreamURL      string `json:"upstream_url"`
}

func resourceKongAPI() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongAPICreate,
		Read:   resourceKongAPIRead,
		Update: resourceKongAPIUpdate,
		Delete: resourceKongAPIDelete,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The API name. If none is specified, will default to the request_host or request_path.",
			},

			"request_host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The public DNS address that points to your API. For example, mockbin.com. At least request_host or request_path or both should be specified.",
			},

			"request_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The public path that points to your API. For example, /someservice. At least request_host or request_path or both should be specified.",
			},

			"strip_request_path": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Strip the request_path value before proxying the request to the final API. For example a request made to /someservice/hello will be resolved to upstream_url/hello. By default is false.",
			},

			"preserve_host": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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

func resourceKongAPICreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	api := &API{
		Name:             d.Get("name").(string),
		RequestHost:      d.Get("request_host").(string),
		RequestPath:      d.Get("request_path").(string),
		StripRequestPath: d.Get("strip_request_path").(bool),
		PreserveHost:     d.Get("preserve_host").(bool),
		UpstreamURL:      d.Get("upstream_url").(string),
	}

	createdAPI := new(API)

	response, error := sling.New().BodyJSON(api).Post("apis/").ReceiveSuccess(createdAPI)
	if error != nil {
		return fmt.Errorf("Error while creating API.")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	d.SetId(createdAPI.ID)
	d.Set("name", createdAPI.Name)
	d.Set("request_host", createdAPI.RequestHost)
	d.Set("request_path", createdAPI.RequestPath)
	d.Set("strip_request_path", createdAPI.StripRequestPath)
	d.Set("preserve_host", createdAPI.PreserveHost)
	d.Set("upstream_url", createdAPI.UpstreamURL)

	return nil
}

func resourceKongAPIRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKongAPIUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKongAPIDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	id := d.Get("id").(string)

	response, error := sling.New().Delete("apis/").Path(id).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("Error while deleting API.")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}
