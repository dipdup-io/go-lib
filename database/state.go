package database

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// State -
type State struct {
	//nolint
	tableName     struct{} `gorm:"-" bun:"-" pg:"dipdup_state" json:"-" comment:"Indexer state table"`
	bun.BaseModel `gorm:"-" pg:"-" bun:"dipdup_state" json:"-" comment:"Indexer state table"`

	IndexName string    `gorm:"primaryKey" pg:",pk" bun:",pk" json:"index_name" comment:"Index name"`
	IndexType string    `json:"index_type" comment:"Index type"`
	Hash      string    `json:"hash" comment:"Current hash"`
	Timestamp time.Time `json:"timestamp" comment:"Current timestamp"`
	Level     uint64    `json:"level" comment:"Index level"`
	UpdatedAt int       `gorm:"autoUpdateTime" comment:"Last updated timestamp"`
	CreatedAt int       `gorm:"autoCreateTime" comment:"Created timestamp"`
}

var _ bun.BeforeAppendModelHook = (*State)(nil)

func (s *State) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		s.UpdatedAt = int(time.Now().Unix())
		s.CreatedAt = s.UpdatedAt

	case *bun.UpdateQuery:
		s.UpdatedAt = int(time.Now().Unix())
	}
	return nil
}

// TableName -
func (State) TableName() string {
	return "dipdup_state"
}

// StateRepository -
type StateRepository interface {
	State(ctx context.Context, name string) (*State, error)
	UpdateState(sctx context.Context, tate *State) error
	CreateState(ctx context.Context, state *State) error
	DeleteState(ctx context.Context, state *State) error
}
