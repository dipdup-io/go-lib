package data

// ContractJSONSchema -
type ContractJSONSchema struct {
	Storage     JSONSchema             `json:"storageSchema"`
	Entrypoints []EntrypointJSONSchema `json:"entrypoints"`
	BigMaps     []BigMapJSONSchema     `json:"bigMaps"`
}

// EntrypointJSONSchema -
type EntrypointJSONSchema struct {
	Name      string     `json:"name"`
	Parameter JSONSchema `json:"parameterSchema"`
}
