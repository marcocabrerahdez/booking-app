package models

import (
	"backend/database"
	"backend/pkg/common"
)

func init() {
	database.RegisterModel(&database.MigrationTask{
		Model:           &Business{},
		DropOnFlush:     true,
		TruncateOnFlush: true,
	})
}

type Business struct {
	common.CommonEntity `gorm:"embedded"`
	Name                string `gorm:"type:varchar(255);not null"`
	Type                string `gorm:"type:varchar(255);not null"`
	Location            string `gorm:"type:varchar(255);not null"`
	OwnerID             string `gorm:"type:varchar(255);not null"`
	Capacity            int    `gorm:"type:int;not null"`

	// Relationships
	Owner        User          `gorm:"foreignKey:OwnerID"`
	Reservations []Reservation `gorm:"foreignKey:BusinessID"`
}

type BusinessDTO struct {
	common.CommonDTO `json:",inline,omitempty" tstype:",extends,required"`
	Name             string `json:"name" tstype:"string,required"`
	Type             string `json:"type" tstype:"string,required"`
	Location         string `json:"location" tstype:"string,required"`
	OwnerID          string `json:"owner_id" tstype:"string,required"`
	Capacity         int    `json:"capacity" tstype:"number,required"`
}

func (b Business) ToDTO() common.DTO {
	dto := &BusinessDTO{
		CommonDTO: common.CommonDTO{
			ID:        b.ID,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		},
		Name:     b.Name,
		Type:     b.Type,
		Location: b.Location,
		OwnerID:  b.OwnerID,
		Capacity: b.Capacity,
	}
	return dto
}

func (b BusinessDTO) ToEntity() common.Entity {
	entity := &Business{
		CommonEntity: common.CommonEntity{
			ID:        b.ID,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		},
		Name:     b.Name,
		Type:     b.Type,
		Location: b.Location,
		OwnerID:  b.OwnerID,
		Capacity: b.Capacity,
	}
	return entity
}
