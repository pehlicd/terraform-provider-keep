package keep

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strings"
)

func resourceProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateProvider,
		ReadContext:   resourceReadProvider,
		UpdateContext: resourceUpdateProvider,
		DeleteContext: resourceDeleteProvider,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the keep provider",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the keep provider",
			},
			"auth_config": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "Configuration of the keep provider authentication",
			},
		},
	}
}

func resourceCreateProvider(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	providerType := d.Get("type").(string)
	providerName := d.Get("name").(string)

	// create new request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/providers", client.HostURL), nil)
	if err != nil {
		return diag.Errorf("cannot create request: %s", err)
	}

	// do request
	body, err := client.doReq(req)
	if err != nil {
		return diag.Errorf("cannot send request: %s", err)
	}

	// parse response
	var providers map[string]interface{}
	err = json.Unmarshal(body, &providers)
	if err != nil {
		return diag.Errorf("cannot parse response: %s", err)
	}

	found := false
	authConfigs := d.Get("auth_config").(map[string]interface{})
	availableProviders := providers["providers"].([]interface{})

	for _, provider := range availableProviders {
		providerMap, ok := provider.(map[string]interface{})
		if ok && providerMap["type"] == providerType {
			// provider is supported
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("provider not found")
	}

	// Prepare the payload for the provider installation request
	providerInstallPayload := map[string]interface{}{
		"provider_id":   providerType,
		"provider_name": providerName,
	}

	// Add the auth config to the payload
	for key, value := range authConfigs {
		providerInstallPayload[key] = value
	}

	// Marshal the payload
	payload, err := json.Marshal(providerInstallPayload)
	if err != nil {
		return diag.Errorf("cannot marshal payload: %s", err)
	}

	// Create a new request
	req, err = http.NewRequest("POST", fmt.Sprintf("%s/providers/install", client.HostURL), strings.NewReader(string(payload)))
	if err != nil {
		return diag.Errorf("cannot create request: %s", err)
	}

	// Do the request
	body, err = client.doReq(req)
	if err != nil {
		return diag.Errorf("cannot send request: %s", err)
	}

	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return diag.Errorf("cannot parse response: %s", err)
	}
	if response == nil {
		return diag.Errorf("couldn't create provider properly, response is nil")
	}

	// Set the ID
	d.SetId(response["id"].(string))
	d.Set("type", providerType)

	return nil
}

func resourceDeleteProvider(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	id := d.Id()
	providerType := d.Get("type").(string)

	// create new request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/providers/%s/%s", client.HostURL, providerType, id), nil)
	if err != nil {
		return diag.Errorf("cannot create request: %s", err)
	}

	// send request
	_, err = client.doReq(req)
	if err != nil {
		return diag.Errorf("cannot send request: %s", err)
	}

	return nil
}

func resourceReadProvider(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	id := d.Id()

	// create new request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/providers/", client.HostURL), nil)
	if err != nil {
		return diag.Errorf("cannot create request: %s", err)
	}

	// send request
	body, err := client.doReq(req)
	if err != nil {
		return diag.Errorf("cannot send request: %s", err)
	}

	// unmarshal response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return diag.Errorf("cannot unmarshal response: %s", err)
	}

	installedProviders := response["installed_providers"].([]interface{})
	for _, provider := range installedProviders {
		if provider.(map[string]interface{})["id"] == id {
			// provider exists
			// in the future we can set the provider data here
			d.SetId(id)
			return nil
		}
	}

	return nil
}

func resourceUpdateProvider(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	id := d.Id()
	providerType := d.Get("type").(string)
	providerName := d.Get("name").(string)
	authConfig := d.Get("auth_config").(map[string]interface{})

	// Prepare the payload for the provider update request
	providerUpdatePayload := map[string]interface{}{
		"provider_id":   providerType,
		"provider_name": providerName,
	}

	// Add the auth config to the payload
	for key, value := range authConfig {
		providerUpdatePayload[key] = value
	}

	// Marshal the payload
	payload, err := json.Marshal(providerUpdatePayload)
	if err != nil {
		return diag.Errorf("cannot marshal payload: %s", err)
	}

	// Create a new request
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/providers/%s", client.HostURL, id), strings.NewReader(string(payload)))
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
