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
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The API name. If none is specified, will default to the request_host or request_path.",
			},

			"request_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The public DNS address that points to your API. For example, mockbin.com. At least request_host or request_path or both should be specified.",
			},

			"request_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The public path that points to your API. For example, /someservice. At least request_host or request_path or both should be specified.",
			},

			"strip_request_path": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Strip the request_path value before proxying the request to the final API. For example a request made to /someservice/hello will be resolved to upstream_url/hello. By default is false.",
			},

			"preserve_host": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Preserves the original Host header sent by the client, instead of replacing it with the hostname of the upstream_url. By default is false.",
			},

			"upstream_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The base target URL that points to your API server, this URL will be used for proxying requests. For example, https://mockbin.com.",
			},
		},
	}
}

func resourceKongAPICreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	api := getAPIFromResourceData(d)

	createdAPI := new(API)

	response, error := sling.New().BodyJSON(api).Post("apis/").ReceiveSuccess(createdAPI)
	if error != nil {
		return fmt.Errorf("Error while creating API.")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setAPIToResourceData(d, createdAPI)

	return nil
}

func resourceKongAPIRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	id := d.Get("id").(string)
	api := new(API)

	response, error := sling.New().Path("apis/").Get(id).ReceiveSuccess(api)
	if error != nil {
		return fmt.Errorf("Error while updating API.")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setAPIToResourceData(d, api)

	return nil
}

func resourceKongAPIUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	api := getAPIFromResourceData(d)

	updatedAPI := new(API)

	response, error := sling.New().BodyJSON(api).Patch("apis/").Path(api.ID).ReceiveSuccess(updatedAPI)
	if error != nil {
		return fmt.Errorf("Error while updating API.")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setAPIToResourceData(d, updatedAPI)

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

func getAPIFromResourceData(d *schema.ResourceData) *API {
	api := &API{
		Name:             d.Get("name").(string),
		RequestHost:      d.Get("request_host").(string),
		RequestPath:      d.Get("request_path").(string),
		StripRequestPath: d.Get("strip_request_path").(bool),
		PreserveHost:     d.Get("preserve_host").(bool),
		UpstreamURL:      d.Get("upstream_url").(string),
	}

	if id, ok := d.GetOk("id"); ok {
		api.ID = id.(string)
	}

	return api
}

func setAPIToResourceData(d *schema.ResourceData, api *API) {
	d.SetId(api.ID)
	d.Set("name", api.Name)
	d.Set("request_host", api.RequestHost)
	d.Set("request_path", api.RequestPath)
	d.Set("strip_request_path", api.StripRequestPath)
	d.Set("preserve_host", api.PreserveHost)
	d.Set("upstream_url", api.UpstreamURL)
}
