package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

type ConsumerACLGroup struct {
	ID       string `json:"id,omitempty"`
	Group    string `json:"group,omitempty"`
	Consumer string `json:"-"`
}

func resourceKongConsumerACLGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongConsumerACLGroupCreate,
		Read:   resourceKongConsumerACLGroupRead,
		Update: resourceKongConsumerACLGroupUpdate,
		Delete: resourceKongConsumerACLGroupDelete,

		Schema: map[string]*schema.Schema{
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The arbitrary group name to associate to the consumer.",
			},
			"consumer": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceKongConsumerACLGroupCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	consumerACLGroup := getConsumerACLGroupFromResourceData(d)

	createdConsumerACLGroup := getConsumerACLGroupFromResourceData(d)

	response, error := sling.New().BodyJSON(consumerACLGroup).Path("consumers/").Path(consumerACLGroup.Consumer + "/").Post("acls/").ReceiveSuccess(createdConsumerACLGroup)
	if error != nil {
		return fmt.Errorf("Error while creating consumer ACL group.")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setConsumerACLGroupToResourceData(d, createdConsumerACLGroup)

	return nil
}

func resourceKongConsumerACLGroupRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	consumerACLGroup := getConsumerACLGroupFromResourceData(d)

	response, error := sling.New().Path("consumers/").Path(consumerACLGroup.Consumer + "/").Path("acls/").Get(consumerACLGroup.ID).ReceiveSuccess(consumerACLGroup)
	if error != nil {
		return fmt.Errorf("Error while updating consumer ACL group.")
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setConsumerACLGroupToResourceData(d, consumerACLGroup)

	return nil
}

func resourceKongConsumerACLGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	consumerACLGroup := getConsumerACLGroupFromResourceData(d)

	updatedConsumerACLGroup := getConsumerACLGroupFromResourceData(d)

	response, error := sling.New().BodyJSON(consumerACLGroup).Path("consumers/").Path(consumerACLGroup.Consumer + "/").Patch("acls/").Path(consumerACLGroup.ID).ReceiveSuccess(updatedConsumerACLGroup)
	if error != nil {
		return fmt.Errorf("Error while updating consumer ACL group.")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setConsumerACLGroupToResourceData(d, updatedConsumerACLGroup)

	return nil
}

func resourceKongConsumerACLGroupDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	consumerACLGroup := getConsumerACLGroupFromResourceData(d)

	response, error := sling.New().Path("consumers/").Path(consumerACLGroup.Consumer + "/").Path("acls/").Delete(consumerACLGroup.ID).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("Error while deleting consumer ACL group.")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getConsumerACLGroupFromResourceData(d *schema.ResourceData) *ConsumerACLGroup {
	consumerACLGroup := &ConsumerACLGroup{
		ID:       d.Id(),
		Group:    d.Get("group").(string),
		Consumer: d.Get("consumer").(string),
	}

	return consumerACLGroup
}

func setConsumerACLGroupToResourceData(d *schema.ResourceData, consumerACLGroup *ConsumerACLGroup) {
	d.SetId(consumerACLGroup.ID)
	d.Set("group", consumerACLGroup.Group)
	d.Set("consumer", consumerACLGroup.Consumer)
}
