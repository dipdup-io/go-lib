package hasura

type replaceMetadataResponse struct {
	Message string `json:"message"`
}

// ExportMetadataResponse -
type ExportMetadataResponse struct {
	Tables []Table `json:"tables"`
}
