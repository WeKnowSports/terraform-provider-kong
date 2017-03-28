package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

type API struct {
	ID                     string      `json:"id,omitempty"`
	Name                   string      `json:"name,omitempty"`
	Hosts                  interface{} `json:"hosts,omitempty"`
	URIs                   interface{} `json:"uris,omitempty"`
	StripURI               bool        `json:"strip_uri"`
	PreserveHost           bool        `json:"preserve_host,omitempty"`
	UpstreamURL            string      `json:"upstream_url,omitempty"`
	Methods                interface{} `json:"methods,omitempty"`
	Retries                int         `json:"retries,omitempty"`
	HTTPSOnly              bool        `json:"https_only,omitempty"`
	HTTPIfTerminated       bool        `json:"http_if_terminated,omitempty"`
	UpstreamConnectTimeout int         `json:"upstream_connect_timeout,omitempty"`
	UpstreamSendTimeout    int         `json:"upstream_send_timeout,omitempty"`
	UpstreamReadTimeout    int         `json:"upstream_read_timeout,omitempty"`
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
				Description: "The API name.",
			},

			"hosts": &schema.Schema{
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "A comma-separated list of domain names that point to your API. For example: example.com. At least one of hosts, uris, or methods should be specified.",
			},

			"uris": &schema.Schema{
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "A comma-separated list of URIs prefixes that point to your API. For example: /my-path. At least one of hosts, uris, or methods should be specified.",
			},

			"strip_uri": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
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

			"methods": &schema.Schema{
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "A comma-separated list of HTTP methods that point to your API. For example: GET,POST. At least one of hosts, uris, or methods should be specified.",
			},

			"retries": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The number of retries to execute upon failure to proxy. The default is 5.",
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
		return fmt.Errorf("Error while reading API.")
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
		Name:                   d.Get("name").(string),
		Hosts:                  d.Get("hosts"),
		URIs:                   d.Get("uris"),
		StripURI:               d.Get("strip_uri").(bool),
		PreserveHost:           d.Get("preserve_host").(bool),
		UpstreamURL:            d.Get("upstream_url").(string),
		Methods:                d.Get("methods"),
		Retries:                d.Get("retries").(int),
		HTTPSOnly:              d.Get("https_only").(bool),
		HTTPIfTerminated:       d.Get("http_if_terminated").(bool),
		UpstreamConnectTimeout: d.Get("upstream_connect_timeout").(int),
		UpstreamSendTimeout:    d.Get("upstream_send_timeout").(int),
		UpstreamReadTimeout:    d.Get("upstream_read_timeout").(int),
	}

	if id, ok := d.GetOk("id"); ok {
		api.ID = id.(string)
	}

	return api
}

func setAPIToResourceData(d *schema.ResourceData, api *API) {
	d.SetId(api.ID)
	d.Set("name", api.Name)
	d.Set("hosts", api.Hosts)
	d.Set("uris", api.URIs)
	d.Set("strip_uri", api.StripURI)
	d.Set("preserve_host", api.PreserveHost)
	d.Set("upstream_url", api.UpstreamURL)
	d.Set("methods", api.Methods)
	d.Set("retries", api.Retries)
	d.Set("https_only", api.HTTPSOnly)
	d.Set("http_if_terminated", api.HTTPIfTerminated)
	d.Set("upstream_connect_timeout", api.UpstreamConnectTimeout)
	d.Set("upstream_send_timeout", api.UpstreamSendTimeout)
	d.Set("upstream_read_timeout", api.UpstreamReadTimeout)
}
