package kong

import (
	"fmt"
	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strconv"
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
	Url            string `json:"url,omitempty"`
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return (new == "" && d.Get("url").(string) != "") || (old == "http" && (new == ""))
				},
			},

			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The host of the upstream server.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "" && d.Get("url").(string) != ""
				},
			},

			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The upstream server port. Defaults to 80.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "0" && d.Get("url").(string) != "" || (old == "80" && new == "0")
				},
			},

			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to be used in requests to the upstream server. Empty by default.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "" && d.Get("url").(string) != ""
				},
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

			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Shorthand attribute to set protocol, host, port and path at once. This attribute is write-only (the Admin API never \"returns\" the url).",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					service := getServiceFromResourceData(d)

					oldUrl := service.Protocol + "://" + service.Host + ":" + strconv.FormatInt(int64(service.Port), 10) + service.Path
					oldUrlNoPort := service.Protocol + "://" + service.Host + service.Path

					return new == oldUrl || new == oldUrlNoPort
				},
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
		Protocol:       d.Get("protocol").(string),
		Host:           d.Get("host").(string),
		Port:           d.Get("port").(int),
		Path:           d.Get("path").(string),
		Retries:        d.Get("retries").(int),
		ConnectTimeout: d.Get("connect_timeout").(int),
		WriteTimeout:   d.Get("write_timeout").(int),
		ReadTimeout:    d.Get("read_timeout").(int),
		Url:            d.Get("url").(string),
	}

	return service
}

func setServiceToResourceData(d *schema.ResourceData, service *Service) {
	d.SetId(service.ID)
	_ = d.Set("name", service.Name)
	_ = d.Set("protocol", service.Protocol)
	_ = d.Set("host", service.Host)
	_ = d.Set("port", service.Port)
	_ = d.Set("path", service.Path)
	_ = d.Set("retries", service.Retries)
	_ = d.Set("connect_timeout", service.ConnectTimeout)
	_ = d.Set("write_timeout", service.WriteTimeout)
	_ = d.Set("read_timeout", service.ReadTimeout)
	_ = d.Set("url", service.Url)
}
