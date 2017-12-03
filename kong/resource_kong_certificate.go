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
	SNIs string `json:"snis,omitempty"`
}

func resourceKongCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongCertificateCreate,
		Read:   resourceKongCertificateRead,
		Update: resourceKongCertificateUpdate,
		Delete: resourceKongCertificateDelete,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "PEM-encoded public certificate of the SSL key pair.",
			},
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "PEM-encoded private key of the SSL key pair.",
			},
			"snis": &schema.Schema{
				Type: schema.TypeString,
				Description: "One or more hostnames to associate with this certificate as an SNI. This is a sugar " +
					"parameter that will, under the hood, create an SNI object and associate it with this certificate " +
					"for your convenience.",
				Optional: true,
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

	response, error := sling.New().BodyJSON(certificate).Path("certificates/").Path(certificate.ID).ReceiveSuccess(updatedCertificate)
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
	stSnis := d.Get("snis").([]interface{})

	snis := make([]string, len(stSnis))
	for i, v := range stSnis {
		snis[i] = v.(string)
	}

	certificate := &Certificate{
		Cert: d.Get("cert").(string),
		Key:  d.Get("key").(string),
		SNIs: strings.Join(snis, ","),
	}

	if id, ok := d.GetOk("id"); ok {
		certificate.ID = id.(string)
	}

	return certificate
}

func setCertificateToResourceData(d *schema.ResourceData, certificate *Certificate) {
	d.SetId(certificate.ID)
	d.Set("cert", certificate.Cert)
	d.Set("key", certificate.Key)

	snis := make([]interface{}, len(certificate.SNIs))

	for i, v := range strings.Split(certificate.SNIs, ",") {
		snis[i] = v
	}

	d.Set("snis", snis)
}
