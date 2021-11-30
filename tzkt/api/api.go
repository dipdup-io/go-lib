package api

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// API -
type API struct {
	url    string
	client *http.Client
}

// New -
func New(baseURL string) *API {
	return &API{
		url:    baseURL,
		client: http.DefaultClient,
	}
}

func (tzkt *API) get(ctx context.Context, endpoint string, args map[string]string) (*http.Response, error) {
	u, err := url.Parse(tzkt.url)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, endpoint)

	values := u.Query()
	for key, value := range args {
		values.Add(key, value)
	}
	u.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	return tzkt.client.Do(req)
}

func (tzkt *API) post(ctx context.Context, endpoint string, args map[string]string, body interface{}) (*http.Response, error) {
	u, err := url.Parse(tzkt.url)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, endpoint)

	values := u.Query()
	for key, value := range args {
		values.Add(key, value)
	}
	u.RawQuery = values.Encode()

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	return tzkt.client.Do(req)
}

func (tzkt *API) json(ctx context.Context, endpoint string, args map[string]string, output interface{}) error {
	resp, err := tzkt.get(ctx, endpoint, args)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return json.NewDecoder(resp.Body).Decode(output)
	case http.StatusNoContent:
		return nil
	default:
		return errors.New(fmt.Sprintf("%s: %s %v", resp.Status, endpoint, args))
	}
}

func (tzkt *API) count(ctx context.Context, endpoint string) (uint64, error) {
	resp, err := tzkt.get(ctx, endpoint, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(string(data), 10, 64)
}

// GetHead -
func (tzkt *API) GetHead(ctx context.Context) (head Head, err error) {
	err = tzkt.json(ctx, "/v1/head", nil, &head)
	return
}
