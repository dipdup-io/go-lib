package database

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// State -
type State struct {
	//nolint
	tableName     struct{} `bun:"-"            comment:"Indexer state table" gorm:"-" json:"-" pg:"dipdup_state"`
	bun.BaseModel `bun:"dipdup_state" comment:"Indexer state table" gorm:"-" json:"-" pg:"-"`

	IndexName string    `bun:",pk"                        comment:"Index name"  gorm:"primaryKey" json:"index_name" pg:",pk"`
	IndexType string    `comment:"Index type"             json:"index_type"`
	Hash      string    `comment:"Current hash"           json:"hash"`
	Timestamp time.Time `comment:"Current timestamp"      json:"timestamp"`
	Level     uint64    `comment:"Index level"            json:"level"`
	UpdatedAt int       `comment:"Last updated timestamp" gorm:"autoUpdateTime"`
	CreatedAt int       `comment:"Created timestamp"      gorm:"autoCreateTime"`
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
