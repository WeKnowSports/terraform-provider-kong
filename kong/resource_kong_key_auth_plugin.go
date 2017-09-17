package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

type KeyAuthPlugin struct {
	ID               string                 `json:"id,omitempty"`
	Name             string                 `json:"name,omitempty"`
	KeyNames      	 string                 `json:"config.key_names,omitempty"`
	HideCredentials  bool                   `json:"config.hide_credentials,omitempty"`
	Anonymous        string                  `json:"config.anonymous,omitempty"`
	API              string                 `json:"api_id,omitempty"`
}

func resourceKongKeyAuthPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyAuthPluginCreate,
		Read:   resourceKeyAuthPluginRead,
		Update: resourceKeyAuthPluginUpdate,
		Delete: resourceKeyAuthPluginDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

      "key_names": &schema.Schema{
        Type:   schema.TypeString,
        Required: true,
        Default: nil,
        Description: "The name of the API key header to use.",
      },

      "hide_credentials": &schema.Schema{
        Type: schema.TypeBool,
        Optional: true,
        Default: nil,
        Description: "Whether credentials should be hidden.",
      },

      "anonymous": &schema.Schema{
        Type: schema.TypeString,
        Optional: true,
        Default: nil,
        Description: "String (consumer UUID) to use as an anonymous 'consumer', if authentication fails.",
      },

			"api": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
		},
	}
}

func resourceKeyAuthPluginCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getKeyAuthPluginFromResourceData(d)

	createdPlugin := getKeyAuthPluginFromResourceData(d)

	request := sling.New().BodyJSON(plugin)
	if plugin.API != "" {
		request = request.Path("apis/").Path(plugin.API + "/")
	}
	response, error := request.Post("plugins/").ReceiveSuccess(createdPlugin)
	if error != nil {
		return fmt.Errorf("Error while creating plugin.")
	}

	if response.StatusCode == http.StatusConflict {
		return fmt.Errorf("409 Conflict - use terraform import to manage this plugin.")
	} else if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setKeyAuthPluginFromResourceData(d, createdPlugin)

	return nil
}

func resourceKeyAuthPluginRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getKeyAuthPluginFromResourceData(d)

	response, error := sling.New().Path("plugins/").Get(plugin.ID).ReceiveSuccess(plugin)
	if error != nil {
		return fmt.Errorf("Error while updating plugin.")
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setKeyAuthPluginFromResourceData(d, plugin)

	return nil
}

func resourceKeyAuthPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getKeyAuthPluginFromResourceData(d)

	updatedPlugin := getKeyAuthPluginFromResourceData(d)

	response, error := sling.New().BodyJSON(plugin).Path("plugins/").Patch(plugin.ID).ReceiveSuccess(updatedPlugin)
	if error != nil {
		return fmt.Errorf("Error while updating plugin.")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setKeyAuthPluginFromResourceData(d, updatedPlugin)

	return nil
}

func resourceKeyAuthPluginDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getKeyAuthPluginFromResourceData(d)

	response, error := sling.New().Path("plugins/").Delete(plugin.ID).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("Error while deleting plugin.")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getKeyAuthPluginFromResourceData(d *schema.ResourceData) *KeyAuthPlugin {
	plugin := &KeyAuthPlugin{
		Name:            "key-auth",
		KeyNames:        d.Get("key_names").(string),
		HideCredentials: d.Get("hide_credentials").(bool),
		Anonymous:       d.Get("anonymous").(string),
		API:             d.Get("api").(string),
	}

	if id, ok := d.GetOk("id"); ok {
		plugin.ID = id.(string)
	}

	return plugin
}

func setKeyAuthPluginFromResourceData(d *schema.ResourceData, plugin *KeyAuthPlugin) {
	d.SetId(plugin.ID)
	d.Set("name", plugin.Name)
	d.Set("key_names", plugin.KeyNames)
	d.Set("hide_credentials", plugin.HideCredentials)
	d.Set("anonymous", plugin.Anonymous)
	d.Set("api", plugin.API)
}
