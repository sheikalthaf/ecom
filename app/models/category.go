package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;column:Id" json:"Id"`
	Name        string    `gorm:"type:varchar(200);column:Name" json:"Name"`
	Description string    `gorm:"column:Description" json:"Description"`
	IsDeleted   bool      `gorm:"type:bool;column:IsDeleted" json:"IsDeleted"`
	Product     []Product `gorm:"foreignKey:CategoryID" json:"Product"`
	CreatedAt   time.Time `gorm:"column:CreatedAt" json:"CreatedAt"`
}

func (category *Category) TableName() string {
	return "Categories"
}

var CategoryColumns = struct {
	ID          string
	Name        string
	Description string
	IsDeleted   string
	CreatedAt   string
}{
	ID:          "Id",
	Name:        "Name",
	Description: "Description",
	IsDeleted:   "IsDeleted",
	CreatedAt:   "CreatedAt",
}

func (category *Category) BeforeCreate(tx *gorm.DB) (err error) {
	category.ID = uuid.New()
	return
}
