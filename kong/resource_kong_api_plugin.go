package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

// Plugin : Kong API plugin request object structure
type Plugin struct {
	ID            string                 `json:"id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Configuration map[string]interface{} `json:"config,omitempty"`
	API           string                 `json:"-"`
	Consumer      string                 `json:"consumer_id,omitempty"`
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
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"consumer": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The id of the consumer to scope this plugin to.",
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The name of the plugin to use.",
			},

			"config": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
				Default:  nil,
			},

			"api": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
		},
	}
}

func resourceKongPluginCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)

	createdPlugin := getPluginFromResourceData(d)

	request := sling.New().BodyJSON(plugin)
	if plugin.API != "" {
		request = request.Path("apis/").Path(plugin.API + "/")
	}
	response, error := request.Post("plugins/").ReceiveSuccess(createdPlugin)
	if error != nil {
		return fmt.Errorf("error while creating plugin: " + error.Error())
	}

	if response.StatusCode == http.StatusConflict {
		return fmt.Errorf("409 Conflict - use terraform import to manage this plugin.")
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

	response, error := sling.New().BodyJSON(plugin).Path("plugins/").Patch(plugin.ID).ReceiveSuccess(updatedPlugin)
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
		Consumer:      d.Get("consumer").(string),
	}

	if id, ok := d.GetOk("id"); ok {
		plugin.ID = id.(string)
	}

	return plugin
}

func setPluginToResourceData(d *schema.ResourceData, plugin *Plugin) {
	d.SetId(plugin.ID)
	d.Set("name", plugin.Name)
	d.Set("config", plugin.Configuration)
	d.Set("api", plugin.API)
	d.Set("consumer", plugin.Consumer)
}
