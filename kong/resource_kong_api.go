package kong

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

// APIRequest : Kong API request object structure
type APIRequest struct {
	ID                     string `json:"id,omitempty"`
	Name                   string `json:"name"`
	Hosts                  string `json:"hosts,omitempty"`
	Uris                   string `json:"uris,omitempty"`
	Methods                string `json:"methods,omitempty"`
	UpstreamURL            string `json:"upstream_url"`
	StripURI               bool   `json:"strip_uri"`
	PreserveHost           bool   `json:"preserve_host"`
	Retries                int    `json:"retries,omitempty"`
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout,omitempty"`
	UpstreamSendTimeout    int    `json:"upstream_send_timeout,omitempty"`
	UpstreamReadTimeout    int    `json:"upstream_read_timeout,omitempty"`
	HTTPSOnly              bool   `json:"https_only"`
	HTTPIfTerminated       bool   `json:"http_if_terminated"`
}

// APIResponse : Kong API response object structure
type APIResponse struct {
	ID                     string   `json:"id,omitempty"`
	Name                   string   `json:"name"`
	Hosts                  []string `json:"hosts,omitempty"`
	Uris                   []string `json:"uris,omitempty"`
	Methods                []string `json:"methods,omitempty"`
	UpstreamURL            string   `json:"upstream_url"`
	StripURI               bool     `json:"strip_uri"`
	PreserveHost           bool     `json:"preserve_host"`
	Retries                int      `json:"retries,omitempty"`
	UpstreamConnectTimeout int      `json:"upstream_connect_timeout,omitempty"`
	UpstreamSendTimeout    int      `json:"upstream_send_timeout,omitempty"`
	UpstreamReadTimeout    int      `json:"upstream_read_timeout,omitempty"`
	HTTPSOnly              bool     `json:"https_only"`
	HTTPIfTerminated       bool     `json:"http_if_terminated"`
}

func resourceKongAPI() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongAPICreate,
		Read:   resourceKongAPIRead,
		Update: resourceKongAPIUpdate,
		Delete: resourceKongAPIDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The API name.",
			},

			"hosts": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "A comma-separated list of domain names that point to your API. For example: example.com. At least one of hosts, uris, or methods should be specified.",
			},

			"uris": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "A comma-separated list of URIs prefixes that point to your API. For example: /my-path. At least one of hosts, uris, or methods should be specified",
			},

			"methods": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "A comma-separated list of HTTP methods that point to your API. For example: GET,POST. At least one of hosts, uris, or methods should be specified.",
			},

			"upstream_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The base target URL that points to your API server. This URL will be used for proxying requests. For example: https://example.com.",
			},

			"strip_uri": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When matching an API via one of the uris prefixes, strip that matching prefix from the upstream URI to be requested. Default: true.",
			},

			"preserve_host": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When matching an API via one of the hosts domain names, make sure the request Host header is forwarded to the upstream service. By default, this is false, and the upstream Host header will be extracted from the configured upstream_url.",
			},

			"retries": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "A comma-separated list of HTTP methods that point to your API. For example: GET,POST. At least one of hosts, uris, or methods should be specified.",
			},

			"upstream_connect_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60000,
				Description: "The timeout in milliseconds for establishing a connection to your upstream service. Defaults to 60000.",
			},

			"upstream_send_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60000,
				Description: "The timeout in milliseconds between two successive write operations for transmitting a request to your upstream service Defaults to 60000.",
			},

			"upstream_read_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60000,
				Description: "The timeout in milliseconds between two successive read operations for transmitting a request to your upstream service Defaults to 60000.",
			},

			"https_only": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     nil,
				Description: "To be enabled if you wish to only serve an API through HTTPS, on the appropriate port (8443 by default). Default: false.",
			},

			"http_if_terminated": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Consider the X-Forwarded-Proto header when enforcing HTTPS only traffic. Default: true.",
			},
		},
	}
}

func resourceKongAPICreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	api := getAPIFromResourceData(d)

	createdAPI := new(APIResponse)
	response, error := sling.New().BodyJSON(api).Post("apis/").ReceiveSuccess(createdAPI)

	if error != nil {
		return fmt.Errorf("error while creating API: " + error.Error())
	}

	if response.StatusCode == http.StatusConflict {
		return fmt.Errorf("409 Conflict - use terraform import to manage this api.")
	} else if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	setAPIToResourceData(d, createdAPI)

	return nil
}

func resourceKongAPIRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	id := d.Get("id").(string)
	api := new(APIResponse)

	response, error := sling.New().Path("apis/").Get(id).ReceiveSuccess(api)

	if error != nil {
		return fmt.Errorf("error while updating API" + error.Error())
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	setAPIToResourceData(d, api)

	return nil
}

func resourceKongAPIUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	api := getAPIFromResourceData(d)

	updatedAPI := new(APIResponse)

	response, error := sling.New().BodyJSON(api).Patch("apis/").Path(api.ID).ReceiveSuccess(updatedAPI)

	if error != nil {
		return fmt.Errorf("error while updating API" + error.Error())
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	setAPIToResourceData(d, updatedAPI)

	return nil
}

func resourceKongAPIDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	id := d.Get("id").(string)

	response, error := sling.New().Delete("apis/").Path(id).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("error while deleting API" + error.Error())
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	return nil
}

func getAPIFromResourceData(d *schema.ResourceData) *APIRequest {
	api := &APIRequest{
		Name:                   d.Get("name").(string),
		Hosts:                  d.Get("hosts").(string),
		Uris:                   d.Get("uris").(string),
		Methods:                d.Get("methods").(string),
		UpstreamURL:            d.Get("upstream_url").(string),
		StripURI:               d.Get("strip_uri").(bool),
		PreserveHost:           d.Get("preserve_host").(bool),
		Retries:                d.Get("retries").(int),
		UpstreamConnectTimeout: d.Get("upstream_connect_timeout").(int),
		UpstreamSendTimeout:    d.Get("upstream_send_timeout").(int),
		UpstreamReadTimeout:    d.Get("upstream_read_timeout").(int),
		HTTPSOnly:              d.Get("https_only").(bool),
		HTTPIfTerminated:       d.Get("http_if_terminated").(bool),
	}

	if id, ok := d.GetOk("id"); ok {
		api.ID = id.(string)
	}

	return api
}

func setAPIToResourceData(d *schema.ResourceData, api *APIResponse) {
	d.SetId(api.ID)
	d.Set("name", api.Name)
	d.Set("hosts", api.Hosts)
	d.Set("uris", strings.Join(api.Uris, ","))
	d.Set("methods", api.Methods)
	d.Set("upstream_url", api.UpstreamURL)
	d.Set("strip_uri", api.StripURI)
	d.Set("preserve_host", api.PreserveHost)
	d.Set("retries", api.Retries)
	d.Set("upstream_connect_timeout", api.UpstreamConnectTimeout)
	d.Set("upstream_send_timeout", api.UpstreamSendTimeout)
	d.Set("upstream_read_timeout", api.UpstreamReadTimeout)
	d.Set("https_only", api.HTTPSOnly)
	d.Set("http_if_terminated", api.HTTPIfTerminated)
}
