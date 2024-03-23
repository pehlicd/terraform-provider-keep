package keep

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client struct with Api Key needed to authenticate against keep
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	ApiKey     string
}

// NewClient func creates new client
func NewClient(hostUrl string, apiKey string) *Client {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    hostUrl,
		ApiKey:     apiKey,
	}

	return &c
}

// doReq func does the api requests
func (c *Client) doReq(req *http.Request) ([]byte, error) {
	req.Header.Add("X-API-KEY", c.ApiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	statusOk := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOk {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func ClientConfigurer(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	host, err := url.Parse(d.Get("backend_url").(string))
	if err != nil {
		return nil, diag.Errorf("backend_url was not a valid url: %s", err.Error())
	}

	return NewClient(host.String(), d.Get("api_key").(string)), nil
}
