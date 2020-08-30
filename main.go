package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

    "github.com/rudexi/terraform-provider-isilon/isilon"
)

func main() {

    plugin.Serve(&plugin.ServeOpts{
            ProviderFunc: func() *schema.Provider {
                return isilon.Provider()
            },
    })

}
