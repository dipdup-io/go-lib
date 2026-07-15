package hasura

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/dipdup-io/go-lib/config"
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

// New creates an API client talking to the Hasura instance at baseURL, sending
// secret as the X-Hasura-Admin-Secret header on every request when non-empty. The
// returned client reuses a single *http.Client with connection pooling tuned for
// concurrent use (up to 100 idle connections, 100 per host) and a one-minute
// per-request timeout.
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
		return "", errors.Wrap(err, "parse url")
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return "", errors.Errorf("invalid scheme: %s", u.Scheme)
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

	return api.client.Do(req) //nolint:gosec
}

// nolint
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

// Health calls Hasura's /healthz endpoint and returns any error from performing the
// request itself (network failure, non-2xx transport error, etc).
func (api *API) Health(ctx context.Context) error {
	resp, err := api.get(ctx, "/healthz", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return errors.Errorf("invalid status code: %d", resp.StatusCode)
}

// AddSource issues a pg_add_source metadata request that (re)configures the source
// named hasura.Source.Name to point at the Postgres database described by cfg,
// replacing any existing connection configuration. hasura.Source.DatabaseHost
// overrides cfg.Host when set, and hasura.Source.IsolationLevel overrides the
// default "read-committed" isolation level.
func (api *API) AddSource(ctx context.Context, hasura *config.Hasura, cfg config.Database) error {
	host := cfg.Host
	if hasura.Source.DatabaseHost != "" {
		host = hasura.Source.DatabaseHost
	}

	databaseUrl := DatabaseUrl(fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", cfg.User, cfg.Password, host, cfg.Port, cfg.Database))

	isolationLevel := "read-committed"
	if hasura.Source.IsolationLevel != "" {
		isolationLevel = hasura.Source.IsolationLevel
	}

	req := Request{
		Type: "pg_add_source",
		Args: map[string]interface{}{
			"name": hasura.Source.Name,
			"configuration": Configuration{
				ConnectionInfo: ConnectionInfo{
					DatabaseUrl:           databaseUrl,
					UsePreparedStatements: hasura.Source.UsePreparedStatements,
					IsolationLevel:        isolationLevel,
				},
			},
			"replace_configuration": true,
		},
	}
	err := api.post(ctx, "/v1/metadata", nil, req, nil)
	return err
}

// ExportMetadata fetches the full metadata currently applied to the Hasura instance
// via the versioned export_metadata metadata API (v2).
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

// ReplaceMetadata replaces Hasura's entire metadata with data via the versioned
// replace_metadata metadata API (v2), allowing inconsistent metadata so that
// partially invalid state does not block the call.
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

// TrackTable starts tracking the Postgres table name from source in Hasura via
// pg_track_table, exposing it through the GraphQL API.
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

// CustomConfiguration sends conf as-is to Hasura's /v1/metadata endpoint, allowing
// arbitrary metadata API requests — such as those produced by ReadCustomConfigs or
// IterateCustomConfigs — that are not otherwise wrapped by a dedicated API method.
func (api *API) CustomConfiguration(ctx context.Context, conf interface{}) error {
	return api.post(ctx, "/v1/metadata", nil, conf, nil)
}

// CreateSelectPermissions grants role select access to table (in source) via
// pg_create_select_permission, restricted to the columns and row filter described
// by perm.
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

// DropSelectPermissions removes any select permission previously granted to role on
// table (in source) via pg_drop_select_permission.
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

// CreateRestEndpoint exposes queryName from collectionName as a GET REST endpoint
// named name, reachable at url, via create_rest_endpoint.
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
