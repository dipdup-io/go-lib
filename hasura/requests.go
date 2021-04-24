package hasura

type request struct {
	Type string      `json:"type"`
	Args interface{} `json:"args"`
}
