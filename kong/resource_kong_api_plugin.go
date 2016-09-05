package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

type Plugin struct {
	ID            string                 `json:"id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Configuration map[string]interface{} `json:"config,omitempty"`
	API           string                 `json:"-"`
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

	response, error := sling.New().BodyJSON(plugin).Path("apis/").Path(plugin.API + "/").Post("plugins/").ReceiveSuccess(createdPlugin)
	if error != nil {
		return fmt.Errorf("Error while creating plugin.")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setPluginToResourceData(d, createdPlugin)

	return nil
}

func resourceKongPluginRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)

	response, error := sling.New().Path("apis/").Path(plugin.API + "/").Path("plugins/").Get(plugin.ID).ReceiveSuccess(plugin)
	if error != nil {
		return fmt.Errorf("Error while updating plugin.")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setPluginToResourceData(d, plugin)

	return nil
}

func resourceKongPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)

	updatedPlugin := getPluginFromResourceData(d)

	response, error := sling.New().BodyJSON(plugin).Path("apis/").Path(plugin.API + "/").Path("plugins/").Patch(plugin.ID).ReceiveSuccess(updatedPlugin)
	if error != nil {
		return fmt.Errorf("Error while updating plugin.")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setPluginToResourceData(d, updatedPlugin)

	return nil
}

func resourceKongPluginDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)

	response, error := sling.New().Path("apis/").Path(plugin.API + "/").Path("plugins/").Delete(plugin.ID).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("Error while deleting plugin.")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
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

	return plugin
}

func setPluginToResourceData(d *schema.ResourceData, plugin *Plugin) {
	d.SetId(plugin.ID)
	d.Set("name", plugin.Name)
	d.Set("config", plugin.Configuration)
	d.Set("api", plugin.API)
}
