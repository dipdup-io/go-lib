package node

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// BlockAPI -
type BlockAPI interface {
	Blocks(ctx context.Context, args BlocksArgs) ([][]string, error)
	Block(ctx context.Context, blockID string) (Block, error)
	Head(ctx context.Context) (Block, error)
	Header(ctx context.Context, blockID string) (Header, error)
	HeaderRaw(ctx context.Context, blockID string) (string, error)
	HeaderShell(ctx context.Context, blockID string) (HeaderShell, error)
	Metadata(ctx context.Context, blockID string) (BlockMetadata, error)
	MetadataHash(ctx context.Context, blockID string) (string, error)
	Hash(ctx context.Context, blockID string) (string, error)
	ProtocolData(ctx context.Context, blockID string) (ProtocolData, error)
	ProtocolDataRaw(ctx context.Context, blockID string) (string, error)
	OperationHashes(ctx context.Context, blockID string) ([][]string, error)
	OperationMetadataHash(ctx context.Context, blockID string) (string, error)
	OperationMetadataHashes(ctx context.Context, blockID string) ([][]string, error)
	Operations(ctx context.Context, blockID string) ([][]OperationGroup, error)
	OperationsOffset(ctx context.Context, blockID string, listOffset int) ([]OperationGroup, error)
	Operation(ctx context.Context, blockID string, listOffset, operationOffset int) (OperationGroup, error)
	BlockProtocols(ctx context.Context, blockID string) (BlockProtocols, error)
	VotesBallotList(ctx context.Context, blockID string) ([]BlockBallot, error)
	VotesBallots(ctx context.Context, blockID string) (BlockBallots, error)
	VotesCurrentPeriod(ctx context.Context, blockID string) (VotingPeriod, error)
	VotesCurrentProposal(ctx context.Context, blockID string) (string, error)
	VotesQuorum(ctx context.Context, blockID string) (int, error)
	VotesListing(ctx context.Context, blockID string) ([]Rolls, error)
	VotesProposals(ctx context.Context, blockID string) ([]string, error)
	VotesSuccessorPeriod(ctx context.Context, blockID string) (VotingPeriod, error)
	VotesTotalVotingPower(ctx context.Context, blockID string) (int, error)
}

// BlockRPC -
type BlockRPC struct {
	baseURL string
	chainID string
	client  *client
}

// NewChain -
func NewBlockRPC(baseURL, chainID string) *BlockRPC {
	return &BlockRPC{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		chainID: chainID,
		client:  newClient(),
	}
}

// NewMainBlockRPC -
func NewMainBlockRPC(baseURL string) *BlockRPC {
	return NewBlockRPC(baseURL, "main")
}

// Blocks -
func (api *BlockRPC) Blocks(ctx context.Context, args BlocksArgs) ([][]string, error) {
	queryArgs := make(url.Values)
	if args.HeadHash != "" {
		queryArgs.Add("head", args.HeadHash)
	}
	if args.Length > 0 {
		queryArgs.Add("length", fmt.Sprintf("%d", args.Length))
	}

	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/", api.chainID), queryArgs)
	if err != nil {
		return nil, err
	}
	var blocks [][]string
	err = req.doWithJSONResponse(ctx, api.client, &blocks)
	return blocks, err
}

// Block -
func (api *BlockRPC) Block(ctx context.Context, blockID string) (Block, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s", api.chainID, blockID), nil)
	if err != nil {
		return Block{}, err
	}
	var block Block
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// Head -
func (api *BlockRPC) Head(ctx context.Context) (Block, error) {
	return api.Block(ctx, "head")
}

// Header -
func (api *BlockRPC) Header(ctx context.Context, blockID string) (Header, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/header", api.chainID, blockID), nil)
	if err != nil {
		return Header{}, err
	}
	var block Header
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// HeaderRaw -
func (api *BlockRPC) HeaderRaw(ctx context.Context, blockID string) (string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/header/raw", api.chainID, blockID), nil)
	if err != nil {
		return "", err
	}
	var block string
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// HeaderShell -
func (api *BlockRPC) HeaderShell(ctx context.Context, blockID string) (HeaderShell, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/header/shell", api.chainID, blockID), nil)
	if err != nil {
		return HeaderShell{}, err
	}
	var block HeaderShell
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// Metadata -
func (api *BlockRPC) Metadata(ctx context.Context, blockID string) (BlockMetadata, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/metadata", api.chainID, blockID), nil)
	if err != nil {
		return BlockMetadata{}, err
	}
	var block BlockMetadata
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// MetadataHash -
func (api *BlockRPC) MetadataHash(ctx context.Context, blockID string) (string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/metadata_hash", api.chainID, blockID), nil)
	if err != nil {
		return "", err
	}
	var block string
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// Hash -
func (api *BlockRPC) Hash(ctx context.Context, blockID string) (string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/hash", api.chainID, blockID), nil)
	if err != nil {
		return "", err
	}
	var block string
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// ProtocolData -
func (api *BlockRPC) ProtocolData(ctx context.Context, blockID string) (ProtocolData, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/header/protocol_data", api.chainID, blockID), nil)
	if err != nil {
		return ProtocolData{}, err
	}
	var block ProtocolData
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// ProtocolDataRaw -
func (api *BlockRPC) ProtocolDataRaw(ctx context.Context, blockID string) (string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/header/protocol_data/raw", api.chainID, blockID), nil)
	if err != nil {
		return "", err
	}
	var block string
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// OperationHashes -
func (api *BlockRPC) OperationHashes(ctx context.Context, blockID string) ([][]string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/header/protocol_data/raw", api.chainID, blockID), nil)
	if err != nil {
		return nil, err
	}
	var block [][]string
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// OperationMetadataHash -
func (api *BlockRPC) OperationMetadataHash(ctx context.Context, blockID string) (string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/operations_metadata_hash", api.chainID, blockID), nil)
	if err != nil {
		return "", err
	}
	var block string
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// OperationMetadataHashes -
func (api *BlockRPC) OperationMetadataHashes(ctx context.Context, blockID string) ([][]string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/operation_metadata_hashes", api.chainID, blockID), nil)
	if err != nil {
		return nil, err
	}
	var block [][]string
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// Operations -
func (api *BlockRPC) Operations(ctx context.Context, blockID string) ([][]OperationGroup, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/operations", api.chainID, blockID), nil)
	if err != nil {
		return nil, err
	}
	var block [][]OperationGroup
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// OperationsOffset -
func (api *BlockRPC) OperationsOffset(ctx context.Context, blockID string, listOffset int) ([]OperationGroup, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/operations/%d", api.chainID, blockID, listOffset), nil)
	if err != nil {
		return nil, err
	}
	var block []OperationGroup
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// Operation -
func (api *BlockRPC) Operation(ctx context.Context, blockID string, listOffset, operationOffset int) (OperationGroup, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/operations/%d/%d", api.chainID, blockID, listOffset, operationOffset), nil)
	if err != nil {
		return OperationGroup{}, err
	}
	var block OperationGroup
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// BlockProtocols -
func (api *BlockRPC) BlockProtocols(ctx context.Context, blockID string) (BlockProtocols, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/protocols", api.chainID, blockID), nil)
	if err != nil {
		return BlockProtocols{}, err
	}
	var block BlockProtocols
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

func (api *BlockRPC) VotesBallotList(ctx context.Context, blockID string) ([]BlockBallot, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/votes/ballot_list", api.chainID, blockID), nil)
	if err != nil {
		return nil, err
	}
	var block []BlockBallot
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// VotesBallots -
func (api *BlockRPC) VotesBallots(ctx context.Context, blockID string) (BlockBallots, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/votes/ballots", api.chainID, blockID), nil)
	if err != nil {
		return BlockBallots{}, err
	}
	var block BlockBallots
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// VotesCurrentPeriod -
func (api *BlockRPC) VotesCurrentPeriod(ctx context.Context, blockID string) (VotingPeriod, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/votes/current_period", api.chainID, blockID), nil)
	if err != nil {
		return VotingPeriod{}, err
	}
	var block VotingPeriod
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// VotesCurrentProposal -
func (api *BlockRPC) VotesCurrentProposal(ctx context.Context, blockID string) (string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/votes/current_proposal", api.chainID, blockID), nil)
	if err != nil {
		return "", err
	}
	var block string
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// VotesQuorum -
func (api *BlockRPC) VotesQuorum(ctx context.Context, blockID string) (int, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/votes/current_quorum", api.chainID, blockID), nil)
	if err != nil {
		return 0, err
	}
	var block int
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// VotesListing
func (api *BlockRPC) VotesListing(ctx context.Context, blockID string) ([]Rolls, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/votes/listings", api.chainID, blockID), nil)
	if err != nil {
		return nil, err
	}
	var block []Rolls
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// VotesProposals -
func (api *BlockRPC) VotesProposals(ctx context.Context, blockID string) ([]string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/votes/proposals", api.chainID, blockID), nil)
	if err != nil {
		return nil, err
	}
	var block []string
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// VotesSuccessorPeriod -
func (api *BlockRPC) VotesSuccessorPeriod(ctx context.Context, blockID string) (VotingPeriod, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/votes/successor_period", api.chainID, blockID), nil)
	if err != nil {
		return VotingPeriod{}, err
	}
	var block VotingPeriod
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}

// VotesTotalVotingPower -
func (api *BlockRPC) VotesTotalVotingPower(ctx context.Context, blockID string) (int, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/votes/total_voting_power", api.chainID, blockID), nil)
	if err != nil {
		return 0, err
	}
	var block int
	err = req.doWithJSONResponse(ctx, api.client, &block)
	return block, err
}
