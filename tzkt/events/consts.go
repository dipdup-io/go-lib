package events

// Methods
const (
	MethodHead           = "SubscribeToHead"
	MethodBlocks         = "SubscribeToBlocks"
	MethodOperations     = "SubscribeToOperations"
	MethodBigMap         = "SubscribeToBigMaps"
	MethodAccounts       = "SubscribeToAccounts"
	MethodTokenTransfers = "SubscribeToTokenTransfers"
)

// Channels
const (
	ChannelHead       = "head"
	ChannelBlocks     = "blocks"
	ChannelOperations = "operations"
	ChannelBigMap     = "bigmaps"
	ChannelAccounts   = "accounts"
	ChannelTransfers  = "transfers"
)

// Big map tags
const (
	BigMapTagMetadata      = "metadata"
	BigMapTagTokenMetadata = "token_metadata"
)
