package kong

import (
	"fmt"
	"net/http"

	"github.com/WeKnowSports/terraform-provider-kong/helper"
	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Consumer struct {
	ID       string   `json:"id,omitempty"`
	Username string   `json:"username,omitempty"`
	CustomID string   `json:"custom_id,omitempty"`
	Tags     []string `json:"tags"`
}

func resourceKongConsumer() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongConsumerCreate,
		Read:   resourceKongConsumerRead,
		Update: resourceKongConsumerUpdate,
		Delete: resourceKongConsumerDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The username of the consumer. You must send either this field or custom_id with the request.",
			},

			"custom_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     nil,
				Description: "Field for storing an existing ID for the consumer, useful for mapping Kong with users in your existing database. You must send either this field or username with the request.",
			},

			"tags": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An optional set of strings associated with the Service for grouping and filtering.",
			},
		},
	}
}

func resourceKongConsumerCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	consumer := getConsumerFromResourceData(d)

	createdConsumer := new(Consumer)

	response, error := sling.New().BodyJSON(consumer).Post("consumers/").ReceiveSuccess(createdConsumer)
	if error != nil {
		return fmt.Errorf("error while creating consumer")
	}

	if response.StatusCode == http.StatusConflict {
		return fmt.Errorf("409 Conflict - use terraform import to manage this consumer")
	} else if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setConsumerToResourceData(d, createdConsumer)

	return nil
}

func resourceKongConsumerRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	id := d.Id()
	consumer := new(Consumer)

	response, error := sling.New().Path("consumers/").Get(id).ReceiveSuccess(consumer)
	if error != nil {
		return fmt.Errorf("error while updating consumer")
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setConsumerToResourceData(d, consumer)

	return nil
}

func resourceKongConsumerUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	consumer := getConsumerFromResourceData(d)

	updatedConsumer := new(Consumer)

	response, error := sling.New().BodyJSON(consumer).Patch("consumers/").Path(consumer.ID).ReceiveSuccess(updatedConsumer)
	if error != nil {
		return fmt.Errorf("error while updating consumer")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setConsumerToResourceData(d, updatedConsumer)

	return nil
}

func resourceKongConsumerDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	id := d.Id()

	response, error := sling.New().Delete("consumers/").Path(id).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("error while deleting consumer")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getConsumerFromResourceData(d *schema.ResourceData) *Consumer {
	consumer := &Consumer{
		ID:       d.Id(),
		Username: d.Get("username").(string),
		CustomID: d.Get("custom_id").(string),
		Tags:     helper.ConvertInterfaceArrToStrings(d.Get("tags").([]interface{})),
	}

	return consumer
}

func setConsumerToResourceData(d *schema.ResourceData, consumer *Consumer) {
	d.SetId(consumer.ID)
	d.Set("username", consumer.Username)
	d.Set("custom_id", consumer.CustomID)
	d.Set("tags", consumer.Tags)
}
