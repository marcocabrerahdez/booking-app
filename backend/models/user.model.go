package models

import (
	"backend/database"
	"backend/pkg/common"
)

func init() {
	database.RegisterModel(&database.MigrationTask{
		Model:           &User{},
		DropOnFlush:     true,
		TruncateOnFlush: true,
	})
}

type User struct {
	common.CommonEntity `gorm:"embedded"`
	Name                string `gorm:"type:varchar(255);not null"`
	Email               string `gorm:"type:varchar(255);not null;unique"`
	Password            string `gorm:"type:varchar(255);not null"`
	Role                string `gorm:"type:varchar(255);not null"`
}

type UserDTO struct {
	common.CommonDTO `json:",inline,omitempty" tstype:",extends,required"`
	Name             string `json:"name" tstype:"string,required"`
	Email            string `json:"email" tstype:"string,required"`
	Password         string `json:"password" tstype:"string,required"`
	Role             string `json:"role" tstype:"string,required"`
}

func (u User) ToDTO() common.DTO {
	dto := &UserDTO{
		CommonDTO: common.CommonDTO{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
	}
	return dto
}

func (u UserDTO) ToEntity() common.Entity {
	entity := &User{
		CommonEntity: common.CommonEntity{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
	}
	return entity
}
