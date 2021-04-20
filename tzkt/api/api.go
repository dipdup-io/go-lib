package api

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

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

func (tzkt *API) get(endpoint string, args map[string]string, output interface{}) error {
	u, err := url.Parse(tzkt.url)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, endpoint)

	values := u.Query()
	for key, value := range args {
		values.Add(key, value)
	}
	u.RawQuery = values.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := tzkt.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return json.NewDecoder(resp.Body).Decode(output)
	}

	return errors.New(resp.Status)
}

// GetHead -
func (tzkt *API) GetHead() (head Head, err error) {
	err = tzkt.get("/v1/head", nil, &head)
	return
}

// GetBlock -
func (tzkt *API) GetBlock(level uint64) (b Block, err error) {
	err = tzkt.get(fmt.Sprintf("/v1/blocks/%d", level), nil, &b)
	return
}

// GetBlocks -
func (tzkt *API) GetBlocks(filters map[string]string) (b []Block, err error) {
	err = tzkt.get("/v1/blocks", filters, &b)
	return
}

// GetEndorsements -
func (tzkt *API) GetEndorsements(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.get("/v1/operations/endorsements", filters, &operations)
	return
}

// GetBallots -
func (tzkt *API) GetBallots(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.get("/v1/operations/ballots", filters, &operations)
	return
}

// GetProposals -
func (tzkt *API) GetProposals(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.get("/v1/operations/proposals", filters, &operations)
	return
}

// GetActivations -
func (tzkt *API) GetActivations(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.get("/v1/operations/activations", filters, &operations)
	return
}

// GetDoubleBakings -
func (tzkt *API) GetDoubleBakings(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.get("/v1/operations/double_baking", filters, &operations)
	return
}

// GetDoubleEndorsings -
func (tzkt *API) GetDoubleEndorsings(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.get("/v1/operations/double_endorsing", filters, &operations)
	return
}

// GetNonceRevelations -
func (tzkt *API) GetNonceRevelations(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.get("/v1/operations/nonce_revelations", filters, &operations)
	return
}

// GetDelegations -
func (tzkt *API) GetDelegations(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.get("/v1/operations/delegations", filters, &operations)
	return
}

// GetOriginations -
func (tzkt *API) GetOriginations(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.get("/v1/operations/originations", filters, &operations)
	return
}

// GetTransactions -
func (tzkt *API) GetTransactions(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.get("/v1/operations/transactions", filters, &operations)
	return
}

// GetReveals -
func (tzkt *API) GetReveals(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.get("/v1/operations/reveals", filters, &operations)
	return
}
