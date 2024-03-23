package keep

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
)

func dataSourceWorkflows() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadWorkflow,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the workflow.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the workflow.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the workflow.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who created the workflow.",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the workflow was created.",
			},
			"triggers": {
				Type: schema.TypeString,
				//Elem:        false,
				Computed:    true,
				Description: "The triggers of the workflow.",
			},
			"interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The interval of the workflow.",
			},
			"last_execution_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the workflow was last executed.",
			},
			"last_execution_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the last execution of the workflow.",
			},
			"keep_providers": {
				Type: schema.TypeString,
				//Elem:        false,
				Computed:    true,
				Description: "The providers of the workflow.",
			},
			"workflow_raw_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the raw workflow.",
			},
			"workflow_raw": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The raw workflow.",
			},
			"revision": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The revision of the workflow.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the workflow was last updated.",
			},
			"invalid": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The invalid status of the workflow.",
			},
		},
	}
}

func dataSourceReadWorkflow(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	id := d.Get("id").(string)

	req, err := http.NewRequest("GET", client.HostURL+"/workflows/", nil)

	body, err := client.doReq(req)
	if err != nil {
		return diag.Errorf("cannot send request: %s", err)
	}

	var response []map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return diag.Errorf("cannot unmarshal response: %s", err)
	}

	for _, workflow := range response {
		if workflow["id"] == id {
			d.SetId(workflow["id"].(string))
			d.Set("name", workflow["name"])
			d.Set("description", workflow["description"])
			d.Set("created_by", workflow["created_by"])
			d.Set("creation_time", workflow["creation_time"])
			d.Set("triggers", workflow["triggers"])
			d.Set("interval", workflow["interval"])
			d.Set("last_execution_time", workflow["last_execution_time"])
			d.Set("last_execution_status", workflow["last_execution_status"])
			d.Set("keep_providers", workflow["providers"])
			d.Set("workflow_raw_id", workflow["workflow_raw_id"])
			d.Set("workflow_raw", workflow["workflow_raw"])
			d.Set("revision", workflow["revision"])
			d.Set("last_updated", workflow["last_updated"])
			d.Set("invalid", workflow["invalid"])
			break
		}
	}

	return nil
}
