package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

type Certificate struct {
	ID   string `json:"id,omitempty"`
	Cert string `json:"cert,omitempty"`
	Key  string `json:"key,omitempty"`
}

func resourceKongCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongCertificateCreate,
		Read:   resourceKongCertificateRead,
		Update: resourceKongCertificateUpdate,
		Delete: resourceKongCertificateDelete,

		Schema: map[string]*schema.Schema{
			"cert": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PEM-encoded public certificate of the SSL key pair.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "PEM-encoded private key of the SSL key pair.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
			},
		},
	}
}

func resourceKongCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	certificate := getCertificateFromResourceData(d)

	createdCertificate := getCertificateFromResourceData(d)

	response, error := sling.New().BodyJSON(certificate).Post("certificates/").ReceiveSuccess(createdCertificate)
	if error != nil {
		return fmt.Errorf("Error while creating certificate.")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setCertificateToResourceData(d, createdCertificate)

	return nil
}

func resourceKongCertificateRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	certificate := getCertificateFromResourceData(d)

	response, error := sling.New().Path("certificates/").Get(certificate.ID).ReceiveSuccess(certificate)
	if error != nil {
		return fmt.Errorf("Error while updating certificate.")
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setCertificateToResourceData(d, certificate)

	return nil
}

func resourceKongCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	certificate := getCertificateFromResourceData(d)

	updatedCertificate := getCertificateFromResourceData(d)

	response, error := sling.New().BodyJSON(certificate).Path("certificates/").Patch(certificate.ID).ReceiveSuccess(updatedCertificate)
	if error != nil {
		return fmt.Errorf("Error while updating certificate")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setCertificateToResourceData(d, updatedCertificate)

	return nil
}

func resourceKongCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	certificate := getCertificateFromResourceData(d)

	response, error := sling.New().Path("certificates/").Delete(certificate.ID).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("Error while deleting certificate.")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getCertificateFromResourceData(d *schema.ResourceData) *Certificate {
	certificate := &Certificate{
		ID:   d.Id(),
		Cert: d.Get("cert").(string),
		Key:  d.Get("key").(string),
	}

	return certificate
}

func setCertificateToResourceData(d *schema.ResourceData, certificate *Certificate) {
	d.SetId(certificate.ID)
	d.Set("cert", certificate.Cert)
	d.Set("key", certificate.Key)
}
