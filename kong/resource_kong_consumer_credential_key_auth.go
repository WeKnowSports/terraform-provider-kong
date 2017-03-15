package kong

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

type KeyAuthCredential struct {
	ID       string `json:"id,omitempty"`
	Key      string `json:"key,omitempty"`
	Consumer string `json:"-"`
}

func resourceKongKeyAuthCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongKeyAuthCredentialCreate,
		Read:   resourceKongKeyAuthCredentialRead,
		Update: resourceKongKeyAuthCredentialUpdate,
		Delete: resourceKongKeyAuthCredentialDelete,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The key to use in the Key Authentication.",
			},

			"consumer": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceKongKeyAuthCredentialCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	keyAuthCredential := getKeyAuthCredentialFromResourceData(d)
	log.Print("KEY AUTH")
	log.Print(spew.Sdump(keyAuthCredential))

	createdKeyAuthCredential := getKeyAuthCredentialFromResourceData(d)

	response, error := sling.New().BodyJSON(keyAuthCredential).Path("consumers/").Path(keyAuthCredential.Consumer + "/").Post("key-auth/").ReceiveSuccess(createdKeyAuthCredential)
	log.Print("KEY AUTH RESPONSE")
	log.Print(spew.Sdump(response))
	if error != nil {
		return fmt.Errorf("Error while creating keyAuthCredential.")
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
		return fmt.Errorf("Error while updating keyAuthCredential.")
	}

	if response.StatusCode != http.StatusOK {
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
		return fmt.Errorf("Error while updating keyAuthCredential.")
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
		return fmt.Errorf("Error while deleting keyAuthCredential.")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getKeyAuthCredentialFromResourceData(d *schema.ResourceData) *KeyAuthCredential {
	keyAuthCredential := &KeyAuthCredential{
		Key:      d.Get("key").(string),
		Consumer: d.Get("consumer").(string),
	}

	if id, ok := d.GetOk("id"); ok {
		keyAuthCredential.ID = id.(string)
	}

	return keyAuthCredential
}

func setKeyAuthCredentialToResourceData(d *schema.ResourceData, keyAuthCredential *KeyAuthCredential) {
	d.SetId(keyAuthCredential.ID)
	d.Set("key", keyAuthCredential.Key)
	d.Set("consumer", keyAuthCredential.Consumer)
}
