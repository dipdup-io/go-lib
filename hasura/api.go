package hasura

import (
	"bytes"
	"net/http"
	"net/url"
	"path"

	jsoniter "github.com/json-iterator/go"

	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// API -
type API struct {
	baseURL string
	secret  string

	client http.Client
}

// New -
func New(baseURL, secret string) *API {
	return &API{baseURL, secret, *http.DefaultClient}
}

func (api *API) buildURL(endpoint string, args map[string]string) (string, error) {
	u, err := url.Parse(api.baseURL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, endpoint)

	values := u.Query()
	for key, value := range args {
		values.Add(key, value)
	}
	u.RawQuery = values.Encode()
	return u.String(), nil
}

func (api *API) get(endpoint string, args map[string]string) (*http.Response, error) {
	url, err := api.buildURL(endpoint, args)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return api.client.Do(req)
}

//nolint
func (api *API) post(endpoint string, args map[string]string, body interface{}, output interface{}) error {
	url, err := api.buildURL(endpoint, args)
	if err != nil {
		return err
	}

	postBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(postBody))
	if err != nil {
		return err
	}

	if api.secret != "" {
		req.Header.Add("X-Hasura-Admin-Secret", api.secret)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiError APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return errors.Wrap(err, "Hasura's response decoding error")
		}
		return apiError
	}

	if output == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(output)
}

// Health
func (api *API) Health() error {
	resp, err := api.get("/healthz", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return err
}

// ExportMetadata -
func (api *API) ExportMetadata(data *Metadata) (ExportMetadataResponse, error) {
	req := request{
		Type: "export_metadata",
		Args: data,
	}
	var resp ExportMetadataResponse
	err := api.post("/v1/query", nil, req, &resp)
	return resp, err
}

// ReplaceMetadata -
func (api *API) ReplaceMetadata(data *Metadata) error {
	req := request{
		Type: "replace_metadata",
		Args: data,
	}
	var resp replaceMetadataResponse
	if err := api.post("/v1/query", nil, req, &resp); err != nil {
		return err
	}
	if resp.Message == "success" {
		return nil
	}
	return errors.Errorf("Can't replace hasura's metadata: %s", resp.Message)
}

// TrackTable -
func (api *API) TrackTable(schema, name string) error {
	req := request{
		Type: "track_table",
		Args: map[string]string{
			"schema": schema,
			"name":   name,
		},
	}
	return api.post("/v1/query", nil, req, nil)
}

// CreateSelectPermissions - A select permission is used to restrict access to only the specified columns and rows.
func (api *API) CreateSelectPermissions(table, role string, perm Permission) error {
	req := request{
		Type: "create_select_permission",
		Args: map[string]interface{}{
			"table":      table,
			"role":       role,
			"permission": perm,
		},
	}
	return api.post("/v1/query", nil, req, nil)
}

// DropSelectPermissions -
func (api *API) DropSelectPermissions(table, role string) error {
	req := request{
		Type: "drop_select_permission",
		Args: map[string]interface{}{
			"table": table,
			"role":  role,
		},
	}
	return api.post("/v1/query", nil, req, nil)
}

// CreateRestEndpoint -
func (api *API) CreateRestEndpoint(name, url, queryName, collectionName string) error {
	req := request{
		Type: "create_rest_endpoint",
		Args: map[string]interface{}{
			"name":    name,
			"url":     url,
			"methods": []string{"GET"},
			"definition": map[string]interface{}{
				"query": map[string]interface{}{
					"query_name":      queryName,
					"collection_name": collectionName,
				},
			},
		},
	}
	return api.post("/v1/query", nil, req, nil)
}
