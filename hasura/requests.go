package hasura

type request struct {
	Type string      `json:"type"`
	Args interface{} `json:"args"`
}

// Permission -
type Permission struct {
	Columns   string      `json:"columns"`
	Limit     uint64      `json:"limit"`
	AllowAggs bool        `json:"allow_aggregations"`
	Filter    interface{} `json:"filter,omitempty"`
}
