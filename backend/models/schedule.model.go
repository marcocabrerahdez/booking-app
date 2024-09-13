package models

import (
	"backend/database"
	"backend/pkg/common"
	"time"

	"github.com/google/uuid"
)

func init() {
	database.RegisterModel(&database.MigrationTask{
		Model:           &Schedule{},
		DropOnFlush:     true,
		TruncateOnFlush: true,
	})
}

type Schedule struct {
	common.CommonEntity `gorm:"embedded"`
	BusinessID          uuid.UUID `gorm:"type:uuid;not null"`
	DayOfWeek           int       `gorm:"type:int;not null"`
	StartTime           time.Time `gorm:"type:timestamp;not null"`
	EndTime             time.Time `gorm:"type:timestamp;not null"`

	// Relationships
	Business Business `gorm:"foreignKey:BusinessID"`
}

type ScheduleDTO struct {
	common.CommonDTO `json:",inline,omitempty" tstype:",extends,required"`
	BusinessID       uuid.UUID `json:"business_id" tstype:"string,required"`
	DayOfWeek        int       `json:"day_of_week" tstype:"number,required"`
	StartTime        time.Time `json:"start_time" tstype:"string,required"`
	EndTime          time.Time `json:"end_time" tstype:"string,required"`
}

func (s Schedule) ToDTO() common.DTO {
	dto := &ScheduleDTO{
		CommonDTO: common.CommonDTO{
			ID:        s.ID,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		},
		BusinessID: s.BusinessID,
		DayOfWeek:  s.DayOfWeek,
		StartTime:  s.StartTime,
		EndTime:    s.EndTime,
	}
	return dto
}

func (s ScheduleDTO) ToEntity() common.Entity {
	entity := &Schedule{
		CommonEntity: common.CommonEntity{
			ID:        s.ID,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		},
		BusinessID: s.BusinessID,
		DayOfWeek:  s.DayOfWeek,
		StartTime:  s.StartTime,
		EndTime:    s.EndTime,
	}
	return entity
}
