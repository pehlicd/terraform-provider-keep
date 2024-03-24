package keep

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strconv"
)

type Mapping struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	FileName    string   `json:"file_name"`
	Matchers    []string `json:"matchers"`
	Attributes  []string `json:"attributes"`
	CreatedAt   string   `json:"created_at"`
	CreatedBy   string   `json:"created_by"`
}

func dataSourceMapping() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadMapping,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the mapping",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the mapping",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the mapping",
			},
			"file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the mapping file",
			},
			"matchers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of matchers",
			},
			"attributes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of attributes",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the mapping",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creator of the mapping",
			},
		},
	}
}

func dataSourceReadMapping(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	id := d.Get("id").(int)

	req, err := http.NewRequest("GET", client.HostURL+"/mapping/", nil)

	body, err := client.doReq(req)
	if err != nil {
		return diag.Errorf("cannot send request: %s", err)
	}

	var response []Mapping
	err = json.Unmarshal(body, &response)
	if err != nil {
		return diag.Errorf("cannot unmarshal response: %s", err)
	}

	for _, mapping := range response {
		if mapping.ID == id {
			d.SetId(strconv.Itoa(id))
			d.Set("id", strconv.Itoa(id))
			d.Set("name", mapping.Name)
			d.Set("description", mapping.Description)
			d.Set("file_name", mapping.FileName)
			d.Set("matchers", mapping.Matchers)
			d.Set("attributes", mapping.Attributes)
			d.Set("created_at", mapping.CreatedAt)
			d.Set("created_by", mapping.CreatedBy)
			break
		}
	}

	return nil
}
