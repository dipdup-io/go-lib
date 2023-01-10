package data

import (
	stdJSON "encoding/json"
	"time"
)

// Token -
type Token struct {
	ID             uint64             `json:"id"`
	Contract       Address            `json:"contract"`
	TokenID        string             `json:"tokenId"`
	Standard       string             `json:"standard"`
	Metadata       stdJSON.RawMessage `json:"metadata,omitempty"`
	FirstLevel     uint64             `json:"firstLevel"`
	FirstTime      time.Time          `json:"firstTime"`
	LastLevel      uint64             `json:"lastLevel"`
	LastTime       time.Time          `json:"lastTime"`
	TransfersCount uint64             `json:"transfersCount"`
	BalancesCount  uint64             `json:"balancesCount"`
	HoldersCount   uint64             `json:"holdersCount"`
	TotalMinted    string             `json:"totalMinted"`
	TotalBurned    string             `json:"totalBurned"`
	TotalSupply    string             `json:"totalSupply"`
}

// TokenBalance -
type TokenBalance struct {
	ID             uint64    `json:"id"`
	Account        *Address  `json:"account,omitempty"`
	Token          *Token    `json:"token,omitempty"`
	Balance        string    `json:"balance,omitempty"`
	TransfersCount uint64    `json:"transfersCount"`
	FirstLevel     uint64    `json:"firstLevel"`
	FirstTime      time.Time `json:"firstTime"`
	LastLevel      uint64    `json:"lastLevel"`
	LastTime       time.Time `json:"lastTime"`
}
