package kong

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func ImportConsumerACL(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// This can either be any combination of the consumer name/id and acl group/id separated with a slash.
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Expected a string in the format \"<consumer>/<acl>\" to import.")
	}

	d.Set("consumer", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
