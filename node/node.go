package node

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// NodeRPC -
type NodeRPC struct {
	baseURL string
}

// NewNodeRPC -
func NewNodeRPC(baseURL string) *NodeRPC {
	return &NodeRPC{
		baseURL: strings.TrimSuffix(baseURL, "/"),
	}
}

// URL -
func (rpc *NodeRPC) URL() string {
	return rpc.baseURL
}

func (rpc *NodeRPC) parseResponse(resp *http.Response, response interface{}) error {
	if resp.StatusCode != http.StatusOK {
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
	return json.NewDecoder(resp.Body).Decode(response)
}

func (rpc *NodeRPC) makeRequest(method, uri string, queryArgs url.Values, body interface{}, opts RequestOpts) (*http.Response, error) {
	link, err := url.Parse(rpc.baseURL)
	if err != nil {
		return nil, err
	}
	if uri != "" {
		link.Path = path.Join(link.Path, uri)
	}
	if len(queryArgs) > 0 {
		link.RawQuery = queryArgs.Encode()
	}

	var bodyReader io.ReadWriter
	if body != nil {
		bodyReader = new(bytes.Buffer)
		if err := json.NewEncoder(bodyReader).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(opts.ctx, method, link.String(), bodyReader)
	if err != nil {
		return nil, errors.Errorf("makeGetRequest.NewRequest: %v", err)
	}
	return http.DefaultClient.Do(req)
}

//nolint
func (rpc *NodeRPC) get(uri string, queryArgs url.Values, opts RequestOpts, response interface{}) error {
	resp, err := rpc.makeRequest(http.MethodGet, uri, queryArgs, nil, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return rpc.parseResponse(resp, response)
}

//nolint
func (rpc *NodeRPC) post(uri string, queryArgs url.Values, body interface{}, opts RequestOpts, response interface{}) error {
	resp, err := rpc.makeRequest(http.MethodPost, uri, queryArgs, body, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return rpc.parseResponse(resp, response)
}

// PendingOperations -
func (rpc *NodeRPC) PendingOperations(opts ...RequestOption) (res MempoolResponse, err error) {
	options := newRequestOpts(opts...)
	err = rpc.get("chains/main/mempool/pending_operations", nil, options, &res)
	return
}

// Constants -
func (rpc *NodeRPC) Constants(opts ...RequestOption) (constants Constants, err error) {
	options := newRequestOpts(opts...)
	err = rpc.get("chains/main/blocks/head/context/constants", nil, options, &constants)
	return
}

// ActiveDelegatesWithRolls -
func (rpc *NodeRPC) ActiveDelegatesWithRolls(opts ...RequestOption) (delegates []string, err error) {
	options := newRequestOpts(opts...)
	err = rpc.get("chains/main/blocks/head/context/raw/json/active_delegates_with_rolls", nil, options, &delegates)
	return
}

// Delegates -
func (rpc *NodeRPC) Delegates(active *bool, opts ...RequestOption) (delegates []string, err error) {
	queryArgs := make(url.Values)
	if active != nil && *active {
		queryArgs.Add("active", "true")
	}
	options := newRequestOpts(opts...)
	err = rpc.get("chains/main/blocks/head/context/delegates", queryArgs, options, &delegates)
	return
}

// StakingBalance -
func (rpc *NodeRPC) StakingBalance(address string, opts ...RequestOption) (balance string, err error) {
	options := newRequestOpts(opts...)
	uri := fmt.Sprintf("chains/main/blocks/head/context/delegates/%s/staking_balance", address)
	err = rpc.get(uri, nil, options, &balance)
	return
}

// InjectOperaiton -
func (rpc *NodeRPC) InjectOperaiton(request InjectOperationRequest, opts ...RequestOption) (hash string, err error) {
	queryArgs := make(url.Values)
	if request.Async {
		queryArgs.Add("async", "true")
	}
	if request.ChainID != "" {
		queryArgs.Add("chain", request.ChainID)
	}
	options := newRequestOpts(opts...)
	err = rpc.post("injection/operation", queryArgs, request.Operation, options, &hash)
	return
}

// Counter -
func (rpc *NodeRPC) Counter(contract string, block string, opts ...RequestOption) (counter string, err error) {
	options := newRequestOpts(opts...)
	uri := fmt.Sprintf("chains/main/blocks/%s/context/contracts/%s/counter", block, contract)
	err = rpc.get(uri, nil, options, &counter)
	return
}

// Header -
func (rpc *NodeRPC) Header(block string, opts ...RequestOption) (head Header, err error) {
	options := newRequestOpts(opts...)
	err = rpc.get(fmt.Sprintf("chains/main/blocks/%s/header", block), nil, options, &head)
	return
}

// HeadMetadata -
func (rpc *NodeRPC) HeadMetadata(block string, opts ...RequestOption) (head HeadMetadata, err error) {
	options := newRequestOpts(opts...)
	err = rpc.get(fmt.Sprintf("chains/main/blocks/%s/metadata", block), nil, options, &head)
	return
}

// Operations -
func (rpc *NodeRPC) Operations(block string, opts ...RequestOption) (operations [][]Operation, err error) {
	options := newRequestOpts(opts...)
	err = rpc.get(fmt.Sprintf("chains/main/blocks/%s/operations", block), nil, options, &operations)
	return
}

// ManagerOperations -
func (rpc *NodeRPC) ManagerOperations(block string, opts ...RequestOption) (operations []Operation, err error) {
	options := newRequestOpts(opts...)
	err = rpc.get(fmt.Sprintf("chains/main/blocks/%s/operations/3", block), nil, options, &operations)
	return
}

// ContractStorage -
func (rpc *NodeRPC) ContractStorage(block, contract string, output interface{}, opts ...RequestOption) (err error) {
	options := newRequestOpts(opts...)
	err = rpc.get(fmt.Sprintf("chains/main/blocks/%s/context/contracts/%s/storage", block, contract), nil, options, &output)
	return
}
