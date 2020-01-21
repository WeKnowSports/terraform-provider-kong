package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

type Upstream struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	Slots              int    `json:"slots,omitempty"`
	HashOn             string `json:"hash_on,omitempty"`
	HashFallback       string `json:"hash_fallback,omitempty"`
	HashOnHeader       string `json:"hash_on_header,omitempty"`
	HashFallbackHeader string `json:"hash_fallback_header,omitempty"`
	HashOnCookie       string `json:"hash_on_cookie,omitempty"`
	HashOnCookiePath   string `json:"hash_on_cookie_path,omitempty"`
	Algorithm          string `json:"algorithm,omitempty"`
}

func resourceKongUpstream() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongUpstreamCreate,
		Read:   resourceKongUpstreamRead,
		Update: resourceKongUpstreamUpdate,
		Delete: resourceKongUpstreamDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is a hostname, which must be equal to the host of a Service.",
			},
			"slots": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of slots in the loadbalancer algorithm (10-65536, defaults to 1000).",
				Default:     1000,
				ValidateFunc: func(i interface{}, s string) (strings []string, errors []error) {
					slots := i.(int)

					if slots >= 10 && slots <= 65536 {
						return nil, nil
					}

					return nil, []error{fmt.Errorf("slots value of %d not in the range of 10-65536", slots)}
				},
			},
			"hash_on": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "What to use as hashing input: none, consumer, ip, header, or cookie (defaults to none resulting in a weighted-round-robin scheme).",
				Default:     "none",
				ValidateFunc: func(i interface{}, s string) (strings []string, errors []error) {
					// TODO: validate against [none, consume, ip, header, cookie]
					return nil, nil
				},
			},
			"hash_fallback": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "What to use as hashing input if the primary hash_on does not return a hash (eg. header is missing, or no consumer identified). One of: none, consumer, ip, header, or cookie (defaults to none, not available if hash_on is set to cookie).",
				Default:     "none",
				ValidateFunc: func(i interface{}, s string) (strings []string, errors []error) {
					// TODO: validate against [none, consume, ip, header, cookie]
					return nil, nil
				},
			},
			"hash_on_header": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The header name to take the value from as hash input (only required when hash_on is set to header).",
			},
			"hash_fallback_header": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The header name to take the value from as hash input (only required when hash_fallback is set to header).",
			},
			"hash_on_cookie": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cookie name to take the value from as hash input (only required when hash_on or hash_fallback is set to cookie). If the specified cookie is not in the request, Kong will generate a value and set the cookie in the response.",
			},
			"hash_on_cookie_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return (old == "" && new == "/") || (old == "/" && new == "")
				},
				Description: "The cookie path to set in the response headers (only required when hash_on or hash_fallback is set to cookie, defaults to \"/\")",
			},
			"algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Which load balancing algorithm to use. One of: round-robin, consistent-hashing, or least-connections. Defaults to \"round-robin\". Kong 1.3.0 and up.",
				ValidateFunc: func(i interface{}, s string) (strings []string, errors []error) {
					a := i.(string)
					algs := []string{"round-robin", "consistent-hashing", "least-connections"}
					for i := 0; i < len(algs); i++ {
						if algs[i] == a {
							return nil, nil
						}
					}

					return nil, append(errors, fmt.Errorf("algorithm must be one of %v. %s was provided instead.", algs, s))
				},
			},
		},
	}
}

func resourceKongUpstreamCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	upstream := getUpstreamFromResourceData(d)

	createdUpstream := getUpstreamFromResourceData(d)

	response, error := sling.New().BodyJSON(upstream).Post("upstreams/").ReceiveSuccess(createdUpstream)
	if error != nil {
		return fmt.Errorf("Error while creating upstream.")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setUpstreamToResourceData(d, createdUpstream)

	return nil
}

func resourceKongUpstreamRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	upstream := getUpstreamFromResourceData(d)

	response, error := sling.New().Path("upstreams/").Get(upstream.ID).ReceiveSuccess(upstream)
	if error != nil {
		return fmt.Errorf("Error while updating upstream.")
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setUpstreamToResourceData(d, upstream)

	return nil
}

func resourceKongUpstreamUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	upstream := getUpstreamFromResourceData(d)

	updatedUpstream := getUpstreamFromResourceData(d)

	response, error := sling.New().BodyJSON(upstream).Path("upstreams/").Patch(upstream.ID).ReceiveSuccess(updatedUpstream)
	if error != nil {
		return fmt.Errorf("Error while updating upstream")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setUpstreamToResourceData(d, updatedUpstream)

	return nil
}

func resourceKongUpstreamDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	upstream := getUpstreamFromResourceData(d)

	response, error := sling.New().Path("upstreams/").Delete(upstream.ID).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("Error while deleting upstream.")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getUpstreamFromResourceData(d *schema.ResourceData) *Upstream {
	upstream := &Upstream{
		ID:                 d.Id(),
		Name:               d.Get("name").(string),
		Slots:              d.Get("slots").(int),
		HashOn:             d.Get("hash_on").(string),
		HashFallback:       d.Get("hash_fallback").(string),
		HashOnHeader:       d.Get("hash_on_header").(string),
		HashFallbackHeader: d.Get("hash_fallback_header").(string),
		HashOnCookie:       d.Get("hash_on_cookie").(string),
		HashOnCookiePath:   d.Get("hash_on_cookie_path").(string),
		Algorithm:          d.Get("algorithm").(string),
	}

	return upstream
}

func setUpstreamToResourceData(d *schema.ResourceData, upstream *Upstream) {
	d.SetId(upstream.ID)
	d.Set("name", upstream.Name)
	d.Set("slots", upstream.Slots)
	d.Set("hash_on", upstream.HashOn)
	d.Set("hash_fallback", upstream.HashFallback)
	d.Set("hash_on_header", upstream.HashOnHeader)
	d.Set("hash_fallback_header", upstream.HashFallbackHeader)
	d.Set("hash_on_cookie", upstream.HashOnCookie)
	d.Set("algorithm", upstream.Algorithm)
}
