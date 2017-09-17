package kong

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func ImportConsumerCredential(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Expected a string in the format \"<consumer_id>/<credential_id>\" to import.")
	}

	d.Set("consumer", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
