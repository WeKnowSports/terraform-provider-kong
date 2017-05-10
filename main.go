package main

import (
	"github.com/hashicorp/terraform/plugin"
	"./kong"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kong.Provider,
	})
}
