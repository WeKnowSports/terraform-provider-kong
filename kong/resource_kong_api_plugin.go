package kong

import (
	"fmt"
	"net/http"
	"sort"
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
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The name of the plugin to use.",
			},

			"config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type: schema.TypeString,
							Required: true,
						},
						"value": {
							Type: schema.TypeString,
							Required: true,
						},
					},
				},
				Default:  nil,
			},

			"api": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceKongPluginCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)

	createdPlugin := getPluginFromResourceData(d)
	response, err := s.New().BodyJSON(plugin).Path("apis/").Path(plugin.API + "/").Post("plugins/").ReceiveSuccess(createdPlugin)
	if err != nil {
		return fmt.Errorf("Error while creating plugin.")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setPluginToResourceData(d, createdPlugin)

	return nil
}

func resourceKongPluginRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)

	response, err := s.New().Path("apis/").Path(plugin.API + "/").Path("plugins/").Get(plugin.ID).ReceiveSuccess(plugin)
	if err != nil {
		return fmt.Errorf("Error while updating plugin.")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setPluginToResourceData(d, plugin)

	return nil
}

func resourceKongPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)

	updatedPlugin := getPluginFromResourceData(d)

	response, err := s.New().BodyJSON(plugin).Path("apis/").Path(plugin.API + "/").Path("plugins/").Patch(plugin.ID).ReceiveSuccess(updatedPlugin)
	if err != nil {
		return fmt.Errorf("Error while updating plugin.")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setPluginToResourceData(d, updatedPlugin)

	return nil
}

func resourceKongPluginDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*sling.Sling)

	plugin := getPluginFromResourceData(d)

	response, err := s.New().Path("apis/").Path(plugin.API + "/").Path("plugins/").Delete(plugin.ID).ReceiveSuccess(nil)
	if err != nil {
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
		API:           d.Get("api").(string),
	}

	if id, ok := d.GetOk("id"); ok {
		plugin.ID = id.(string)

	}

	configs := d.Get("config").(*schema.Set).List()
	configuration := make(map[string]interface{}, len(configs))

	for _, config := range configs {
		name := config.(map[string]interface{})["name"].(string)
		value := config.(map[string]interface{})["value"].(string)
		configuration[name] = value
	}
	plugin.Configuration = configuration

	return plugin
}

func setPluginToResourceData(d *schema.ResourceData, plugin *Plugin) {
	config := plugin.Configuration
	keys := make([]string, len(config))

	for k := range config {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	configs := make([]interface{}, len(keys))
	for k := range keys {
		configs = append(configs, map[string]interface{}{
			"name": keys[k],
			"value": config[keys[k]],
		})
	}

	d.SetId(plugin.ID)
	d.Set("name", plugin.Name)
	d.Set("config", configs)
	d.Set("api", plugin.API)
}
