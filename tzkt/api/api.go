package api

import (
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

func (tzkt *API) get(endpoint string, args map[string]string) (*http.Response, error) {
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

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	return tzkt.client.Do(req)
}

func (tzkt *API) json(endpoint string, args map[string]string, output interface{}) error {
	resp, err := tzkt.get(endpoint, args)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return json.NewDecoder(resp.Body).Decode(output)
	}

	return errors.New(resp.Status)
}

func (tzkt *API) count(endpoint string) (uint64, error) {
	resp, err := tzkt.get(endpoint, nil)
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
func (tzkt *API) GetHead() (head Head, err error) {
	err = tzkt.json("/v1/head", nil, &head)
	return
}
