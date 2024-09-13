package common

import (
	"time"

	"backend/pkg/helpers"

	"github.com/google/uuid"
)

// CommonDTO is a struct that contains common fields for all DTOs
type CommonDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
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
