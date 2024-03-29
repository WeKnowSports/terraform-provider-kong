package kong

import (
	"fmt"
	"net/http"

	"github.com/WeKnowSports/terraform-provider-kong/helper"
	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type KeyAuthCredential struct {
	ID       string   `json:"id,omitempty"`
	Key      string   `json:"key,omitempty"`
	Consumer string   `json:"-"`
	TTL      int      `json:"ttl,omitempty"`
	Tags     []string `json:"tags"`
}

func resourceKongKeyAuthCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongKeyAuthCredentialCreate,
		Read:   resourceKongKeyAuthCredentialRead,
		Update: resourceKongKeyAuthCredentialUpdate,
		Delete: resourceKongKeyAuthCredentialDelete,

		Importer: &schema.ResourceImporter{
			State: ImportConsumerCredential,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Default:     nil,
				Sensitive:   true,
				Description: "The key to use in the Key Authentication.",
			},

			"consumer": {
				Type:     schema.TypeString,
				Required: true,
			},

			"tags": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An optional set of strings associated with the Service for grouping and filtering.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of seconds the key is going to be valid",
			},
		},
	}
}

func resourceKongKeyAuthCredentialCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	keyAuthCredential := getKeyAuthCredentialFromResourceData(d)

	createdKeyAuthCredential := getKeyAuthCredentialFromResourceData(d)

	response, error := sling.New().BodyJSON(keyAuthCredential).Path("consumers/").Path(keyAuthCredential.Consumer + "/").Post("key-auth/").ReceiveSuccess(createdKeyAuthCredential)
	if error != nil {
		return fmt.Errorf("error while creating keyAuthCredential")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setKeyAuthCredentialToResourceData(d, createdKeyAuthCredential)

	return nil
}

func resourceKongKeyAuthCredentialRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	keyAuthCredential := getKeyAuthCredentialFromResourceData(d)

	response, error := sling.New().Path("consumers/").Path(keyAuthCredential.Consumer + "/").Path("key-auth/").Get(keyAuthCredential.ID).ReceiveSuccess(keyAuthCredential)
	if error != nil {
		return fmt.Errorf("error while updating keyAuthCredential")
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setKeyAuthCredentialToResourceData(d, keyAuthCredential)

	return nil
}

func resourceKongKeyAuthCredentialUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	keyAuthCredential := getKeyAuthCredentialFromResourceData(d)

	updatedKeyAuthCredential := getKeyAuthCredentialFromResourceData(d)

	response, error := sling.New().BodyJSON(keyAuthCredential).Path("consumers/").Path(keyAuthCredential.Consumer + "/").Patch("key-auth/").Path(keyAuthCredential.ID).ReceiveSuccess(updatedKeyAuthCredential)
	if error != nil {
		return fmt.Errorf("error while updating keyAuthCredential")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setKeyAuthCredentialToResourceData(d, updatedKeyAuthCredential)

	return nil
}

func resourceKongKeyAuthCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	keyAuthCredential := getKeyAuthCredentialFromResourceData(d)

	response, error := sling.New().Path("consumers/").Path(keyAuthCredential.Consumer + "/").Path("key-auth/").Delete(keyAuthCredential.ID).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("error while deleting keyAuthCredential")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getKeyAuthCredentialFromResourceData(d *schema.ResourceData) *KeyAuthCredential {
	keyAuthCredential := &KeyAuthCredential{
		ID:       d.Id(),
		Key:      d.Get("key").(string),
		Consumer: d.Get("consumer").(string),
		Tags:     helper.ConvertInterfaceArrToStrings(d.Get("tags").([]interface{})),
		TTL:      d.Get("ttl").(int),
	}

	return keyAuthCredential
}

func setKeyAuthCredentialToResourceData(d *schema.ResourceData, keyAuthCredential *KeyAuthCredential) {
	d.SetId(keyAuthCredential.ID)
	d.Set("key", keyAuthCredential.Key)
	d.Set("consumer", keyAuthCredential.Consumer)
	d.Set("tags", keyAuthCredential.Tags)
	d.Set("ttls", keyAuthCredential.TTL)
}
