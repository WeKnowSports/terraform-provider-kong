package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SNI struct {
	Name             string      `json:"name,omitempty"`
	SSLCertificateID Certificate `json:"certificate,omitempty"`
}

func resourceKongSNI() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongSNICreate,
		Read:   resourceKongSNIRead,
		Update: resourceKongSNIUpdate,
		Delete: resourceKongSNIDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The SNI name to associate with the given sni.",
			},
			"certificate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id (a UUID) of the certificate with which to associate the SNI hostname.",
			},
		},
	}
}

func resourceKongSNICreate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	sni := getSNIFromResourceData(d)

	createdSNI := getSNIFromResourceData(d)

	response, error := sling.New().BodyJSON(sni).Post("snis/").ReceiveSuccess(createdSNI)
	if error != nil {
		return fmt.Errorf("error while creating SNI")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf(response.Status)
	}

	setSNIToResourceData(d, createdSNI)

	return nil
}

func resourceKongSNIRead(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	sni := getSNIFromResourceData(d)

	response, error := sling.New().Path("snis/").Get(sni.Name).ReceiveSuccess(sni)
	if error != nil {
		return fmt.Errorf("error while updating SNI")
	}

	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setSNIToResourceData(d, sni)

	return nil
}

func resourceKongSNIUpdate(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	sni := getSNIFromResourceData(d)

	updatedSNI := getSNIFromResourceData(d)

	response, error := sling.New().BodyJSON(sni).Path("snis/").Patch(sni.Name).ReceiveSuccess(updatedSNI)
	if error != nil {
		return fmt.Errorf("error while updating SNI")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	setSNIToResourceData(d, updatedSNI)

	return nil
}

func resourceKongSNIDelete(d *schema.ResourceData, meta interface{}) error {
	sling := meta.(*sling.Sling)

	sni := getSNIFromResourceData(d)

	response, error := sling.New().Path("snis/").Delete(sni.Name).ReceiveSuccess(nil)
	if error != nil {
		return fmt.Errorf("error while deleting SNI")
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(response.Status)
	}

	return nil
}

func getSNIFromResourceData(d *schema.ResourceData) *SNI {
	sni := &SNI{
		Name: d.Get("name").(string),
		SSLCertificateID: Certificate{
			ID: d.Get("certificate").(string),
		},
	}

	return sni
}

func setSNIToResourceData(d *schema.ResourceData, sni *SNI) {
	d.SetId(sni.Name)
	d.Set("name", sni.Name)
	d.Set("certificate", sni.SSLCertificateID)
}
