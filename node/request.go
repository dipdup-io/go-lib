package node

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

type request struct {
	link       *url.URL
	method     string
	bodyReader io.ReadWriter
}

func newRequest(baseURL, uri, method string, query url.Values, body interface{}) (*request, error) {
	req := new(request)
	link, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if uri != "" {
		link.Path = path.Join(link.Path, uri)
	}
	if len(query) > 0 {
		link.RawQuery = query.Encode()
	}
	req.link = link

	if body != nil {
		bodyReader := new(bytes.Buffer)
		if err := json.NewEncoder(bodyReader).Encode(body); err != nil {
			return nil, err
		}
		req.bodyReader = bodyReader
	}

	req.method = method
	return req, nil
}

func newGetRequest(baseURL, uri string, query url.Values) (*request, error) {
	return newRequest(baseURL, uri, http.MethodGet, query, nil)
}

func newPostRequest(baseURL, uri string, query url.Values, body interface{}) (*request, error) {
	return newRequest(baseURL, uri, http.MethodPost, query, body)
}

func (r *request) do(ctx context.Context, client *client) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, r.method, r.link.String(), r.bodyReader)
	if err != nil {
		return nil, errors.Errorf("request.do: %v", err)
	}
	return client.Do(req)
}

func (r *request) checkStatusCode(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RequestError{
			Code: resp.StatusCode,
			Body: resp.Status,
			Err:  err,
		}
	}
	return RequestError{
		Code: resp.StatusCode,
		Body: string(data),
	}
}

func (r *request) doWithJSONResponse(ctx context.Context, client *client, response interface{}) error {
	resp, err := r.do(ctx, client)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := r.checkStatusCode(resp); err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(response)
}

func (r *request) doWithBytesResponse(ctx context.Context, client *client) ([]byte, error) {
	resp, err := r.do(ctx, client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := r.checkStatusCode(resp); err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}
