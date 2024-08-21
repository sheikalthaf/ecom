package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;column:Id" json:"Id"`
	Name      string    `gorm:"type:varchar(200);column:Name" json:"Name"`
	IsDeleted bool      `gorm:"type:bool;column:IsDeleted" json:"IsDeleted"`
	CreatedAt time.Time `gorm:"column:CreatedAt" json:"CreatedAt"`
}

func (tags *Tag) TableName() string {
	return "Tags"
}

var TagsColumns = struct {
	ID        string
	Name      string
	IsDeleted string
	CreatedAt string
}{
	ID:        "Id",
	Name:      "Name",
	IsDeleted: "IsDeleted",
	CreatedAt: "CreatedAt",
}

func (tags *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	tags.ID = uuid.New()
	return
}
