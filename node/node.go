package node

import (
	"fmt"
	"net/http"
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
		return errors.Errorf("Invalid response %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(response)
}

func (rpc *NodeRPC) makeGetRequest(uri string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", rpc.baseURL, uri)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Errorf("makeGetRequest.NewRequest: %v", err)
	}
	client := http.Client{
		Timeout: rpc.timeout,
	}
	return client.Do(req)
}

//nolint
func (rpc *NodeRPC) get(uri string, response interface{}) error {
	resp, err := rpc.makeGetRequest(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return rpc.parseResponse(resp, response)
}

// PendingOperations -
func (rpc *NodeRPC) PendingOperations() (res MempoolResponse, err error) {
	err = rpc.get("chains/main/mempool/pending_operations", &res)
	return
}

// Constants -
func (rpc *NodeRPC) Constants() (constants Constants, err error) {
	err = rpc.get("chains/main/blocks/head/context/constants", &constants)
	return
}

// Head -
func (rpc *NodeRPC) Head() (head Header, err error) {
	err = rpc.get("chains/main/blocks/head/header", &head)
	return
}
