package hasura

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/dipdup-net/go-lib/config"
	jsoniter "github.com/json-iterator/go"

	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// API -
type API struct {
	baseURL string
	secret  string

	client *http.Client
}

// New -
func New(baseURL, secret string) *API {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &API{baseURL, secret, &http.Client{
		Timeout:   time.Minute,
		Transport: t,
	}}
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

func (api *API) get(ctx context.Context, endpoint string, args map[string]string) (*http.Response, error) {
	url, err := api.buildURL(endpoint, args)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return api.client.Do(req)
}

//nolint
func (api *API) post(ctx context.Context, endpoint string, args map[string]string, body interface{}, output interface{}) error {
	url, err := api.buildURL(endpoint, args)
	if err != nil {
		return err
	}

	postBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(postBody))
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
func (api *API) Health(ctx context.Context) error {
	resp, err := api.get(ctx, "/healthz", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return err
}

// AddSource -
func (api *API) AddSource(ctx context.Context, hasura *config.Hasura, cfg config.Database) error {
	req := Request{
		Type: "pg_add_source",
		Args: map[string]interface{}{
			"name": hasura.Source,
			"configuration": Configuration{
				ConnectionInfo: ConnectionInfo{
					DatabaseUrl:           DatabaseUrl(fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)),
					UsePreparedStatements: true,
					IsolationLevel:        "read-committed",
				},
			},
			"replace_configuration": true,
		},
	}
	err := api.post(ctx, "/v1/metadata", nil, req, nil)
	return err
}

// ExportMetadata -
func (api *API) ExportMetadata(ctx context.Context) (Metadata, error) {
	req := versionedRequest{
		Type:    "export_metadata",
		Version: 2,
		Args:    map[string]interface{}{},
	}
	var resp Metadata
	err := api.post(ctx, "/v1/metadata", nil, req, &resp)
	return resp, err
}

// ReplaceMetadata -
func (api *API) ReplaceMetadata(ctx context.Context, data *Metadata) error {
	req := versionedRequest{
		Type:    "replace_metadata",
		Version: 2,
		Args: map[string]any{
			"metadata":                    data,
			"allow_inconsistent_metadata": true,
		},
	}
	return api.post(ctx, "/v1/metadata", nil, req, nil)
}

// TrackTable -
func (api *API) TrackTable(ctx context.Context, name string, source string) error {
	req := Request{
		Type: "pg_track_table",
		Args: map[string]string{
			"table":  name,
			"source": source,
		},
	}
	return api.post(ctx, "/v1/metadata", nil, req, nil)
}

// CustomConfiguration
func (api *API) CustomConfiguration(ctx context.Context, conf interface{}) error {
	return api.post(ctx, "/v1/metadata", nil, conf, nil)
}

// CreateSelectPermissions - A select permission is used to restrict access to only the specified columns and rows.
func (api *API) CreateSelectPermissions(ctx context.Context, table, source string, role string, perm Permission) error {
	req := Request{
		Type: "pg_create_select_permission",
		Args: map[string]interface{}{
			"table":      table,
			"role":       role,
			"permission": perm,
			"source":     source,
		},
	}
	return api.post(ctx, "/v1/metadata", nil, req, nil)
}

// DropSelectPermissions -
func (api *API) DropSelectPermissions(ctx context.Context, table, source string, role string) error {
	req := Request{
		Type: "pg_drop_select_permission",
		Args: map[string]interface{}{
			"table":  table,
			"role":   role,
			"source": source,
		},
	}
	return api.post(ctx, "/v1/metadata", nil, req, nil)
}

// CreateRestEndpoint -
func (api *API) CreateRestEndpoint(ctx context.Context, name, url, queryName, collectionName string) error {
	req := Request{
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
	return api.post(ctx, "/v1/metadata", nil, req, nil)
}
