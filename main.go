package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/pehlicd/terraform-provider-keep/keep"
)

//go:generate tfplugindocs
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: keep.Provider,
	})
}
