package node

// Version -
type Version struct {
	Version struct {
		Major          int    `json:"major"`
		Minor          int    `json:"minor"`
		AdditionalInfo string `json:"additional_info"`
	} `json:"version"`
	NetworkVersion struct {
		ChainName            string `json:"chain_name"`
		DistributedDbVersion int    `json:"distributed_db_version"`
		P2PVersion           int    `json:"p2p_version"`
	} `json:"network_version"`
	CommitInfo struct {
		CommitHash string `json:"commit_hash"`
		CommitDate string `json:"commit_date"`
	} `json:"commit_info"`
}

// StatsGC -
type StatsGC struct {
	MinorWords             int64 `json:"minor_words"`
	PromotedWords          int64 `json:"promoted_words"`
	MajorWords             int64 `json:"major_words"`
	MinorCollections       int64 `json:"minor_collections"`
	MajorCollections       int64 `json:"major_collections"`
	ForcedMajorCollections int64 `json:"forced_major_collections"`
	HeapWords              int64 `json:"heap_words"`
	HeapChunks             int64 `json:"heap_chunks"`
	LiveWords              int64 `json:"live_words"`
	LiveBlocks             int64 `json:"live_blocks"`
	FreeWords              int64 `json:"free_words"`
	FreeBlocks             int64 `json:"free_blocks"`
	LargestFree            int64 `json:"largest_free"`
	Fragments              int64 `json:"fragments"`
	Compactions            int64 `json:"compactions"`
	TopHeapWords           int64 `json:"top_heap_words"`
	StackSize              int64 `json:"stack_size"`
}

// StatsMemory -
type StatsMemory struct {
	PageSize int64  `json:"page_size"`
	Size     string `json:"size"`
	Resident string `json:"resident"`
	Shared   string `json:"shared"`
	Text     string `json:"text"`
	Lib      string `json:"lib"`
	Data     string `json:"data"`
	Dt       string `json:"dt"`
}
