package kong

import (
	"fmt"
	"net/http"

	"github.com/WeKnowSports/terraform-provider-kong/helper"
	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Route : Kong Route request object structure
type Route struct {
	ID                      string              `json:"id,omitempty"`
	Name                    string              `json:"name,omitempty"`
	Protocols               []string            `json:"protocols"`
	Methods                 []string            `json:"methods"`
	Hosts                   []string            `json:"hosts"`
	Paths                   []string            `json:"paths"`
	Headers                 map[string][]string `json:"headers"`
	HttpsRedirectStatusCode int                 `json:"https_redirect_status_code,omitempty"`
	RegexPriority           int                 `json:"regex_priority"`
	StripPath               bool                `json:"strip_path,omitempty"`
	PathHandling            string              `json:"path_handling,omitempty"`
	PreserveHost            bool                `json:"preserve_host,omitempty"`
	RequestBuffering        bool                `json:"request_buffering"`
	ResponseBuffering       bool                `json:"response_buffering"`
	SNIs                    []string            `json:"snis,omitempty"`
	// Sources                 []string            `json:"sources,omitempty"`
	// Destinations            []string            `json:"destinations,omitempty"`
	Tags    []string `json:"tags"`
	Service Service  `json:"service,omitempty"`
}

func resourceKongRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongRouteCreate,
		Read:   resourceKongRouteRead,
		Update: resourceKongRouteUpdate,
		Delete: resourceKongRouteDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The name of the Route.",
			},

			"protocols": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "A list of the protocols this Route should allow. By default it is [\"http\", \"https\"], which means that the Route accepts both. When set to [\"https\"], HTTP requests are answered with a request to upgrade to HTTPS.",
			},

			"methods": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of HTTP methods that match this Route. For example: [\"GET\", \"POST\"]. At least one of hosts, paths, or methods must be set.",
			},

			"hosts": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of domain names that match this Route. For example: example.com. At least one of hosts, paths, or methods must be set.",
			},

			"paths": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "A list of paths that match this Route. For example: /my-path. At least one of hosts, paths, or methods must be set.",
			},

			"header": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				Description: "One or more lists of values indexed by header name that will cause this Route to match if present in the request.",
			},

			"https_redirect_status_code": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     426,
				Description: "The status code Kong responds with when all properties of a Route match except the protocol i.e. if the protocol of the request is HTTP instead of HTTPS",
			},

			"regex_priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "A number used to choose which route resolves a given request when several routes match it using regexes simultaneously.",
			},
			"strip_path": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When matching a Route via one of the paths, strip the matching prefix from the upstream request URL. Defaults to true.",
			},

			"path_handling": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "v0",
				Description: "Controls how the Service path, Route path and requested path are combined when sending a request to the upstream",
			},

			"preserve_host": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When matching a Route via one of the hosts domain names, use the request Host header in the upstream request headers. By default set to false, and the upstream Host header will be that of the Service's host.",
			},

			"request_buffering": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable request body buffering or not",
			},

			"response_buffering": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable response body buffering or not",
			},

			"snis": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "A list of SNIs that match this Route when using stream routing.",
			},

			// "sources": {
			// 	Type: schema.TypeSet,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"ip": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 			"port": {
			// 				Type:     schema.TypeInt,
			// 				Optional: true,
			// 			},
			// 		},
			// 	},
			// 	Optional:    true,
			// 	Description: "A list of IP sources of incoming connections that match this Route when using stream routing",
			// },

			// "destinations": {
			// 	Type: schema.TypeSet,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"ip": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 			"port": {
			// 				Type:     schema.TypeInt,
			// 				Optional: true,
			// 			},
			// 		},
			// 	},
			// 	Optional:    true,
			// 	Description: "A list of IP destinations of incoming connections that match this Route when using stream routing",
			// },

			"tags": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An optional set of strings associated with the Service for grouping and filtering.",
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

	id := d.Id()
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

	id := d.Id()

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
		ID:                      d.Id(),
		Name:                    d.Get("name").(string),
		Protocols:               helper.ConvertInterfaceArrToStrings(d.Get("protocols").([]interface{})),
		Methods:                 helper.ConvertInterfaceArrToStrings(d.Get("methods").([]interface{})),
		Hosts:                   helper.ConvertInterfaceArrToStrings(d.Get("hosts").([]interface{})),
		Paths:                   helper.ConvertInterfaceArrToStrings(d.Get("paths").([]interface{})),
		Headers:                 readMapStringArrayFromResource(d, "header"),
		HttpsRedirectStatusCode: d.Get("https_redirect_status_code").(int),
		RegexPriority:           d.Get("regex_priority").(int),
		StripPath:               d.Get("strip_path").(bool),
		PathHandling:            d.Get("path_handling").(string),
		PreserveHost:            d.Get("preserve_host").(bool),
		RequestBuffering:        d.Get("request_buffering").(bool),
		ResponseBuffering:       d.Get("response_buffering").(bool),
		SNIs:                    helper.ConvertInterfaceArrToStrings(d.Get("snis").([]interface{})),
		// Sources:                 helper.ConvertInterfaceArrToStrings(d.Get("sources").([]interface{})),
		// Destinations:            helper.ConvertInterfaceArrToStrings(d.Get("destinations").([]interface{})),
		Tags: helper.ConvertInterfaceArrToStrings(d.Get("tags").([]interface{})),
		Service: Service{
			ID: d.Get("service").(string),
		},
	}

	return route
}

func setRouteToResourceData(d *schema.ResourceData, route *Route) {
	d.SetId(route.ID)
	d.Set("name", route.Name)
	d.Set("protocols", route.Protocols)
	d.Set("methods", route.Methods)
	d.Set("hosts", route.Hosts)
	d.Set("paths", route.Paths)
	d.Set("headers", route.Headers)
	d.Set("https_redirect_status_code", route.HttpsRedirectStatusCode)
	d.Set("regex_priority", route.RegexPriority)
	d.Set("strip_path", route.StripPath)
	d.Set("path_handling", route.PathHandling)
	d.Set("preserve_host", route.PreserveHost)
	d.Set("request_buffering", route.RequestBuffering)
	d.Set("response_buffering", route.ResponseBuffering)
	d.Set("snis", route.SNIs)
	// d.Set("sources", route.Sources)
	// d.Set("destinations", route.Destinations)
	d.Set("tags", route.Tags)
	d.Set("service", route.Service.ID)
}

func readMapStringArrayFromResource(d *schema.ResourceData, key string) map[string][]string {
	results := map[string][]string{}
	if attr, ok := d.GetOk(key); ok {
		set := attr.(*schema.Set)
		for _, item := range set.List() {
			m := item.(map[string]interface{})
			if name, ok := m["name"].(string); ok {
				if values, ok := m["values"].([]interface{}); ok {
					var vals []string
					for _, v := range values {
						vals = append(vals, v.(string))
					}
					results[name] = vals
				}
			}
		}
	}

	return results
}
