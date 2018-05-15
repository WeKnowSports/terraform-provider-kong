package kong

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

// Plugin : Kong API plugin request object structure
type Plugin struct {
	ID            string                 `json:"id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Configuration map[string]interface{} `json:"config,omitempty"`
	API           string                 `json:"api_id,omitempty"`
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
				Type:          schema.TypeMap,
				Optional:      true,
				Elem:          schema.TypeString,
				Default:       nil,
				ConflictsWith: []string{"config_json"},
			},

			"config_json": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Default:       nil,
				ConflictsWith: []string{"config"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					var oldConfig map[string]interface{}
					json.Unmarshal([]byte(old), &oldConfig)

					var newConfig map[string]interface{}
					json.Unmarshal([]byte(new), &newConfig)

					return reflect.DeepEqual(oldConfig, newConfig)
				},
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if !json.Valid([]byte(v.(string))) {
						errors = append(errors, fmt.Errorf("Invalid JSON: %v", v))
					}

					return
				},
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

	plugin, err := getPluginFromResourceData(d)
	if err != nil {
		return err
	}

	createdPlugin, err := getPluginFromResourceData(d)
	if err != nil {
		return err
	}

	request := sling.New().BodyJSON(plugin)
	if plugin.API != "" {
		request = request.Path("apis/").Path(plugin.API + "/")
	}
	errorResponse := make(map[string]interface{})
	response, error := request.Post("plugins/").Receive(createdPlugin, &errorResponse)
	if error != nil {
		return fmt.Errorf("error while creating plugin: " + error.Error())
	}

	if response.StatusCode == http.StatusConflict {
		return fmt.Errorf("409 Conflict - use terraform import to manage this plugin.")
	} else if response.StatusCode != http.StatusCreated {
		return ErrorFromResponse(response, errorResponse)
	}

	setPluginToResourceData(d, createdPlugin)

	return nil
}

func resourceKongPluginRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := &Plugin{}
	errorResponse := make(map[string]interface{})
	response, error := sling.New().Path("plugins/").Get(d.Id()).Receive(plugin, &errorResponse)
	if error != nil {
		return fmt.Errorf("error while updating plugin: " + error.Error())
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return ErrorFromResponse(response, errorResponse)
	}

	setPluginToResourceData(d, plugin)

	return nil
}

func resourceKongPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin, err := getPluginFromResourceData(d)
	if err != nil {
		return err
	}

	updatedPlugin, err := getPluginFromResourceData(d)
	if err != nil {
		return err
	}

	errorResponse := make(map[string]interface{})
	response, error := sling.New().BodyJSON(plugin).Path("plugins/").Patch(d.Id()).Receive(updatedPlugin, &errorResponse)
	if error != nil {
		return fmt.Errorf("error while updating plugin: " + error.Error())
	}

	if response.StatusCode != http.StatusOK {
		return ErrorFromResponse(response, errorResponse)
	}

	setPluginToResourceData(d, updatedPlugin)

	return nil
}

func resourceKongPluginDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)
	errorResponse := make(map[string]interface{})
	response, error := sling.New().Path("plugins/").Delete(d.Id()).Receive(nil, &errorResponse)
	if error != nil {
		return fmt.Errorf("error while deleting plugin: " + error.Error())
	}

	if response.StatusCode != http.StatusNoContent {
		return ErrorFromResponse(response, errorResponse)
	}

	return nil
}

func setPluginConfig(d *schema.ResourceData, config map[string]interface{}) error {
	// If a config map was specified in the schema store the raw map.
	if _, ok := d.GetOk("config"); ok {
		err := d.Set("config_json", nil)
		if err != nil {
			return fmt.Errorf("%v: Unable to set config_json to nil.", err)
		}

		return d.Set("config", config)
	}

	// Otherwise marshal the configuration map into a json blob.
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = d.Set("config", nil)
	if err != nil {
		return fmt.Errorf("%v: Unable to set config to nil.", err)
	}

	err = d.Set("config_json", string(jsonBytes))
	if err != nil {
		return fmt.Errorf("%v: Unable to set config_json to %+v.", err, string(jsonBytes))
	}

	return nil
}

func getPluginConfig(d *schema.ResourceData) (map[string]interface{}, error) {
	if config, ok := d.GetOk("config"); ok {
		return config.(map[string]interface{}), nil
	}

	if configBlob, ok := d.GetOk("config_json"); ok {
		var config map[string]interface{}
		err := json.Unmarshal([]byte(configBlob.(string)), &config)
		if err != nil {
			return nil, err
		}

		return config, nil
	}

	// No plugin configuration was specified meaning terraform-provider-kong will not manage the plugin config field.
	return nil, nil
}

func getPluginFromResourceData(d *schema.ResourceData) (*Plugin, error) {
	plugin := &Plugin{
		Name:     d.Get("name").(string),
		API:      d.Get("api").(string),
		Consumer: d.Get("consumer").(string),
	}

	config, err := getPluginConfig(d)
	if err != nil {
		return nil, err
	}

	if config != nil {
		plugin.Configuration = config
	}

	return plugin, nil
}

func setPluginToResourceData(d *schema.ResourceData, plugin *Plugin) {
	d.SetId(plugin.ID)
	d.Set("name", plugin.Name)
	d.Set("api", plugin.API)
	d.Set("consumer", plugin.Consumer)

	setPluginConfig(d, plugin.Configuration)
}
