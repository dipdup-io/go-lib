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
		return errors.Errorf("Invalid status code: %s", resp.Status)
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
func (api *API) ExportMetadata(data interface{}) (ExportMetadataResponse, error) {
	req := request{
		Type: "export_metadata",
		Args: data,
	}
	var resp ExportMetadataResponse
	err := api.post("/v1/query", nil, req, &resp)
	return resp, err
}

// ReplaceMetadata -
func (api *API) ReplaceMetadata(data interface{}) error {
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
