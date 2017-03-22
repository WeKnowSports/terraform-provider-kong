package kong

import (
	"fmt"
	"net/http"
  "log"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
  "github.com/davecgh/go-spew/spew"
)

type KeyAuthConfig struct {
  KeyNames        string  `json:"key_names,omitempty"`
  HideCredentials bool    `json:"hide_credentials,omitempty"`
  Anonymous       bool    `json:"anonymous,omitempty"`
}

type KeyAuthPlugin struct {
	ID            string                 `json:"id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Configuration KeyAuthConfig `json:"config,omitempty"`
	API           string                 `json:"-"`
}

func resourceKongKeyAuthPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyAuthPluginCreate,
		Read:   resourceKeyAuthPluginRead,
		Update: resourceKeyAuthPluginUpdate,
		Delete: resourceKeyAuthPluginDelete,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"config": &schema.Schema{
        Type: schema.TypeMap,
        Required: true,
        Default: nil,
        Elem: &schema.Resource{
          Schema: map[string]*schema.Schema{
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
          },
        },
			},

			"api": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceKeyAuthPluginCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getKeyAuthPluginFromResourceData(d)

	createdPlugin := getKeyAuthPluginFromResourceData(d)

	response, error := sling.New().BodyJSON(plugin).Path("apis/").Path(plugin.API + "/").Post("plugins/").ReceiveSuccess(createdPlugin)
  pluginStr := spew.Sdump(createdPlugin)
  str := spew.Sdump(response)
  log.Print(pluginStr)
  log.Print(str)
	if error != nil {
		return fmt.Errorf("Error while creating plugin.")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setKeyAuthPluginFromResourceData(d, createdPlugin)

	return nil
}

func resourceKeyAuthPluginRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getKeyAuthPluginFromResourceData(d)

	response, error := sling.New().Path("apis/").Path(plugin.API + "/").Path("plugins/").Get(plugin.ID).ReceiveSuccess(plugin)
	if error != nil {
		return fmt.Errorf("Error while updating plugin.")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setKeyAuthPluginFromResourceData(d, plugin)

	return nil
}

func resourceKeyAuthPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	plugin := getKeyAuthPluginFromResourceData(d)

	updatedPlugin := getKeyAuthPluginFromResourceData(d)

	response, error := sling.New().BodyJSON(plugin).Path("apis/").Path(plugin.API + "/").Path("plugins/").Patch(plugin.ID).ReceiveSuccess(updatedPlugin)
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

	response, error := sling.New().Path("apis/").Path(plugin.API + "/").Path("plugins/").Delete(plugin.ID).ReceiveSuccess(nil)
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
		Name:          "key-auth",
		Configuration: d.Get("config").(KeyAuthConfig),
		API:           d.Get("api").(string),
	}

  log.Print("GET KEY AUTH PLUGIN")
  pluginStr := spew.Sdump(plugin)
  log.Print(pluginStr)


  // if config, ok := d.Get("config").(KeyAuthConfig); ok {
  //   if config.HideCredentials != "" {
  //     plugin.Configuration.HideCredentials = strconv.ParseBool(plugin.Configuration.HideCredentials)
  //   }
  // }

	if id, ok := d.GetOk("id"); ok {
		plugin.ID = id.(string)
	}

	return plugin
}

func setKeyAuthPluginFromResourceData(d *schema.ResourceData, plugin *KeyAuthPlugin) {
	d.SetId(plugin.ID)
	d.Set("name", "key-auth")
	d.Set("config", plugin.Configuration)
	d.Set("api", plugin.API)
}
