package models

import (
	"backend/database"
	"backend/pkg/common"

	"time"

	"github.com/google/uuid"
)

func init() {
	database.RegisterModel(&database.MigrationTask{
		Model:           &Reservation{},
		DropOnFlush:     true,
		TruncateOnFlush: true,
	})
}

type Reservation struct {
	common.CommonEntity `gorm:"embedded"`
	UserID              uuid.UUID `gorm:"type:uuid;not null"`
	BusinessID          uuid.UUID `gorm:"type:uuid;not null"`
	Date                time.Time `gorm:"type:timestamp;not null"`
	NumberOfPeople      int       `gorm:"type:int;not null"`
	Status              string    `gorm:"type:varchar(255);not null"`

	// Relationships
	User     User     `gorm:"foreignKey:UserID"`
	Business Business `gorm:"foreignKey:BusinessID"`
}

type ReservationDTO struct {
	common.CommonDTO `json:",inline,omitempty" tstype:",extends,required"`
	UserID           uuid.UUID `json:"user_id" tstype:"string,required"`
	BusinessID       uuid.UUID `json:"business_id" tstype:"string,required"`
	Date             time.Time `json:"date" tstype:"string,required"`
	NumberOfPeople   int       `json:"number_of_people" tstype:"number,required"`
	Status           string    `json:"status" tstype:"string,required"`
}

func (r Reservation) ToDTO() common.DTO {
	dto := &ReservationDTO{
		CommonDTO: common.CommonDTO{
			ID:        r.ID,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
		UserID:         r.UserID,
		BusinessID:     r.BusinessID,
		Date:           r.Date,
		NumberOfPeople: r.NumberOfPeople,
		Status:         r.Status,
	}
	return dto
}

func (r ReservationDTO) ToEntity() common.Entity {
	entity := &Reservation{
		CommonEntity: common.CommonEntity{
			ID:        r.ID,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
		UserID:         r.UserID,
		BusinessID:     r.BusinessID,
		Date:           r.Date,
		NumberOfPeople: r.NumberOfPeople,
		Status:         r.Status,
	}
	return entity
}
