package kong

import (
	"fmt"
	"net/http"

	"github.com/WeKnowSports/terraform-provider-kong/helper"
	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Service : Kong Service request object structure
type Service struct {
	ID                string      `json:"id,omitempty"`
	Name              string      `json:"name,omitempty"`
	Retries           int         `json:"retries,omitempty"`
	Protocol          string      `json:"protocol,omitempty"`
	Host              string      `json:"host,omitempty"`
	Port              int         `json:"port,omitempty"`
	Path              string      `json:"path,omitempty"`
	ConnectTimeout    int         `json:"connect_timeout,omitempty"`
	WriteTimeout      int         `json:"write_timeout,omitempty"`
	ReadTimeout       int         `json:"read_timeout,omitempty"`
	Tags              []string    `json:"tags"`
	ClientCertificate Certificate `json:"-"`                          // TO DO: add if statement which assign value only if Protocol is HTTPS
	TlsVerify         bool        `json:"tls_verify,omitempty"`       //
	TlsVerifyDepth    int         `json:"tls_verify_depth,omitempty"` //
	CACertificates    []string    `json:"-"`                          //
	Enabled           bool        `json:"enabled"`
}

func resourceKongService() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongServiceCreate,
		Read:   resourceKongServiceRead,
		Update: resourceKongServiceUpdate,
		Delete: resourceKongServiceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The Service name.",
			},

			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The protocol used to communicate with the upstream. It can be one of http (default) or https.",
				Default:     "http",
			},

			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The host of the upstream server.",
			},

			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The upstream server port. Defaults to 80.",
				Default:     80,
			},

			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to be used in requests to the upstream server. Empty by default.",
				Default:     "/",
			},

			"retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of retries to execute upon failure to proxy. Default: 5.",
				Default:     5,
			},

			"connect_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout in milliseconds for establishing a connection to the upstream server. Defaults to 60000.",
				Default:     60000,
			},

			"write_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout in milliseconds between two successive write operations for transmitting a request to the upstream server. Defaults to 60000.",
				Default:     60000,
			},

			"read_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout in milliseconds between two successive read operations for transmitting a request to the upstream server. Defaults to 60000.",
				Default:     60000,
			},

			"tags": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An optional set of strings associated with the Service for grouping and filtering.",
			},

			"client_certificate": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Certificate to be used as client certificate while TLS handshaking to the upstream server",
				RequiredWith: []string{"protocol"},
				Default:      nil,
			},

			"tls_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable verification of upstream server TLS certificate. If set to null, then the Nginx default is respected",
				Default:     nil,
			},

			"tls_verify_depth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum depth of chain while verifying Upstream servers TLS certificate. If set to null, then the Nginx default is respected",
				Default:     nil,
			},

			"ca_certificates": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Array of CA Certificate object UUIDs that are used to build the trust store while verifying upstream servers TLS certificate",
				Default:     nil,
			},

			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the Service is active",
				Default:     true,
			},
		},
	}
}

func resourceKongServiceCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*sling.Sling)

	service := getServiceFromResourceData(d)

	createdService := new(Service)
	response, e := s.New().BodyJSON(service).Post("services/").ReceiveSuccess(createdService)

	if e != nil {
		return fmt.Errorf("error while creating Service: " + e.Error())
	}

	if response.StatusCode == http.StatusConflict {
		return fmt.Errorf("409 Conflict - use terraform import to manage this service")
	} else if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	setServiceToResourceData(d, createdService)

	return nil
}

func resourceKongServiceRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*sling.Sling)

	id := d.Id()
	service := new(Service)

	response, e := s.New().Path("services/").Get(id).ReceiveSuccess(service)

	if e != nil {
		return fmt.Errorf("error while updating Service" + e.Error())
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	setServiceToResourceData(d, service)

	return nil
}

func resourceKongServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*sling.Sling)

	service := getServiceFromResourceData(d)

	updatedService := new(Service)

	response, e := s.New().BodyJSON(service).Patch("services/").Path(service.ID).ReceiveSuccess(updatedService)

	if e != nil {
		return fmt.Errorf("error while updating Service" + e.Error())
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	setServiceToResourceData(d, updatedService)

	return nil
}

func resourceKongServiceDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*sling.Sling)

	id := d.Id()

	response, e := s.New().Delete("services/").Path(id).ReceiveSuccess(nil)
	if e != nil {
		return fmt.Errorf("error while deleting Service" + e.Error())
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	return nil
}

func getServiceFromResourceData(d *schema.ResourceData) *Service {
	service := &Service{
		ID:             d.Id(),
		Name:           d.Get("name").(string),
		Retries:        d.Get("retries").(int),
		Protocol:       d.Get("protocol").(string),
		Host:           d.Get("host").(string),
		Port:           d.Get("port").(int),
		Path:           d.Get("path").(string),
		ConnectTimeout: d.Get("connect_timeout").(int),
		WriteTimeout:   d.Get("write_timeout").(int),
		ReadTimeout:    d.Get("read_timeout").(int),
		Tags:           helper.ConvertInterfaceArrToStrings(d.Get("tags").([]interface{})),
		ClientCertificate: Certificate{
			ID: d.Get("client_certificate").(string),
		},
		TlsVerify:      d.Get("tls_verify").(bool),
		TlsVerifyDepth: d.Get("tls_verify_depth").(int),
		CACertificates: helper.ConvertInterfaceArrToStrings(d.Get("ca_certificates").([]interface{})),
		Enabled:        d.Get("enabled").(bool),
	}

	return service
}

func setServiceToResourceData(d *schema.ResourceData, service *Service) {
	d.SetId(service.ID)
	_ = d.Set("name", service.Name)
	_ = d.Set("retries", service.Retries)
	_ = d.Set("protocol", service.Protocol)
	_ = d.Set("host", service.Host)
	_ = d.Set("port", service.Port)
	_ = d.Set("path", service.Path)
	_ = d.Set("connect_timeout", service.ConnectTimeout)
	_ = d.Set("write_timeout", service.WriteTimeout)
	_ = d.Set("read_timeout", service.ReadTimeout)
	_ = d.Set("tags", service.Tags)
	_ = d.Set("client_certificate", service.ClientCertificate)
	_ = d.Set("tls_verify", service.TlsVerify)
	_ = d.Set("tls_verify_depth", service.TlsVerifyDepth)
	_ = d.Set("ca_certificates", service.CACertificates)
	_ = d.Set("enabled", service.Enabled)
}
