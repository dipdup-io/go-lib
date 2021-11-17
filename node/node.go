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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// NodeRPC -
type NodeRPC struct {
	baseURL string

	timeout time.Duration
}

// NewNodeRPC -
func NewNodeRPC(baseURL string, opts ...NodeOption) *NodeRPC {
	node := &NodeRPC{
		baseURL: strings.TrimSuffix(baseURL, "/"),
	}

	for i := range opts {
		opts[i](node)
	}

	return node
}

// URL -
func (rpc *NodeRPC) URL() string {
	return rpc.baseURL
}

func (rpc *NodeRPC) parseResponse(resp *http.Response, response interface{}) error {
	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(ErrInvalidResponse, resp.Status)
		}
		return errors.Wrap(ErrInvalidResponse, string(data))
	}
	return json.NewDecoder(resp.Body).Decode(response)
}

func (rpc *NodeRPC) makeRequest(method, uri string, queryArgs url.Values, body interface{}) (*http.Response, error) {
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

	req, err := http.NewRequest(method, link.String(), bodyReader)
	if err != nil {
		return nil, errors.Errorf("makeGetRequest.NewRequest: %v", err)
	}
	client := http.Client{
		Timeout: rpc.timeout,
	}
	return client.Do(req)
}

//nolint
func (rpc *NodeRPC) get(uri string, queryArgs url.Values, response interface{}) error {
	resp, err := rpc.makeRequest(http.MethodGet, uri, queryArgs, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return rpc.parseResponse(resp, response)
}

//nolint
func (rpc *NodeRPC) post(uri string, queryArgs url.Values, body, response interface{}) error {
	resp, err := rpc.makeRequest(http.MethodPost, uri, queryArgs, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return rpc.parseResponse(resp, response)
}

// PendingOperations -
func (rpc *NodeRPC) PendingOperations() (res MempoolResponse, err error) {
	err = rpc.get("chains/main/mempool/pending_operations", nil, &res)
	return
}

// Constants -
func (rpc *NodeRPC) Constants() (constants Constants, err error) {
	err = rpc.get("chains/main/blocks/head/context/constants", nil, &constants)
	return
}

// ActiveDelegatesWithRolls -
func (rpc *NodeRPC) ActiveDelegatesWithRolls() (delegates []string, err error) {
	err = rpc.get("chains/main/blocks/head/context/raw/json/active_delegates_with_rolls", nil, &delegates)
	return
}

// Delegates -
func (rpc *NodeRPC) Delegates(active *bool) (delegates []string, err error) {
	queryArgs := make(url.Values)
	if active != nil && *active {
		queryArgs.Add("active", "true")
	}
	err = rpc.get("chains/main/blocks/head/context/delegates", queryArgs, &delegates)
	return
}

// StakingBalance -
func (rpc *NodeRPC) StakingBalance(address string) (balance string, err error) {
	uri := fmt.Sprintf("chains/main/blocks/head/context/delegates/%s/staking_balance", address)
	err = rpc.get(uri, nil, &balance)
	return
}

// InjectOperaiton -
func (rpc *NodeRPC) InjectOperaiton(request InjectOperationRequest) (hash string, err error) {
	queryArgs := make(url.Values)
	if request.Async {
		queryArgs.Add("async", "true")
	}
	if request.ChainID != "" {
		queryArgs.Add("chain", request.ChainID)
	}
	err = rpc.post("injection/operation", queryArgs, request.Operation, &hash)
	return
}

// Counter -
func (rpc *NodeRPC) Counter(contract string, block string) (counter string, err error) {
	uri := fmt.Sprintf("chains/main/blocks/%s/context/contracts/%s/counter", block, contract)
	err = rpc.get(uri, nil, &counter)
	return
}

// Header -
func (rpc *NodeRPC) Header(block string) (head Header, err error) {
	err = rpc.get(fmt.Sprintf("chains/main/blocks/%s/header", block), nil, &head)
	return
}

// HeadMetadata -
func (rpc *NodeRPC) HeadMetadata(block string) (head HeadMetadata, err error) {
	err = rpc.get(fmt.Sprintf("chains/main/blocks/%s/metadata", block), nil, &head)
	return
}

// Operations -
func (rpc *NodeRPC) Operations(block string) (operations [][]Operation, err error) {
	err = rpc.get(fmt.Sprintf("chains/main/blocks/%s/operations", block), nil, &operations)
	return
}

// ManagerOperations -
func (rpc *NodeRPC) ManagerOperations(block string) (operations []Operation, err error) {
	err = rpc.get(fmt.Sprintf("chains/main/blocks/%s/operations/3", block), nil, &operations)
	return
}
