package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CACertificate struct {
	ID         string   `json:"id,omitempty"`
	Cert       string   `json:"cert,omitempty"`
	CertDigest string   `json:"cert_digest,omitempty"`
	Tags       []string `json:"tags"`
}

func resourceKongCACertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongCACertificateCreate,
		Read:   resourceKongCACertificateRead,
		Update: resourceKongCACertificateUpdate,
		Delete: resourceKongCACertificateDelete,

		Schema: map[string]*schema.Schema{
			"cert": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "PEM-encoded public certificate of the CA",
			},

			"cert_digest": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SHA256 hex digest of the public certificate",
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

func resourceKongCACertificateCreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	caCertificate := getCACertificateFromResourceData(d)

	createdCACertificate := getCACertificateFromResourceData(d)

	response, error := sling.New().BodyJSON(caCertificate).Post("ca_certificates/").ReceiveSuccess(createdCACertificate)
	if error != nil {
		return fmt.Errorf("error while creating caCertificate")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setCACertificateToResourceData(d, createdCACertificate)

	return nil
}

func resourceKongCACertificateRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	caCertificate := getCACertificateFromResourceData(d)

	response, error := sling.New().Path("ca_certificates/").Get(caCertificate.ID).ReceiveSuccess(caCertificate)
	if error != nil {
		return fmt.Errorf("error while updating caCertificate")
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setCACertificateToResourceData(d, caCertificate)

	return nil
}

func resourceKongCACertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	caCertificate := getCACertificateFromResourceData(d)

	updatedCACertificate := getCACertificateFromResourceData(d)

	response, error := sling.New().BodyJSON(caCertificate).Path("ca_certificates/").Patch(caCertificate.ID).ReceiveSuccess(updatedCACertificate)
	if error != nil {
		return fmt.Errorf("error while updating caCertificate")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setCACertificateToResourceData(d, updatedCACertificate)

	return nil
}

func resourceKongCACertificateDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	caCertificate := getCACertificateFromResourceData(d)

	response, error := sling.New().Path("ca_certificates/").Delete(caCertificate.ID).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("error while deleting caCertificate")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getCACertificateFromResourceData(d *schema.ResourceData) *CACertificate {
	caCertificate := &CACertificate{
		ID:         d.Id(),
		Cert:       d.Get("cert").(string),
		CertDigest: d.Get("cert_digest").(string),
	}

	return caCertificate
}

func setCACertificateToResourceData(d *schema.ResourceData, caCertificate *CACertificate) {
	d.SetId(caCertificate.ID)
	d.Set("cert", caCertificate.Cert)
	d.Set("cert_digest", caCertificate.CertDigest)
}
