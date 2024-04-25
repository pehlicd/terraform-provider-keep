package keep

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func resourceMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateMapping,
		ReadContext:   resourceReadMapping,
		UpdateContext: resourceUpdateMapping,
		DeleteContext: resourceDeleteMapping,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the mapping",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the mapping",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the mapping",
			},
			"matchers": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of matchers",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Priority of the mapping",
			},
			"mapping_file_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path of the mapping file",
			},
		},
	}
}

func resourceCreateMapping(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	mappingFilePath := d.Get("mapping_file_path").(string)

	// read file from mappingFilePath it should be a file path and csv file
	fInfo, err := os.Stat(mappingFilePath)
	if err != nil {
		return diag.Errorf("mapping file not found: %s", mappingFilePath)
	} else if fInfo.IsDir() {
		return diag.Errorf("mapping file is a directory: %s", mappingFilePath)
	}

	file, err := os.OpenFile(mappingFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return diag.Errorf("cannot open file: %s", mappingFilePath)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return diag.Errorf("Error reading CSV file: %s", err)
	}

	headers := records[0]
	records = records[1:]

	rows := make([]map[string]string, len(records))
	for i, record := range records {
		row := make(map[string]string)
		for j, cell := range record {
			row[headers[j]] = cell
		}
		rows[i] = row
	}

	matchers := d.Get("matchers").(*schema.Set).List()
	//convert matchers to string array
	matchersStr := make([]string, len(matchers))
	for i, matcher := range matchers {
		matchersStr[i] = matcher.(string)
	}

	body := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"matchers":    matchersStr,
		"priority":    d.Get("priority").(int),
		"rows":        rows,
		"file_name":   fInfo.Name(),
	}

	// marshal body
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return diag.Errorf("cannot request marshal body: %s", err)
	}

	// create mapping
	req, err := http.NewRequest("POST", client.HostURL+"/mapping/", strings.NewReader(string(bodyBytes)))
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

func resourceReadMapping(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	id := d.Id()

	req, err := http.NewRequest("GET", client.HostURL+"/mapping/", nil)
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

	for _, mapping := range response {
		if mapping["id"] == idFloat {
			d.SetId(id)
			d.Set("name", mapping["name"])
			d.Set("description", mapping["description"])
			d.Set("matchers", mapping["matchers"])
			d.Set("priority", mapping["priority"])
			d.Set("file_name", mapping["file_name"])
			break
		}
	}

	return nil
}

func resourceUpdateMapping(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	id := d.Id()

	req, err := http.NewRequest("GET", client.HostURL+"/mapping/", nil)
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

	if len(response) == 0 {
		return diag.Errorf("no mapping found")
	}

	// check if mapping exists
	for _, mapping := range response {
		if mapping["id"] == id {
			// update mapping
			mappingFilePath := d.Get("mapping_file_path").(string)

			// read file from mappingFilePath it should be a file path and csv file
			fInfo, err := os.Stat(mappingFilePath)
			if err != nil {
				return diag.Errorf("mapping file not found: %s", mappingFilePath)
			} else if fInfo.IsDir() {
				return diag.Errorf("mapping file is a directory: %s", mappingFilePath)
			}

			file, err := os.OpenFile(mappingFilePath, os.O_RDONLY, 0644)
			if err != nil {
				return diag.Errorf("cannot open file: %s", mappingFilePath)
			}
			defer file.Close()

			reader := csv.NewReader(file)
			records, err := reader.ReadAll()
			if err != nil {
				return diag.Errorf("Error reading CSV file: %s", err)
			}

			headers := records[0]
			records = records[1:]

			rows := make([]map[string]string, len(records))
			for i, record := range records {
				row := make(map[string]string)
				for j, cell := range record {
					row[headers[j]] = cell
				}
				rows[i] = row
			}

			matchers := d.Get("matchers").(*schema.Set).List()
			//convert matchers to string array
			matchersStr := make([]string, len(matchers))
			for i, matcher := range matchers {
				matchersStr[i] = matcher.(string)
			}

			reqBody := map[string]interface{}{
				"name":        d.Get("name").(string),
				"description": d.Get("description").(string),
				"matchers":    matchersStr,
				"priority":    d.Get("priority").(int),
				"rows":        rows,
				"file_name":   fInfo.Name(),
			}

			mappingRule := map[string]interface{}{
				"name":        mapping["name"],
				"description": mapping["description"],
				"matchers":    mapping["matchers"],
				"priority":    mapping["priority"],
				"rows":        mapping["rows"],
				"file_name":   mapping["file_name"],
			}

			// compare request body with existing mapping if one or more fields are different, update mapping
			if !reflect.DeepEqual(reqBody, mappingRule) {
				// no changes
				return nil
			}

			bodyBytes, err := json.Marshal(reqBody)
			if err != nil {
				return diag.Errorf("cannot request marshal body: %s", err)
			}

			updateReq, err := http.NewRequest("PUT", client.HostURL+"/mapping/", strings.NewReader(string(bodyBytes)))
			if err != nil {
				return diag.Errorf("cannot create request: %s", err)
			}

			// send request
			respBody, err := client.doReq(updateReq)
			if err != nil {
				return diag.Errorf("cannot send request: %s", err)
			}

			// unmarshal response
			var response map[string]interface{}
			err = json.Unmarshal(respBody, &response)
			if err != nil {
				return diag.Errorf("cannot unmarshal response: %s", err)
			}

			d.SetId(response["id"].(string))
			d.Set("id", response["id"].(string))
			break
		}
	}

	return nil
}

func resourceDeleteMapping(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	id := d.Id()

	req, err := http.NewRequest("DELETE", client.HostURL+"/mapping/"+id, nil)
	if err != nil {
		return diag.Errorf("cannot create request: %s", err)
	}

	_, err = client.doReq(req)
	if err != nil {
		return diag.Errorf("cannot send request: %s", err)
	}

	return nil
}
