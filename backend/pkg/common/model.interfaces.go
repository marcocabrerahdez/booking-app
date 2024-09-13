package common

import (
	"backend/pkg/helpers"

	"github.com/google/uuid"
)

// Composite interface that all DTOs must implement
type DTO interface {
	Identifiable
	Entityable
	Validable
}

// Composite interface that all Entities must implement
type Entity interface {
	Identifiable
	Dtoable
}

// Identifiable interface is used to get ID of Entity or DTO
// Entities that embed base.Entity and DTOs that embed base.DTO have this interface implemented
type Identifiable interface {
	GetID() uuid.UUID
	SetID(id uuid.UUID)
}

// Dtoable interface is used to convert Entity to DTO
// This interface must be implemented by all DTOs
type Dtoable interface {
	ToDTO() DTO
}

// Validable interface is used to validate DTOs
// DTOs that embed base.DTO have this interface implemented
type Validable interface {
	Validate() []*helpers.ValidationErrors
}

// Entityable interface is used to convert DTO to Entity
// This interface must be implemented by all DTOs
type Entityable interface {
	ToEntity() Entity
}
