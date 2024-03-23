package keep

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"net/http"
	"os"
)

func resourceWorkflows() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateWorkflow,
		ReadContext:   resourceReadWorkflow,
		UpdateContext: resourceUpdateWorkflow,
		DeleteContext: resourceDeleteWorkflow,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"workflow_file_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path of the workflow file",
			},
		},
	}
}

func resourceCreateWorkflow(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	workflowFilePath := d.Get("workflow_file_path").(string)

	// read file from workflowFilePath it should be a file path and yaml file
	fInfo, err := os.Stat(workflowFilePath)
	if err != nil {
		return diag.Errorf("workflow file not found: %s", workflowFilePath)
	} else if fInfo.IsDir() {
		return diag.Errorf("workflow file is a directory: %s", workflowFilePath)
	}

	// open file
	file, err := os.OpenFile(workflowFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return diag.Errorf("cannot open file: %s", workflowFilePath)
	}
	defer file.Close()

	// read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return diag.Errorf("cannot read file content: %s", workflowFilePath)
	}

	// create new request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/workflows", client.HostURL), bytes.NewReader(fileContent))
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

	// set id
	d.SetId(response["workflow_id"].(string))

	return nil
}

func resourceDeleteWorkflow(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	id := d.Id()

	// create new request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/workflows/%s", client.HostURL, id), nil)
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

func resourceUpdateWorkflow(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceCreateWorkflow(ctx, d, m)
}

func resourceReadWorkflow(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	id := d.Id()

	// create new request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/workflows/", client.HostURL), nil)
	if err != nil {
		return diag.Errorf("cannot create request: %s", err)
	}

	// send request
	body, err := client.doReq(req)
	if err != nil {
		return diag.Errorf("cannot send request: %s", err)
	}

	// unmarshal response
	var response []map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return diag.Errorf("cannot unmarshal response: %s", err)
	}

	// check if workflow exists
	for _, workflow := range response {
		if workflow["id"] == id {
			return nil
		}
	}

	return nil
}
