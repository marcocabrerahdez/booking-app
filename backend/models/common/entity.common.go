package common

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CommonEntity is a struct that contains common fields for all entities
type CommonEntity struct {
	ID        uuid.UUID `gorm:"primaryKey;type:char(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (e *CommonEntity) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

func (e CommonEntity) GetID() uuid.UUID {
	return e.ID
}

func (e *CommonEntity) SetID(id uuid.UUID) {
	e.ID = id
}
