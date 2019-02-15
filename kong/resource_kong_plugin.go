package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

// Plugin : Kong Service/API plugin request object structure
type Plugin struct {
	ID            string                 `json:"id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Configuration map[string]interface{} `json:"config,omitempty"`
	API           string                 `json:"api_id,omitempty"`
	Service       string                 `json:"service_id,omitempty"`
	Route         string                 `json:"route_id,omitempty"`
	Consumer      string                 `json:"consumer_id,omitempty"`
	CreatedAt     int                    `json:"created_at,omitempty"`
}

func resourceKongPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongPluginCreate,
		Read:   resourceKongPluginRead,
		Update: resourceKongPluginUpdate,
		Delete: resourceKongPluginDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"consumer": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The id of the consumer to scope this plugin to.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The name of the plugin to use.",
			},

			"config": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
				Default:  nil,
			},

			"api": {
				Type:          schema.TypeString,
				Optional:      true,
				Deprecated:    "Use service or route instead.",
				Default:       nil,
				ConflictsWith: []string{"service", "route"},
			},

			"service": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       nil,
				ConflictsWith: []string{"api"},
			},

			"route": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       nil,
				ConflictsWith: []string{"api"},
			},

			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceKongPluginCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)

	createdPlugin := getPluginFromResourceData(d)

	response, error := sling.New().BodyJSON(plugin).Post("plugins/").ReceiveSuccess(createdPlugin)
	if error != nil {
		return fmt.Errorf("error while creating plugin: " + error.Error())
	}

	if response.StatusCode == http.StatusConflict {
		return fmt.Errorf("409 Conflict - use terraform import to manage this plugin")
	} else if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	createdPlugin.Configuration = plugin.Configuration

	setPluginToResourceData(d, createdPlugin)

	return nil
}

func resourceKongPluginRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)

	configuration := make(map[string]interface{})
	for key, value := range plugin.Configuration {
		configuration[key] = value
	}

	response, error := sling.New().Path("plugins/").Get(plugin.ID).ReceiveSuccess(plugin)
	if error != nil {
		return fmt.Errorf("error while updating plugin: " + error.Error())
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	plugin.Configuration = configuration

	setPluginToResourceData(d, plugin)

	return nil
}

func resourceKongPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)

	updatedPlugin := getPluginFromResourceData(d)

	response, error := sling.New().BodyJSON(plugin).Put("plugins/").ReceiveSuccess(updatedPlugin)
	if error != nil {
		return fmt.Errorf("error while updating plugin: " + error.Error())
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	updatedPlugin.Configuration = plugin.Configuration

	setPluginToResourceData(d, updatedPlugin)

	return nil
}

func resourceKongPluginDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)

	response, error := sling.New().Path("plugins/").Delete(plugin.ID).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("error while deleting plugin: " + error.Error())
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	return nil
}

func getPluginFromResourceData(d *schema.ResourceData) *Plugin {
	plugin := &Plugin{
		Name:          d.Get("name").(string),
		Configuration: d.Get("config").(map[string]interface{}),
		API:           d.Get("api").(string),
		Service:       d.Get("service").(string),
		Route:         d.Get("route").(string),
		Consumer:      d.Get("consumer").(string),
		Service:       d.Get("service").(string),
		Route:         d.Get("route").(string),
	}

	if id, ok := d.GetOk("id"); ok {
		plugin.ID = id.(string)
	}

	if createdAt, ok := d.GetOk("created_at"); ok {
		plugin.CreatedAt = createdAt.(int)
	}

	return plugin
}

func setPluginToResourceData(d *schema.ResourceData, plugin *Plugin) {
	d.SetId(plugin.ID)
	d.Set("name", plugin.Name)
	d.Set("config", plugin.Configuration)
	d.Set("api", plugin.API)
	d.Set("service", plugin.Service)
	d.Set("route", plugin.Route)
	d.Set("consumer", plugin.Consumer)
	d.Set("serivce", plugin.Service)
	d.Set("route", plugin.Route)
	d.Set("created_at", plugin.CreatedAt)
}
