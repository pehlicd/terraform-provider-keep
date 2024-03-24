package keep

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider for Keep
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"backend_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Keep backend url",
				DefaultFunc: schema.EnvDefaultFunc("KEEP_BACKEND_URL", nil),
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Keep API Key",
				DefaultFunc: schema.EnvDefaultFunc("KEEP_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"keep_provider": resourceProvider(),
			"keep_workflow": resourceWorkflow(),
			"keep_mapping":  resourceMapping(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"keep_workflow": dataSourceWorkflows(),
			"keep_mapping":  dataSourceMapping(),
		},
		ConfigureContextFunc: ClientConfigurer,
	}
}
