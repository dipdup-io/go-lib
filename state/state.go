package state

import (
	"gorm.io/gorm"
)

// State -
type State struct {
	IndexName string `gorm:"primaryKey"`
	IndexType string
	Hash      string
	Level     uint64
}

// TableName -
func (State) TableName() string {
	return "dipdup_state"
}

// UpdateState -
func (s State) Update(db *gorm.DB) error {
	return db.Save(&s).Error
}

// Get -
func Get(db *gorm.DB, indexName string) (state State, err error) {
	err = db.Where("index_name = ?", indexName).First(&state).Error
	return
}
