package kong

import (
	"fmt"
	"net/http"

	"github.com/WeKnowSports/terraform-provider-kong/helper"
	"github.com/dghubble/sling"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	HealthchecksTypes = []string{"http", "tcp", "https"}
)

type PassiveHealthy struct {
	Successes    int   `json:"successes"`
	HttpStatuses []int `json:"http_statuses,omitempty"`
}

type PassiveUnhealthy struct {
	HttpFailures int   `json:"http_failures"`
	HttpStatuses []int `json:"http_statuses,omitempty"`
	TcpFailures  int   `json:"tcp_failures"`
	Timeouts     int   `json:"timeouts"`
}

type HealthChecksPassive struct {
	Type      string            `json:"type,omitempty"`
	Healthy   *PassiveHealthy   `json:"healthy,omitempty"`
	Unhealthy *PassiveUnhealthy `json:"unhealthy,omitempty"`
}

type ActiveHealthy struct {
	Successes    int   `json:"successes"`
	Interval     int   `json:"interval"`
	HttpStatuses []int `json:"http_statuses,omitempty"`
}

type ActiveUnhealthy struct {
	HttpStatuses []int `json:"http_statuses,omitempty"`
	TcpFailures  int   `json:"tcp_failures"`
	Timeouts     int   `json:"timeouts"`
	HttpFailures int   `json:"http_failures"`
	Interval     int   `json:"interval"`
}

type HealthChecksActive struct {
	HttpsVerifyCertificate bool             `json:"https_verify_certificate"`
	HttpPath               string           `json:"http_path,omitempty"`
	Timeout                int              `json:"timeout,omitempty"`
	HttpsSni               *string          `json:"https_sni,omitempty"`
	Concurrency            int              `json:"concurrency,omitempty"`
	Type                   string           `json:"type,omitempty"`
	Healthy                *ActiveHealthy   `json:"healthy,omitempty"`
	Unhealthy              *ActiveUnhealthy `json:"unhealthy,omitempty"`
}

type UpstreamHealthChecks struct {
	Active  *HealthChecksActive  `json:"active,omitempty"`
	Passive *HealthChecksPassive `json:"passive,omitempty"`
}

type Upstream struct {
	ID                      string                `json:"id,omitempty"`
	Name                    string                `json:"name,omitempty"`
	Algorithm               string                `json:"algorithm,omitempty"`
	HashOn                  string                `json:"hash_on"`
	HashFallback            string                `json:"hash_fallback"`
	HashOnHeader            string                `json:"hash_on_header,omitempty"`
	HashFallbackHeader      string                `json:"hash_fallback_header,omitempty"`
	HashOnCookie            string                `json:"hash_on_cookie,omitempty"`
	HashOnCookiePath        string                `json:"hash_on_cookie_path,omitempty"`
	HashOnQueryArg          string                `json:"hash_on_query_arg,omitempty"`
	HashFallbackOnQueryArg  string                `json:"hash_fallback_query_arg,omitempty"`
	HashOnUriCapture        string                `json:"hash_on_uri_capture,omitempty"`
	HashFallbacOnUriCapture string                `json:"hash_fallback_uri_capture,omitempty"`
	Slots                   int                   `json:"slots,omitempty"`
	HealthChecks            *UpstreamHealthChecks `json:"healthchecks,omitempty"`
	Tags                    []string              `json:"tags"`
	HostHeader              string                `json:"host_header,omitempty"`
	ClientCertificate       Certificate           `json:"client_certificate,omitempty"`
	UseSrvName              bool                  `json:"use_srv_name,omitempty"`
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
			"algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Which load balancing algorithm to use. One of: round-robin, consistent-hashing, or least-connections. Defaults to \"round-robin\". Kong 1.3.0 and up.",
				Default:     "round-robin",
			},
			"hash_on": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "What to use as hashing input: none, consumer, ip, header, or cookie (defaults to none resulting in a weighted-round-robin scheme).",
				Default:     "none",
			},
			"hash_fallback": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "What to use as hashing input if the primary hash_on does not return a hash (eg. header is missing, or no consumer identified). One of: none, consumer, ip, header, or cookie (defaults to none, not available if hash_on is set to cookie).",
				Default:     "none",
			},
			"hash_on_header": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The header name to take the value from as hash input (only required when hash_on is set to header).",
			},
			"hash_fallback_header": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The header name to take the value from as hash input (only required when hash_fallback is set to header).",
			},
			"hash_on_cookie": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The cookie name to take the value from as hash input (only required when hash_on or hash_fallback is set to cookie). If the specified cookie is not in the request, Kong will generate a value and set the cookie in the response.",
			},
			"hash_on_cookie_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/",
				Description: "The cookie path to set in the response headers (only required when hash_on or hash_fallback is set to cookie, defaults to \"/\")",
			},
			"hash_on_query_arg": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the query string argument to take the value from as hash input. Only required when hash_on is set to query_arg",
			},
			"hash_fallback_query_arg": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the query string argument to take the value from as hash input. Only required when hash_fallback is set to query_arg",
			},
			"hash_on_uri_capture": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the route URI capture to take the value from as hash input. Only required when hash_on is set to uri_capture",
			},
			"hash_fallback_uri_capture": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the route URI capture to take the value from as hash input. Only required when hash_fallback is set to uri_capture",
			},
			"slots": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"healthchecks": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										// Kong 1.0.0+
										// Default:  "http",
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  1,
									},
									"concurrency": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  10,
									},
									"http_path": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "/",
									},
									"https_verify_certificate": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
										// Kong 1.0.0+
										// Default:  true,
									},
									"https_sni": {
										Type:     schema.TypeString,
										Optional: true,
										// Kong 1.0.0+
										// Default:  nil,
									},
									"healthy": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"interval": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"http_statuses": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"successes": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"unhealthy": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"interval": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"http_statuses": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"tcp_failures": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"http_failures": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"timeouts": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"passive": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										// Kong 1.0.0+
										// Default:  "http",
									},
									"healthy": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"http_statuses": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"successes": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"unhealthy": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"http_statuses": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"tcp_failures": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"http_failures": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"timeouts": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An optional set of strings associated with the Service for grouping and filtering.",
			},
			"client_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"use_srv_name": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"host_header": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceKongUpstreamCreate(d *schema.ResourceData, meta interface{}) error {
	Sling := meta.(*sling.Sling)

	upstream := getUpstreamFromResourceData(d)

	createdUpstream := getUpstreamFromResourceData(d)

	response, Error := Sling.New().BodyJSON(upstream).Post("upstreams/").ReceiveSuccess(createdUpstream)
	if Error != nil {
		return fmt.Errorf("error while creating upstream")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setUpstreamToResourceData(d, createdUpstream)

	return nil
}

func resourceKongUpstreamRead(d *schema.ResourceData, meta interface{}) error {
	Sling := meta.(*sling.Sling)

	upstream := getUpstreamFromResourceData(d)

	response, Error := Sling.New().Path("upstreams/").Get(upstream.ID).ReceiveSuccess(upstream)
	if Error != nil {
		return fmt.Errorf(Error.Error()) //fmt.Errorf("Error while updating upstream")
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
	Sling := meta.(*sling.Sling)

	upstream := getUpstreamFromResourceData(d)
	updatedUpstream := getUpstreamFromResourceData(d)

	response, Error := Sling.New().BodyJSON(upstream).Path("upstreams/").Patch(upstream.ID).ReceiveSuccess(updatedUpstream)
	if Error != nil {
		return fmt.Errorf(Error.Error()) //fmt.Errorf("Error while updating upstream")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setUpstreamToResourceData(d, updatedUpstream)

	return nil
}

func resourceKongUpstreamDelete(d *schema.ResourceData, meta interface{}) error {
	Sling := meta.(*sling.Sling)

	upstream := getUpstreamFromResourceData(d)

	response, Error := Sling.New().Path("upstreams/").Delete(upstream.ID).ReceiveSuccess(nil)
	if Error != nil {
		return fmt.Errorf("error while deleting upstream")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getActiveHealthyFromMap(d *map[string]interface{}) *ActiveHealthy {
	if d != nil {
		m := *d
		healthy := &ActiveHealthy{}

		if m["interval"] != nil {
			healthy.Interval = m["interval"].(int)
		}
		if m["http_statuses"] != nil {
			healthy.HttpStatuses = readIntArrayFromInterface(m["http_statuses"])
		}
		if m["successes"] != nil {
			healthy.Successes = m["successes"].(int)
		}

		return healthy
	}
	return nil
}

func getActiveUnhealthyFromMap(d *map[string]interface{}) *ActiveUnhealthy {
	if d != nil {
		m := *d
		unhealthy := &ActiveUnhealthy{}

		if m["interval"] != nil {
			unhealthy.Interval = m["interval"].(int)
		}
		if m["http_statuses"] != nil {
			unhealthy.HttpStatuses = readIntArrayFromInterface(m["http_statuses"])
		}
		if m["tcp_failures"] != nil {
			unhealthy.TcpFailures = m["tcp_failures"].(int)
		}
		if m["http_failures"] != nil {
			unhealthy.HttpFailures = m["http_failures"].(int)
		}
		if m["timeouts"] != nil {
			unhealthy.Timeouts = m["timeouts"].(int)
		}

		return unhealthy
	}
	return nil
}

func getActiveHealthChecksFromMap(d *map[string]interface{}) *HealthChecksActive {
	if d != nil {
		m := *d
		active := &HealthChecksActive{}

		if m["type"] != nil {
			active.Type = m["type"].(string)
		}
		if m["timeout"] != nil {
			active.Timeout = m["timeout"].(int)
		}
		if m["concurrency"] != nil {
			active.Concurrency = m["concurrency"].(int)
		}
		if m["http_path"] != nil {
			active.HttpPath = m["http_path"].(string)
		}
		if m["https_verify_certificate"] != nil {
			active.HttpsVerifyCertificate = m["https_verify_certificate"].(bool)
		}
		if m["https_sni"] != nil {
			if len(m["https_sni"].(string)) != 0 {
				httpsSni := m["https_sni"].(string)
				active.HttpsSni = &httpsSni
			}
		}

		if m["healthy"] != nil {
			if healthyArray := m["healthy"].([]interface{}); len(healthyArray) > 0 {
				healthyMap := healthyArray[0].(map[string]interface{})
				active.Healthy = getActiveHealthyFromMap(&healthyMap)
			}
		}

		if m["unhealthy"] != nil {
			if unhealthyArray := m["unhealthy"].([]interface{}); len(unhealthyArray) > 0 {
				unhealthyMap := unhealthyArray[0].(map[string]interface{})
				active.Unhealthy = getActiveUnhealthyFromMap(&unhealthyMap)
			}
		}

		return active
	}

	return nil
}

func getPassiveHealthyFromMap(d *map[string]interface{}) *PassiveHealthy {
	if d != nil {
		m := *d
		healthy := &PassiveHealthy{}

		if m["http_statuses"] != nil {
			healthy.HttpStatuses = readIntArrayFromInterface(m["http_statuses"])
		}
		if m["successes"] != nil {
			healthy.Successes = m["successes"].(int)
		}

		return healthy
	}
	return nil
}

func getPassiveUnhealthyFromMap(d *map[string]interface{}) *PassiveUnhealthy {
	if d != nil {
		m := *d
		unhealthy := &PassiveUnhealthy{}

		if m["http_statuses"] != nil {
			unhealthy.HttpStatuses = readIntArrayFromInterface(m["http_statuses"])
		}
		if m["tcp_failures"] != nil {
			unhealthy.TcpFailures = m["tcp_failures"].(int)
		}
		if m["http_failures"] != nil {
			unhealthy.HttpFailures = m["http_failures"].(int)
		}
		if m["timeouts"] != nil {
			unhealthy.Timeouts = m["timeouts"].(int)
		}

		return unhealthy
	}
	return nil
}

func getPassiveHealthsCheckFromMap(d *map[string]interface{}) *HealthChecksPassive {
	if d != nil {
		m := *d
		passive := &HealthChecksPassive{}

		if m["type"] != nil {
			passive.Type = m["type"].(string)
		}

		if m["healthy"] != nil {
			if healthyArray := m["healthy"].([]interface{}); len(healthyArray) > 0 {
				healthyMap := healthyArray[0].(map[string]interface{})
				passive.Healthy = getPassiveHealthyFromMap(&healthyMap)
			}
		}

		if m["unhealthy"] != nil {
			if unhealthyArray := m["unhealthy"].([]interface{}); len(unhealthyArray) > 0 {
				unhealthyMap := unhealthyArray[0].(map[string]interface{})
				passive.Unhealthy = getPassiveUnhealthyFromMap(&unhealthyMap)
			}
		}

		return passive
	}

	return nil
}

func getHealthChecksFromMap(d *map[string]interface{}) *UpstreamHealthChecks {
	if d != nil {
		m := *d
		healthChecks := &UpstreamHealthChecks{}

		if m["active"] != nil {
			if activeArray := m["active"].([]interface{}); len(activeArray) > 0 {
				activeMap := activeArray[0].(map[string]interface{})
				healthChecks.Active = getActiveHealthChecksFromMap(&activeMap)
			}
		}

		if m["passive"] != nil {
			if passiveArray := m["passive"].([]interface{}); len(passiveArray) > 0 {
				passiveMap := passiveArray[0].(map[string]interface{})
				healthChecks.Passive = getPassiveHealthsCheckFromMap(&passiveMap)
			}
		}

		return healthChecks
	}

	return nil
}

func getUpstreamFromResourceData(d *schema.ResourceData) *Upstream {
	upstream := &Upstream{
		ID:                      d.Id(),
		Name:                    d.Get("name").(string),
		Algorithm:               d.Get("algorithm").(string),
		HashOn:                  d.Get("hash_on").(string),
		HashFallback:            d.Get("hash_fallback").(string),
		HashOnHeader:            d.Get("hash_on_header").(string),
		HashFallbackHeader:      d.Get("hash_fallback_header").(string),
		HashOnCookie:            d.Get("hash_on_cookie").(string),
		HashOnCookiePath:        d.Get("hash_on_cookie_path").(string),
		HashOnQueryArg:          d.Get("hash_on_query_arg").(string),
		HashFallbackOnQueryArg:  d.Get("hash_fallback_query_arg").(string),
		HashOnUriCapture:        d.Get("hash_on_uri_capture").(string),
		HashFallbacOnUriCapture: d.Get("hash_fallback_uri_capture").(string),
		Slots:                   d.Get("slots").(int),
		Tags:                    helper.ConvertInterfaceArrToStrings(d.Get("tags").([]interface{})),
		HostHeader:              d.Get("host_header").(string),
		ClientCertificate: Certificate{
			ID: d.Get("client_certificate").(string),
		},
		UseSrvName: d.Get("use_srv_name").(bool),
	}

	hcArr := d.Get("healthchecks").([]interface{})

	if len(hcArr) > 0 {
		hcMap := hcArr[0].(map[string]interface{})
		upstream.HealthChecks = getHealthChecksFromMap(&hcMap)
	}

	return upstream
}

func convertActiveHealthyToResourceData(ah *ActiveHealthy) []map[string]interface{} {
	if ah == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["interval"] = ah.Interval
	m["http_statuses"] = ah.HttpStatuses
	m["successes"] = ah.Successes

	return []map[string]interface{}{m}
}

func convertActiveUnhealthyToResource(au *ActiveUnhealthy) []map[string]interface{} {
	if au == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["interval"] = au.Interval
	m["http_statuses"] = au.HttpStatuses
	m["tcp_failures"] = au.TcpFailures
	m["http_failures"] = au.HttpFailures
	m["timeouts"] = au.Timeouts

	return []map[string]interface{}{m}
}

func convertHealthCheckActiveToResourceData(hca *HealthChecksActive) []interface{} {
	if hca == nil {
		return []interface{}{}
	}
	m := make(map[string]interface{})

	m["type"] = hca.Type
	m["timeout"] = hca.Timeout
	m["concurrency"] = hca.Concurrency
	m["http_path"] = hca.HttpPath
	m["https_verify_certificate"] = hca.HttpsVerifyCertificate

	if hca.HttpsSni != nil {
		m["https_sni"] = *hca.HttpsSni
	}
	if hca.Healthy != nil {
		m["healthy"] = convertActiveHealthyToResourceData(hca.Healthy)
	}
	if hca.Unhealthy != nil {
		m["unhealthy"] = convertActiveUnhealthyToResource(hca.Unhealthy)
	}

	return []interface{}{m}
}

func convertPassiveHealthyToResourceData(ph *PassiveHealthy) []map[string]interface{} {
	if ph == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["http_statuses"] = ph.HttpStatuses
	m["successes"] = ph.Successes

	return []map[string]interface{}{m}
}

func convertPassiveUnhealthyToResourceData(pu *PassiveUnhealthy) []map[string]interface{} {
	if pu == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["http_statuses"] = pu.HttpStatuses
	m["tcp_failures"] = pu.TcpFailures
	m["http_failures"] = pu.HttpFailures
	m["timeouts"] = pu.Timeouts

	return []map[string]interface{}{m}
}

func convertHealthCheckPassiveToResourceData(in *HealthChecksPassive) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	m["type"] = in.Type

	if in.Healthy != nil {
		m["healthy"] = convertPassiveHealthyToResourceData(in.Healthy)
	}
	if in.Unhealthy != nil {
		m["unhealthy"] = convertPassiveUnhealthyToResourceData(in.Unhealthy)
	}

	return []interface{}{m}
}

func convertHealthCheckResourceData(uhc *UpstreamHealthChecks) []interface{} {
	if uhc == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if uhc.Active != nil {
		m["active"] = convertHealthCheckActiveToResourceData(uhc.Active)
	}
	if uhc.Passive != nil {
		m["passive"] = convertHealthCheckPassiveToResourceData(uhc.Passive)
	}

	return []interface{}{m}
}

func setUpstreamToResourceData(d *schema.ResourceData, upstream *Upstream) {
	d.SetId(upstream.ID)
	d.Set("name", upstream.Name)
	d.Set("algorithm", upstream.Algorithm)
	d.Set("hash_on", upstream.HashOn)
	d.Set("hash_fallback", upstream.HashFallback)
	d.Set("hash_on_header", upstream.HashOnHeader)
	d.Set("hash_fallback_header", upstream.HashFallbackHeader)
	d.Set("hash_on_cookie", upstream.HashOnCookie)
	d.Set("hash_on_query_arg", upstream.HashOnQueryArg)
	d.Set("hash_fallback_query_arg", upstream.HashFallbackOnQueryArg)
	d.Set("hash_on_uri_capture", upstream.HashOnUriCapture)
	d.Set("hash_fallback_uri_capture", upstream.HashFallbacOnUriCapture)
	d.Set("slots", upstream.Slots)
	d.Set("healthchecks", convertHealthCheckResourceData(upstream.HealthChecks))
	d.Set("tags", upstream.Tags)
	d.Set("host_header", upstream.HostHeader)
	d.Set("client_certificate", upstream.ClientCertificate)
	d.Set("use_srv_name", upstream.UseSrvName)
}

func readIntArrayFromInterface(in interface{}) []int {
	if arr := in.([]interface{}); arr != nil {
		array := make([]int, len(arr))
		for i, x := range arr {
			item := x.(int)
			array[i] = item
		}

		return array
	}

	return []int{}
}
