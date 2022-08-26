package data

import "time"

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

// Contract -
type Contract struct {
	ID      int      `json:"id"`
	Type    string   `json:"type"`
	Address string   `json:"address"`
	Kind    string   `json:"kind"`
	Tzips   []string `json:"tzips"`
	Balance int      `json:"balance"`
	Creator struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"creator"`
	NumContracts        int       `json:"numContracts"`
	ActiveTokensCount   int       `json:"activeTokensCount"`
	TokenBalancesCount  int       `json:"tokenBalancesCount"`
	TokenTransfersCount int       `json:"tokenTransfersCount"`
	NumDelegations      int       `json:"numDelegations"`
	NumOriginations     int       `json:"numOriginations"`
	NumTransactions     int       `json:"numTransactions"`
	NumReveals          int       `json:"numReveals"`
	NumMigrations       int       `json:"numMigrations"`
	TransferTicketCount int       `json:"transferTicketCount"`
	FirstActivity       int       `json:"firstActivity"`
	FirstActivityTime   time.Time `json:"firstActivityTime"`
	LastActivity        int       `json:"lastActivity"`
	LastActivityTime    time.Time `json:"lastActivityTime"`
	TypeHash            int       `json:"typeHash"`
	CodeHash            int       `json:"codeHash"`
}
