package kong

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Target struct {
	ID       string `json:"id,omitempty"`
	Upstream string `json:"-"`
	Target   string `json:"target,omitempty"`
}

func resourceKongTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongTargetCreate,
		Read:   resourceKongTargetRead,
		Delete: resourceKongTargetDelete,

		Schema: map[string]*schema.Schema{
			"upstream": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier or the name of the upstream to which to add the target.",
				ForceNew:    true,
			},
			"target": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The target address (ip or hostname) and port. If omitted the port defaults to 8000. If the hostname resolves to an SRV record, the port value will overridden by the value from the dns record.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSuffix(new, ":8000") == strings.TrimSuffix(old, ":8000")
				},
				ForceNew: true,
			},
		},
	}
}

func resourceKongTargetCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	target := getTargetFromResourceData(d)

	createdTarget := getTargetFromResourceData(d)

	response, error := sling.New().Path("upstreams/").Path(target.Upstream + "/").BodyJSON(target).Post("targets/").ReceiveSuccess(createdTarget)
	if error != nil {
		return fmt.Errorf("error while creating target")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setTargetToResourceData(d, createdTarget)

	return nil
}

func resourceKongTargetRead(d *schema.ResourceData, meta interface{}) error {
	// Targets can't be read, so we ignore the read operation.
	return nil
}

func resourceKongTargetDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	target := getTargetFromResourceData(d)

	response, error := sling.New().Path("upstreams/").Path(target.Upstream + "/").Path("targets/").Delete(target.ID).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("error while deleting target")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getTargetFromResourceData(d *schema.ResourceData) *Target {
	target := &Target{
		ID:       d.Id(),
		Target:   d.Get("target").(string),
		Upstream: d.Get("upstream").(string),
	}

	return target
}

func setTargetToResourceData(d *schema.ResourceData, target *Target) {
	d.SetId(target.ID)
	d.Set("target", target.Target)
	d.Set("upstream", target.Upstream)
}
