package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

type Upstream struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func resourceKongUpstream() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongUpstreamCreate,
		Read:   resourceKongUpstreamRead,
		Update: resourceKongUpstreamUpdate,
		Delete: resourceKongUpstreamDelete,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is a hostname, which must be equal to the host of a Service.",
			},
		},
	}
}

func resourceKongUpstreamCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	upstream := getUpstreamFromResourceData(d)

	createdUpstream := getUpstreamFromResourceData(d)

	response, error := sling.New().BodyJSON(upstream).Post("upstreams/").ReceiveSuccess(createdUpstream)
	if error != nil {
		return fmt.Errorf("Error while creating upstream.")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setUpstreamToResourceData(d, createdUpstream)

	return nil
}

func resourceKongUpstreamRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	upstream := getUpstreamFromResourceData(d)

	response, error := sling.New().Path("upstreams/").Get(upstream.ID).ReceiveSuccess(upstream)
	if error != nil {
		return fmt.Errorf("Error while updating upstream.")
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setUpstreamToResourceData(d, upstream)

	return nil
}

func resourceKongUpstreamUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	upstream := getUpstreamFromResourceData(d)

	updatedUpstream := getUpstreamFromResourceData(d)

	response, error := sling.New().BodyJSON(upstream).Path("upstreams/").Patch(upstream.ID).ReceiveSuccess(updatedUpstream)
	if error != nil {
		return fmt.Errorf("Error while updating upstream")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setUpstreamToResourceData(d, updatedUpstream)

	return nil
}

func resourceKongUpstreamDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	upstream := getUpstreamFromResourceData(d)

	response, error := sling.New().Path("upstreams/").Delete(upstream.ID).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("Error while deleting upstream.")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getUpstreamFromResourceData(d *schema.ResourceData) *Upstream {
	upstream := &Upstream{
		Name: d.Get("name").(string),
	}

	if id, ok := d.GetOk("id"); ok {
		upstream.ID = id.(string)
	}

	return upstream
}

func setUpstreamToResourceData(d *schema.ResourceData, upstream *Upstream) {
	d.SetId(upstream.ID)
	d.Set("name", upstream.Name)
}
