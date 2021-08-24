package events

// Methods
const (
	MethodHead       = "SubscribeToHead"
	MethodBlocks     = "SubscribeToBlocks"
	MethodOperations = "SubscribeToOperations"
	MethodBigMap     = "SubscribeToBigMaps"
	MethodAccounts   = "SubscribeToAccounts"
)

// Channels
const (
	ChannelHead       = "head"
	ChannelBlocks     = "blocks"
	ChannelOperations = "operations"
	ChannelBigMap     = "bigmaps"
	ChannelAccounts   = "accounts"
)

// operation kinds
const (
	KindTransaction     = "transaction"
	KindOrigination     = "origination"
	KindDelegation      = "delegation"
	KindEndorsement     = "endorsement"
	KindBallot          = "ballot"
	KindProposal        = "proposal"
	KindActivation      = "activation"
	KindDoubleBaking    = "double_baking"
	KindDoubleEndorsing = "double_endorsing"
	KindNonceRevelation = "nonce_revelation"
	KindReveal          = "reveal"
)

// Base URL
const (
	BaseURL = "https://api.tzkt.io/v1/events"
)

// Big map tags
const (
	BigMapTagMetadata      = "metadata"
	BigMapTagTokenMetadata = "token_metadata"
)
