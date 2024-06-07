package keep

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strconv"
	"strings"
)

func resourceExtraction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateExtraction,
		ReadContext:   resourceReadExtraction,
		UpdateContext: resourceUpdateExtraction,
		DeleteContext: resourceDeleteExtraction,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the extraction",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the extraction",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the extraction",
				Default:     "",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Priority of the extraction",
				Default:     0,
			},
			"attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Attribute of the extraction",
			},
			"condition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Condition of the extraction",
				Default:     "",
			},
			"disabled": {
				Type:     schema.TypeBool,
				Required: true,
				Default:  false,
			},
			"regex": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Regex of the extraction",
			},
			"pre": {
				Type:        schema.TypeBool,
				Required:    true,
				Default:     false,
				Description: "Pre of the extraction",
			},
		},
	}
}

func resourceCreateExtraction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	body := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"priority":    d.Get("priority").(int),
		"attribute":   d.Get("attribute").(string),
		"condition":   d.Get("condition").(string),
		"disabled":    d.Get("disabled").(bool),
		"regex":       d.Get("regex").(string),
		"pre":         d.Get("pre").(bool),
	}

	// marshal body
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return diag.Errorf("cannot marshal extraction body: %s", err)
	}

	// create extraction
	req, err := http.NewRequest("POST", client.HostURL+"/extraction/", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return diag.Errorf("cannot create request: %s", err)
	}

	// send request
	respBody, err := client.doReq(req)
	if err != nil {
		return diag.Errorf("cannot send request: %s", err)
	}

	// unmarshal response
	var response map[string]interface{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return diag.Errorf("cannot unmarshal response: %s", err)
	}

	d.SetId(fmt.Sprintf("%f", response["id"]))
	d.Set("id", fmt.Sprintf("%f", response["id"]))

	return nil
}

func resourceReadExtraction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	id := d.Id()

	req, err := http.NewRequest("GET", client.HostURL+"/extraction/", nil)
	if err != nil {
		return diag.Errorf("cannot create request: %s", err)
	}

	body, err := client.doReq(req)
	if err != nil {
		return diag.Errorf("cannot send request: %s", err)
	}

	var response []map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return diag.Errorf("cannot unmarshal response: %s", err)
	}

	idFloat, err := strconv.ParseFloat(id, 64)
	if err != nil {
		return diag.Errorf("cannot parse id: %s", err)
	}

	for _, extraction := range response {
		if extraction["id"] == idFloat {
			d.SetId(id)
			d.Set("name", extraction["name"])
			d.Set("description", extraction["description"])
			d.Set("priority", extraction["priority"])
			d.Set("attribute", extraction["attribute"])
			d.Set("condition", extraction["condition"])
			d.Set("disabled", extraction["disabled"])
			d.Set("regex", extraction["regex"])
			d.Set("pre", extraction["pre"])
			break
		}
	}

	return nil
}

func resourceUpdateExtraction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	id := d.Id()

	// Prepare the payload for the extraction update request
	extractionUpdatePayload := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"priority":    d.Get("priority").(int),
		"attribute":   d.Get("attribute").(string),
		"condition":   d.Get("condition").(string),
		"disabled":    d.Get("disabled").(bool),
		"regex":       d.Get("regex").(string),
		"pre":         d.Get("pre").(bool),
	}

	if !d.HasChange("name") || !d.HasChange("description") || !d.HasChange("priority") || !d.HasChange("attribute") || !d.HasChange("condition") || !d.HasChange("disabled") || !d.HasChange("regex") || !d.HasChange("pre") {
		return nil
	}

	// Marshal the payload
	payload, err := json.Marshal(extractionUpdatePayload)
	if err != nil {
		return diag.Errorf("cannot marshal payload: %s", err)
	}

	// Create a new request
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/extraction/%s", client.HostURL, id), strings.NewReader(string(payload)))
	if err != nil {
		return diag.Errorf("cannot create request: %s", err)
	}

	// Do the request
	body, err := client.doReq(req)
	if err != nil {
		return diag.Errorf("cannot send request: %s", err)
	}

	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return diag.Errorf("cannot parse response: %s", err)
	}

	d.SetId(id)

	return nil
}

func resourceDeleteExtraction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	id := d.Id()

	req, err := http.NewRequest("DELETE", client.HostURL+"/extraction/"+id, nil)
	if err != nil {
		return diag.Errorf("cannot create request: %s", err)
	}

	_, err = client.doReq(req)
	if err != nil {
		return diag.Errorf("cannot send request: %s", err)
	}

	return nil
}
