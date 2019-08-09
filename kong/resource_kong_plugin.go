package kong

import (
	"encoding/json"
	"fmt"
	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"net/url"
	"strings"
)

// Plugin : Kong Service/API plugin request object structure
type Plugin struct {
	ID            string                 `json:"id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Configuration map[string]interface{} `json:"config,omitempty"`
	API           string                 `json:"-"`
	Service       string                 `json:"-"`
	Route         string                 `json:"-"`
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
				ConflictsWith: []string{ "config_json" },
			},

			"config_json": {
				Type: schema.TypeString,
				Optional: true,
				Default: nil,
				ConflictsWith: []string{ "config" },
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
				ConflictsWith: []string{"api", "route"},
			},

			"route": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       nil,
				ConflictsWith: []string{"api", "service"},
			},
		},
	}
}

func resourceKongPluginCreate(d *schema.ResourceData, meta interface{}) error {
	request := buildModifyRequest(d, meta)
	p := &Plugin{}

	if api, ok := d.GetOk("api"); ok {
		request = request.Path("apis/").Path(api.(string) + "/")
	} else if service, ok := d.GetOk("service"); ok {
		request = request.Path("services/").Path(service.(string) + "/")
	} else if route, ok := d.GetOk("route"); ok {
		request = request.Path("routes/").Path(route.(string) + "/")
	}

	response, err := request.Post("plugins/").ReceiveSuccess(p)
	if err != nil {
		return fmt.Errorf("error while creating plugin: " + err.Error())
	}

	if response.StatusCode == http.StatusConflict {
		return fmt.Errorf("409 Conflict - use terraform import to manage this plugin.")
	} else if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	return setPluginToResourceData(d, p)
}

func resourceKongPluginRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	p := &Plugin{}

	response, err := sling.New().Path("plugins/").Get(d.Id()).ReceiveSuccess(p)
	if err != nil {
		return fmt.Errorf("error while updating plugin: " + err.Error())
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	return setPluginToResourceData(d, p)
}

func resourceKongPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	request := buildModifyRequest(d, meta)

	p := &Plugin{}

	response, err := request.Path("plugins/").Patch(d.Id()).ReceiveSuccess(p)
	if err != nil {
		return fmt.Errorf("error while updating plugin: " + err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	return setPluginToResourceData(d, p)
}

func resourceKongPluginDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	response, error := sling.New().Path("plugins/").Delete(d.Id()).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("error while deleting plugin: " + error.Error())
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	return nil
}

func buildModifyRequest(d *schema.ResourceData, meta interface{}) *sling.Sling {
	request := meta.(*sling.Sling).New()

	plugin := &Plugin{
		ID:            d.Id(),
		Name:          d.Get("name").(string),
		API:           d.Get("api").(string),
		Service:       d.Get("service").(string),
		Route:         d.Get("route").(string),
		Consumer:      d.Get("consumer").(string),
	}

	if c, ok := d.GetOk("config_json"); ok {
		config := make(map[string]interface{})
		err := json.Unmarshal([]byte(c.(string)), &config)
		if err != nil {
			// ...
		}

		plugin.Configuration = config

		request = request.BodyJSON(plugin)
	} else {
		form := url.Values{
			"name": {plugin.Name},
		}

		if c, ok := d.GetOk("config"); ok {
			conf := c.(map[string]interface{})
			for k, v := range conf {
				form.Add("config."+k, v.(string))
			}
		}

		body := strings.NewReader(form.Encode())

		request = request.Body(body).Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return request
}

func setPluginToResourceData(d *schema.ResourceData, plugin *Plugin) error {
	d.SetId(plugin.ID)

	_ = d.Set("name", plugin.Name)

	// There are differences in the way service/route IDs are returned from Kong after creation and update between
	// version before and after 1.0.0. We are risking some drift here. This will be handled in later versions.
	if api, ok := d.GetOk("api"); ok {
		plugin.API = api.(string)
	} else if service, ok := d.GetOk("service"); ok {
		plugin.Service = service.(string)
	} else if route, ok := d.GetOk("route"); ok {
		plugin.Route = route.(string)
	}

	_ = d.Set("api", plugin.API)
	_ = d.Set("service", plugin.Service)
	_ = d.Set("route", plugin.Route)
	_ = d.Set("consumer", plugin.Consumer)

	return nil
}