package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

type CreateAclRequest struct {
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
				ForceNew: true,
				Description: "The name of the group to place the specified consumer in.",
			},
		},
	}
}

func resourceKongConsumerACLCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	createRequest := &CreateAclRequest{
		Group: d.Get("group").(string),
	}
	consumer := d.Get("consumer").(string)
	updated := &ConsumerACL{}
	response, error := sling.New().BodyJSON(createRequest).Path("consumers/").Path(consumer + "/").Post("acls/").ReceiveSuccess(updated)
	if error != nil {
		return fmt.Errorf("error while creating ACL" + error.Error())
	}

	if response.StatusCode != http.StatusCreated {
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
	response, error := sling.New().Path("consumers/").Path(consumer).Path("acls/").Get(d.Id()).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("error while reading ACL" + error.Error())
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	return nil
}

func resourceKongConsumerACLDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	id := d.Get("id").(string)
	consumer := d.Get("consumer").(string)
	response, error := sling.New().Delete("consumers/").Path(consumer).Path("acls/").Path(id).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("error while deleting ACL" + error.Error())
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code received: " + response.Status)
	}

	return nil
}
