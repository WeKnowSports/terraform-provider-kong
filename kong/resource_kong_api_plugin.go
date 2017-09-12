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
	ConsumerID    string                 `json:"consumer_id,omitempty"`
	Configuration map[string]interface{} `json:"config,omitempty"`
	API           string                 `json:"-"`
}

type KongAPIError struct {
	Message string `json:"message"`
}

func resourceKongPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongPluginCreate,
		Read:   resourceKongPluginRead,
		Update: resourceKongPluginUpdate,
		Delete: resourceKongPluginDelete,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"consumer_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The consumer_id of the plugin to use.",
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
				Required: true,
			},
		},
	}
}

func resourceKongPluginCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)
	createdPlugin := getPluginFromResourceData(d)

	apiError := new(KongAPIError)
	response, error := sling.New().BodyJSON(plugin).Path("apis/").Path(plugin.API+"/").Post("plugins/").Receive(createdPlugin, apiError)
	if error != nil {
		return fmt.Errorf("error while creating plugin: " + error.Error())
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status + " - " + apiError.Message)
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

	apiError := new(KongAPIError)
	response, error := sling.New().Path("apis/").Path(plugin.API + "/").Path("plugins/").Get(plugin.ID).ReceiveSuccess(plugin)
	if error != nil {
		return fmt.Errorf("error while updating plugin: " + error.Error())
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status + " - " + apiError.Message)
	}

	plugin.Configuration = configuration

	setPluginToResourceData(d, plugin)

	return nil
}

func resourceKongPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)

	updatedPlugin := getPluginFromResourceData(d)
	apiError := new(KongAPIError)
	// Disable saving state until we've confirmed there were no errors updating the plugin.
	d.Partial(true)
	response, error := sling.New().BodyJSON(plugin).Path("apis/").Path(plugin.API+"/").Path("plugins/").Patch(plugin.ID).Receive(updatedPlugin, apiError)
	if error != nil {
		return fmt.Errorf("error while updating plugin: " + error.Error() + plugin.ConsumerID + plugin.ID)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status + " - " + apiError.Message)
	}

	// We can re-enable saving state.
	d.Partial(false)
	updatedPlugin.Configuration = plugin.Configuration

	setPluginToResourceData(d, updatedPlugin)

	return nil
}

func resourceKongPluginDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)
	apiError := new(KongAPIError)
	response, error := sling.New().Path("apis/").Path(plugin.API+"/").Path("plugins/").Delete(plugin.ID).Receive(nil, apiError)
	if error != nil {
		return fmt.Errorf("error while deleting plugin: " + error.Error())
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status + " - " + apiError.Message)
	}

	return nil
}

func getPluginFromResourceData(d *schema.ResourceData) *Plugin {
	plugin := &Plugin{
		Name:          d.Get("name").(string),
		Configuration: d.Get("config").(map[string]interface{}),
		API:           d.Get("api").(string),
	}

	if id, ok := d.GetOk("id"); ok {
		plugin.ID = id.(string)
	}

	if consumer_id, ok := d.GetOk("consumer_id"); ok {
		plugin.ConsumerID = consumer_id.(string)
	}

	return plugin
}

func setPluginToResourceData(d *schema.ResourceData, plugin *Plugin) {
	d.SetId(plugin.ID)
	d.Set("name", plugin.Name)
	d.Set("config", plugin.Configuration)
	d.Set("consumer_id", plugin.ConsumerID)
	d.Set("api", plugin.API)
}
