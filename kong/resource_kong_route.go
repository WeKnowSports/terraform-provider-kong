package kong

import (
	"fmt"
	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
)

// Route : Kong Route request object structure
type Route struct {
	ID           string   `json:"id,omitempty"`
	Protocols    []string `json:"protocols,omitempty"`
	Methods      []string `json:"methods"`
	Hosts        []string `json:"hosts"`
	Paths        []string `json:"paths"`
	StripPath    bool     `json:"strip_path,omitempty"`
	PreserveHost bool     `json:"preserve_host,omitempty"`
	Service      Service  `json:"service,omitempty"`
}

func resourceKongRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongRouteCreate,
		Read:   resourceKongRouteRead,
		Update: resourceKongRouteUpdate,
		Delete: resourceKongRouteDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"protocols": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					protocols := convertInterfaceArrToStrings(d.Get("protocols").([]interface{}))

					// TODO: Yeah...
					return len(protocols) == 2 &&
						((protocols[0] == "http" && protocols[1] == "https") ||
						(protocols[0] == "https" && protocols[1] == "http"))
				},
				Description: "A list of the protocols this Route should allow. By default it is [\"http\", \"https\"], which means that the Route accepts both. When set to [\"https\"], HTTP requests are answered with a request to upgrade to HTTPS.",
			},

			"methods": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "A list of HTTP methods that match this Route. For example: [\"GET\", \"POST\"]. At least one of hosts, paths, or methods must be set.",
			},

			"hosts": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "A list of domain names that match this Route. For example: example.com. At least one of hosts, paths, or methods must be set.",
			},

			"paths": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "A list of paths that match this Route. For example: /my-path. At least one of hosts, paths, or methods must be set.",
			},

			"strip_path": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When matching a Route via one of the paths, strip the matching prefix from the upstream request URL. Defaults to true.",
			},

			"preserve_host": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When matching a Route via one of the hosts domain names, use the request Host header in the upstream request headers. By default set to false, and the upstream Host header will be that of the Service's host.",
			},

			"connect_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60000,
				Description: "The timeout in milliseconds for establishing a connection to the upstream server. Defaults to 60000.",
			},

			"service": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Service this Route is associated to. This is where the Route proxies traffic to.",
			},
		},
	}
}

func resourceKongRouteCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	route := getRouteFromResourceData(d)

	createdRoute := new(Route)
	response, error := sling.New().BodyJSON(route).Post("routes/").ReceiveSuccess(createdRoute)

	if error != nil {
		return fmt.Errorf("error while creating Route: " + error.Error())
	}

	if response.StatusCode == http.StatusConflict {
		return fmt.Errorf("409 Conflict - use terraform import to manage this route")
	} else if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	setRouteToResourceData(d, createdRoute)

	return nil
}

func resourceKongRouteRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	id := d.Get("id").(string)
	route := new(Route)

	response, error := sling.New().Path("routes/").Get(id).ReceiveSuccess(route)

	if error != nil {
		return fmt.Errorf("error while updating Route" + error.Error())
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	setRouteToResourceData(d, route)

	return nil
}

func resourceKongRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	route := getRouteFromResourceData(d)

	updatedRoute := new(Route)

	response, error := sling.New().BodyJSON(route).Patch("routes/").Path(route.ID).ReceiveSuccess(updatedRoute)

	if error != nil {
		return fmt.Errorf("error while updating Route" + error.Error())
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	setRouteToResourceData(d, updatedRoute)

	return nil
}

func resourceKongRouteDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	id := d.Get("id").(string)

	response, error := sling.New().Delete("routes/").Path(id).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("error while deleting Route" + error.Error())
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	return nil
}

func getRouteFromResourceData(d *schema.ResourceData) *Route {
	route := &Route{
		Protocols:    convertInterfaceArrToStrings(d.Get("protocols").([]interface{})),
		Methods:      convertInterfaceArrToStrings(d.Get("methods").([]interface{})),
		Hosts:        convertInterfaceArrToStrings(d.Get("hosts").([]interface{})),
		Paths:        convertInterfaceArrToStrings(d.Get("paths").([]interface{})),
		StripPath:    d.Get("strip_path").(bool),
		PreserveHost: d.Get("preserve_host").(bool),
		Service: Service{
			ID: d.Get("service").(string),
		},
	}

	if id, ok := d.GetOk("id"); ok {
		route.ID = id.(string)
	}

	return route
}

func setRouteToResourceData(d *schema.ResourceData, route *Route) {
	d.SetId(route.ID)
	d.Set("protocols", route.Protocols)
	d.Set("methods", route.Methods)
	d.Set("hosts", route.Hosts)
	d.Set("paths", route.Paths)
	d.Set("strip_path", route.StripPath)
	d.Set("preserve_host", route.PreserveHost)
	d.Set("service", route.Service.ID)
}

func convertInterfaceArrToStrings(strs []interface{}) []string {
	arr := make([]string, len(strs))
	for i, str := range strs {
		arr[i] = str.(string)
	}
	return arr
}
