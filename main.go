package main

import (
	"github.com/WeKnowSports/terraform-provider-kong/kong"
    "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kong.Provider,
	})
}
