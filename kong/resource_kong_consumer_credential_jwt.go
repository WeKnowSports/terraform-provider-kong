package kong

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
)

type JWTCredential struct {
	ID           string `json:"id,omitempty"`
	Key          string `json:"key,omitempty"`
	Algorithm    string `json:"algorithm,omitempty"`
	RSAPublicKey string `json:"rsa_public_key,omitempty"`
	Secret       string `json:"secret,omitempty"`
	Consumer     string `json:"-"`
}

func resourceKongJWTCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongJWTCredentialCreate,
		Read:   resourceKongJWTCredentialRead,
		Update: resourceKongJWTCredentialUpdate,
		Delete: resourceKongJWTCredentialDelete,

		Importer: &schema.ResourceImporter{
			State: ImportConsumerCredential,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "TA unique string identifying the credential. If left out, it will be auto-generated.",
				Sensitive: true,
			},

			"algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The algorithm used to verify the token's signature. Can be HS256 or RS256.",
			},

			"rsa_public_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "If algorithm is RS256, the public key (in PEM format) to use to verify the token's signature.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Sensitive: true,
			},

			"secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "If algorithm is HS256, the secret used to sign JWTs for this credential. If left out, will be auto-generated.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "" && d.Get("rsa_public_key") != ""
				},
				Sensitive: true,
			},

			"consumer": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceKongJWTCredentialCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	jwtCredential := getJWTCredentialFromResourceData(d)

	createdJWTCredential := getJWTCredentialFromResourceData(d)

	response, error := sling.New().BodyJSON(jwtCredential).Path("consumers/").Path(jwtCredential.Consumer + "/").Post("jwt/").ReceiveSuccess(createdJWTCredential)
	if error != nil {
		return fmt.Errorf("Error while creating jwtCredential.")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setJWTCredentialToResourceData(d, createdJWTCredential)

	return nil
}

func resourceKongJWTCredentialRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	jwtCredential := getJWTCredentialFromResourceData(d)

	response, error := sling.New().Path("consumers/").Path(jwtCredential.Consumer + "/").Path("jwt/").Get(jwtCredential.ID).ReceiveSuccess(jwtCredential)
	if error != nil {
		return fmt.Errorf("Error while updating jwtCredential.")
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setJWTCredentialToResourceData(d, jwtCredential)

	return nil
}

func resourceKongJWTCredentialUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	jwtCredential := getJWTCredentialFromResourceData(d)

	updatedJWTCredential := getJWTCredentialFromResourceData(d)

	response, error := sling.New().BodyJSON(jwtCredential).Path("consumers/").Path(jwtCredential.Consumer + "/").Patch("jwt/").Path(jwtCredential.ID).ReceiveSuccess(updatedJWTCredential)
	if error != nil {
		return fmt.Errorf("Error while updating jwtCredential.")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setJWTCredentialToResourceData(d, updatedJWTCredential)

	return nil
}

func resourceKongJWTCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	jwtCredential := getJWTCredentialFromResourceData(d)

	response, error := sling.New().Path("consumers/").Path(jwtCredential.Consumer + "/").Path("jwt/").Delete(jwtCredential.ID).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("Error while deleting jwtCredential.")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getJWTCredentialFromResourceData(d *schema.ResourceData) *JWTCredential {
	jwtCredential := &JWTCredential{
		ID:           d.Id(),
		Key:          d.Get("key").(string),
		Algorithm:    d.Get("algorithm").(string),
		RSAPublicKey: d.Get("rsa_public_key").(string),
		Secret:       d.Get("secret").(string),
		Consumer:     d.Get("consumer").(string),
	}

	return jwtCredential
}

func setJWTCredentialToResourceData(d *schema.ResourceData, jwtCredential *JWTCredential) {
	d.SetId(jwtCredential.ID)
	d.Set("key", jwtCredential.Key)
	d.Set("algorithm", jwtCredential.Algorithm)
	d.Set("rsa_public_key", jwtCredential.RSAPublicKey)
	d.Set("secret", jwtCredential.Secret)
	d.Set("consumer", jwtCredential.Consumer)
}
