package kong

import (
	"fmt"
	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
)

// Service : Kong Service request object structure
type Service struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Protocol       string `json:"protocol,omitempty"`
	Host           string `json:"host,omitempty"`
	Port           int    `json:"port,omitempty"`
	Path           string `json:"path,omitempty"`
	Retries        int    `json:"retries,omitempty"`
	ConnectTimeout int    `json:"connect_timeout,omitempty"`
	WriteTimeout   int    `json:"write_timeout,omitempty"`
	ReadTimeout    int    `json:"read_timeout,omitempty"`
	Url            string `json:"-"`
}

func resourceKongService() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongServiceCreate,
		Read:   resourceKongServiceRead,
		Update: resourceKongServiceUpdate,
		Delete: resourceKongServiceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The Service name.",
			},

			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "http",
				Description: "The protocol used to communicate with the upstream. It can be one of http (default) or https.",
			},

			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The host of the upstream server.",
			},

			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     80,
				Description: "The upstream server port. Defaults to 80.",
			},

			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The path to be used in requests to the upstream server. Empty by default.",
			},

			"retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "A comma-separated list of HTTP methods that point to your Service. For example: GET,POST. At least one of hosts, uris, or methods should be specified.",
			},

			"connect_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60000,
				Description: "The timeout in milliseconds for establishing a connection to the upstream server. Defaults to 60000.",
			},

			"write_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60000,
				Description: "The timeout in milliseconds between two successive write operations for transmitting a request to the upstream server. Defaults to 60000.",
			},

			"read_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60000,
				Description: "The timeout in milliseconds between two successive read operations for transmitting a request to the upstream server. Defaults to 60000.",
			},

			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "Shorthand attribute to set protocol, host, port and path at once. This attribute is write-only (the Admin API never \"returns\" the url).",
			},
		},
	}
}

func resourceKongServiceCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	service := getServiceFromResourceData(d)

	createdService := new(Service)
	response, error := sling.New().BodyJSON(service).Post("services/").ReceiveSuccess(createdService)

	if error != nil {
		return fmt.Errorf("error while creating Service: " + error.Error())
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
	sling := meta.(*sling.Sling)

	id := d.Get("id").(string)
	service := new(Service)

	response, error := sling.New().Path("services/").Get(id).ReceiveSuccess(service)

	if error != nil {
		return fmt.Errorf("error while updating Service" + error.Error())
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
	sling := meta.(*sling.Sling)

	service := getServiceFromResourceData(d)

	updatedService := new(Service)

	response, error := sling.New().BodyJSON(service).Patch("services/").Path(service.ID).ReceiveSuccess(updatedService)

	if error != nil {
		return fmt.Errorf("error while updating Service" + error.Error())
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	setServiceToResourceData(d, updatedService)

	return nil
}

func resourceKongServiceDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	id := d.Get("id").(string)

	response, error := sling.New().Delete("services/").Path(id).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("error while deleting Service" + error.Error())
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	return nil
}

func getServiceFromResourceData(d *schema.ResourceData) *Service {
	service := &Service{
		Name:           d.Get("name").(string),
		Protocol:       d.Get("protocol").(string),
		Host:           d.Get("host").(string),
		Port:           d.Get("port").(int),
		Path:           d.Get("path").(string),
		Retries:        d.Get("retries").(int),
		ConnectTimeout: d.Get("connect_timeout").(int),
		WriteTimeout:   d.Get("write_timeout").(int),
		ReadTimeout:    d.Get("read_timeout").(int),
	}

	if id, ok := d.GetOk("id"); ok {
		service.ID = id.(string)
	}

	return service
}

func setServiceToResourceData(d *schema.ResourceData, service *Service) {
	d.SetId(service.ID)
	d.Set("name", service.Name)
	d.Set("protocol", service.Protocol)
	d.Set("host", service.Host)
	d.Set("port", service.Port)
	d.Set("path", service.Path)
	d.Set("retries", service.Retries)
	d.Set("connect_timeout", service.ConnectTimeout)
	d.Set("write_timeout", service.WriteTimeout)
	d.Set("read_timeout", service.ReadTimeout)
}
