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

// BeforeInsert -
func (s *State) BeforeInsert(ctx context.Context) (context.Context, error) {
	s.UpdatedAt = int(time.Now().Unix())
	s.CreatedAt = s.UpdatedAt
	return ctx, nil
}

// BeforeUpdate -
func (s *State) BeforeUpdate(ctx context.Context) (context.Context, error) {
	s.UpdatedAt = int(time.Now().Unix())
	return ctx, nil
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
