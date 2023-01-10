package events

// Methods
const (
	MethodHead           = "SubscribeToHead"
	MethodBlocks         = "SubscribeToBlocks"
	MethodOperations     = "SubscribeToOperations"
	MethodBigMap         = "SubscribeToBigMaps"
	MethodAccounts       = "SubscribeToAccounts"
	MethodTokenTransfers = "SubscribeToTokenTransfers"
	MethodTokenBalances  = "SubscribeToTokenBalances"
	MethodCycles         = "SubscribeToCycles"
)

// Channels
const (
	ChannelHead          = "head"
	ChannelBlocks        = "blocks"
	ChannelOperations    = "operations"
	ChannelBigMap        = "bigmaps"
	ChannelAccounts      = "accounts"
	ChannelTransfers     = "transfers"
	ChannelCycles        = "cycles"
	ChannelTokenBalances = "token_balances"
)

// Big map tags
const (
	BigMapTagMetadata      = "metadata"
	BigMapTagTokenMetadata = "token_metadata"
)
