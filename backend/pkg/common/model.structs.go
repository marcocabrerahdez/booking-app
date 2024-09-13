package common

import (
	"backend/pkg/helpers"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Embedded struct that contains common fields for all entities
// Must be embedded in all entities
// Must use gorm:"embedded" tag
// For example: type User struct { common.CommonEntity `gorm:"embedded"` }
type CommonEntity struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Embedded struct that contains common fields for all DTOs
// Must be embedded in all DTOs
// Must use json:",inline,omitempty" tag
// For example: type UserDTO struct { common.CommonDTO `json:",inline,omitempty"` }
type CommonDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

func (e CommonDTO) GetID() uuid.UUID {
	return e.ID
}

func (e *CommonDTO) SetID(id uuid.UUID) {
	e.ID = id
}

func (e CommonDTO) Validate() []*helpers.ValidationErrors {
	return helpers.ValidateStruct(e)
}
