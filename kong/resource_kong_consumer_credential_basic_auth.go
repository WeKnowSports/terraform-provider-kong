package kong

import (
	"fmt"
	"net/http"

	"crypto/sha1"
	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"strings"
)

type BasicAuthCredential struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Consumer string `json:"-"`
}

func resourceKongBasicAuthCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongBasicAuthCredentialCreate,
		Read:   resourceKongBasicAuthCredentialRead,
		Update: resourceKongBasicAuthCredentialUpdate,
		Delete: resourceKongBasicAuthCredentialDelete,

		Importer: &schema.ResourceImporter{
			State: ImportConsumerCredential,
		},

		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username to use in the Basic Authentication.",
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Sensitive:   true,
				Description: "The password to use in the Basic Authentication.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					sha1 := sha1.New()
					io.WriteString(sha1, new)
					io.WriteString(sha1, d.Get("consumer").(string))
					return strings.TrimSpace(old) == fmt.Sprintf("%x", sha1.Sum(nil))
				},
			},

			"consumer": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceKongBasicAuthCredentialCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	basicAuthCredential := getBasicAuthCredentialFromResourceData(d)

	createdBasicAuthCredential := getBasicAuthCredentialFromResourceData(d)

	response, error := sling.New().BodyJSON(basicAuthCredential).Path("consumers/").Path(basicAuthCredential.Consumer + "/").Post("basic-auth/").ReceiveSuccess(createdBasicAuthCredential)
	if error != nil {
		return fmt.Errorf("Error while creating basicAuthCredential.")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setBasicAuthCredentialToResourceData(d, createdBasicAuthCredential)

	return nil
}

func resourceKongBasicAuthCredentialRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	basicAuthCredential := getBasicAuthCredentialFromResourceData(d)

	response, error := sling.New().Path("consumers/").Path(basicAuthCredential.Consumer + "/").Path("basic-auth/").Get(basicAuthCredential.ID).ReceiveSuccess(basicAuthCredential)
	if error != nil {
		return fmt.Errorf("Error while updating basicAuthCredential.")
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setBasicAuthCredentialToResourceData(d, basicAuthCredential)

	return nil
}

func resourceKongBasicAuthCredentialUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	basicAuthCredential := getBasicAuthCredentialFromResourceData(d)

	updatedBasicAuthCredential := getBasicAuthCredentialFromResourceData(d)

	response, error := sling.New().BodyJSON(basicAuthCredential).Path("consumers/").Path(basicAuthCredential.Consumer + "/").Patch("basic-auth/").Path(basicAuthCredential.ID).ReceiveSuccess(updatedBasicAuthCredential)
	if error != nil {
		return fmt.Errorf("Error while updating basicAuthCredential.")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setBasicAuthCredentialToResourceData(d, updatedBasicAuthCredential)

	return nil
}

func resourceKongBasicAuthCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	basicAuthCredential := getBasicAuthCredentialFromResourceData(d)

	response, error := sling.New().Path("consumers/").Path(basicAuthCredential.Consumer + "/").Path("basic-auth/").Delete(basicAuthCredential.ID).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("Error while deleting basicAuthCredential.")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getBasicAuthCredentialFromResourceData(d *schema.ResourceData) *BasicAuthCredential {
	basicAuthCredential := &BasicAuthCredential{
		ID:       d.Id(),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Consumer: d.Get("consumer").(string),
	}

	return basicAuthCredential
}

func setBasicAuthCredentialToResourceData(d *schema.ResourceData, basicAuthCredential *BasicAuthCredential) {
	d.SetId(basicAuthCredential.ID)
	d.Set("username", basicAuthCredential.Username)
	d.Set("password", basicAuthCredential.Password)
	d.Set("consumer", basicAuthCredential.Consumer)
}
