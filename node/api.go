package node

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// API -
type API interface {
	BlockAPI
	ChainAPI
	ContextAPI
	ConfigAPI
	GeneralAPI
	ProtocolsAPI
	NetworkAPI
	InjectAPI
}

// RPC -
type RPC struct {
	*BlockRPC
	*Chain
	*Context
	*Config
	*General
	*Protocols
	*Network
	*Inject
}

// NewRPC -
func NewRPC(baseURL, chainID string) *RPC {
	return &RPC{
		BlockRPC:  NewBlockRPC(baseURL, chainID),
		Chain:     NewChain(baseURL, chainID),
		Context:   NewContext(baseURL, chainID),
		Config:    NewConfig(baseURL),
		General:   NewGeneral(baseURL),
		Protocols: NewProtocols(baseURL),
		Network:   NewNetwork(baseURL),
		Inject:    NewInject(baseURL),
	}
}

// NewMainRPC -
func NewMainRPC(baseURL string) *RPC {
	return &RPC{
		BlockRPC:  NewMainBlockRPC(baseURL),
		Chain:     NewMainChain(baseURL),
		Context:   NewMainContext(baseURL),
		Config:    NewConfig(baseURL),
		General:   NewGeneral(baseURL),
		Protocols: NewProtocols(baseURL),
		Network:   NewNetwork(baseURL),
		Inject:    NewInject(baseURL),
	}
}
