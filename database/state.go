package database

import (
	"context"
	"time"
)

// State -
type State struct {
	//nolint
	tableName struct{} `gorm:"-" pg:"dipdup_state" json:"-"`

	IndexName string    `gorm:"primaryKey" pg:",pk" json:"index_name"`
	IndexType string    `json:"index_type"`
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
	Level     uint64    `json:"level"`
	UpdatedAt int       `gorm:"autoUpdateTime"`
}

// BeforeInsert -
func (s *State) BeforeInsert(ctx context.Context) (context.Context, error) {
	s.UpdatedAt = int(time.Now().Unix())
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
	State(name string) (State, error)
	UpdateState(state State) error
	CreateState(state State) error
	DeleteState(state State) error
}
