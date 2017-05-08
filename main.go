package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/ContainerLabs/terraform-provider-kong/kong"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kong.Provider,
	})
}
