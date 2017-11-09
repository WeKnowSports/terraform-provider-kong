package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

type UpdateAclRequest struct {
	Group    string `json:"group"`
}

type ConsumerACL struct {
	ID       string `json:"id,omitempty"`
	Consumer string `json:"consumer_id"`
	Group    string `json:"group"`
}

func resourceKongConsumerACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongConsumerACLCreate,
		Read:   resourceKongConsumerACLRead,
		Update: resourceKongConsumerACLUpdate,
		Delete: resourceKongConsumerACLDelete,

		Importer: &schema.ResourceImporter{
			State: ImportConsumerACL,
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the consumer-acl association.",
			},

			"consumer": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew: true,
				Description: "The id of the consumer to associate this group with.",
			},

			"group": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the group to place the specified consumer in.",
			},
		},
	}
}

func resourceKongConsumerACLCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	createRequest := &UpdateAclRequest{
		Group: d.Get("group").(string),
	}
	consumer := d.Get("consumer").(string)
	updated := &ConsumerACL{}
	failure := &UpdateAclRequest{}
	response, error := sling.New().BodyJSON(createRequest).Path("consumers/").Path(consumer + "/").Post("acls/").Receive(updated, failure)
	if error != nil {
		return fmt.Errorf("error while creating ACL" + error.Error())
	}

	if response.StatusCode == http.StatusBadRequest {
		return fmt.Errorf(response.Status + " - " + failure.Group + "; use `terraform import <consumer>/<acl>` to manage this resource with terraform.")
	} else if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	d.Set("consumer", updated.Consumer)
	d.Set("group", updated.Group)
	d.SetId(updated.ID)

	return nil
}

func resourceKongConsumerACLRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	consumer := d.Get("consumer").(string)
	updated := &ConsumerACL{}
	response, error := sling.New().Path("consumers/").Path(consumer + "/").Path("acls/").Get(d.Id()).ReceiveSuccess(updated)
	if error != nil {
		return fmt.Errorf("error while reading ACL" + error.Error())
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	d.Set("consumer", updated.Consumer)
	d.Set("group", updated.Group)
	// Update the id field initally imported via group name.
	d.SetId(updated.ID)

	return nil
}

func resourceKongConsumerACLUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	updateRequest := &UpdateAclRequest{
		Group: d.Get("group").(string),
	}
	consumer := d.Get("consumer").(string)
	updated := &ConsumerACL{}
	failure := &UpdateAclRequest{}
	response, error := sling.New().BodyJSON(updateRequest).Path("consumers/").Path(consumer + "/").Path("acls/").Patch(d.Id()).Receive(updated, failure)
	if error != nil {
		return fmt.Errorf("error while creating ACL" + error.Error())
	}

	if response.StatusCode == http.StatusBadRequest {
		return fmt.Errorf(response.Status + " - " + failure.Group + "; use `terraform import <consumer>/<acl>` to manage this resource with terraform.")
	} else if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	d.Set("consumer", updated.Consumer)
	d.Set("group", updated.Group)
	d.SetId(updated.ID)

	return nil
}


func resourceKongConsumerACLDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	consumer := d.Get("consumer").(string)
	response, error := sling.New().Path("consumers/").Path(consumer + "/").Path("acls/").Delete(d.Id()).ReceiveSuccess(nil)

	if error != nil {
		return fmt.Errorf("error while deleting ACL" + error.Error())
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	return nil
}