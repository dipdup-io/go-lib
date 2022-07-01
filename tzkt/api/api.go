package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/dipdup-net/go-lib/tools/crypto"
	"github.com/dipdup-net/go-lib/tzkt/data"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// API -
type API struct {
	url    string
	client *http.Client

	user       string
	privateKey string
}

// New -
func New(baseURL string, opts ...Option) *API {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	api := &API{
		url: baseURL,
		client: &http.Client{
			Timeout:   time.Minute,
			Transport: t,
		},
	}
	for i := range opts {
		opts[i](api)
	}

	return api
}

func (tzkt *API) get(ctx context.Context, endpoint string, args map[string]string, withAuth bool) (*http.Response, error) {
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

	if withAuth {
		if err := tzkt.auth(req); err != nil {
			return nil, err
		}
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

func (tzkt *API) json(ctx context.Context, endpoint string, args map[string]string, withAuth bool, output interface{}) error {
	resp, err := tzkt.get(ctx, endpoint, args, withAuth)
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
		if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("%s: %s %v", resp.Status, endpoint, args))
	}
}

func (tzkt *API) count(ctx context.Context, endpoint string) (uint64, error) {
	resp, err := tzkt.get(ctx, endpoint, nil, false)
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
func (tzkt *API) GetHead(ctx context.Context) (head data.Head, err error) {
	err = tzkt.json(ctx, "/v1/head", nil, false, &head)
	return
}

func (tzkt *API) auth(request *http.Request) error {
	if tzkt.privateKey == "" || tzkt.user == "" {
		return errors.Errorf("you have to set auth data")
	}
	key, err := crypto.NewKeyFromBase58(tzkt.privateKey)
	if err != nil {
		return err
	}
	nonceString := fmt.Sprintf("%d", time.Now().UTC().UnixMilli())

	sign, err := key.Sign([]byte(nonceString))
	if err != nil {
		return err
	}

	if !key.Verify([]byte(nonceString), sign.Bytes()) {
		return errors.Errorf("invalid signature")
	}

	signature, err := sign.Base58()
	if err != nil {
		return err
	}

	request.Header.Add("X-TZKT-USER", tzkt.user)
	request.Header.Add("X-TZKT-NONCE", nonceString)
	request.Header.Add("X-TZKT-SIGNATURE", signature)

	return nil
}
